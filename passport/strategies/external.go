package strategies

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/kanthorlabs/common/idx"
	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/passport/config"
	"github.com/kanthorlabs/common/passport/entities"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/sender"
	senderconfig "github.com/kanthorlabs/common/sender/config"
	senderentities "github.com/kanthorlabs/common/sender/entities"
)

var ExternalDefaultHeaders = http.Header{
	"Content-Type": []string{"application/json; charset=utf-8"},
	"Accept":       []string{"application/json"},
	"User-Agent":   []string{fmt.Sprintf("Kanthor/%s", project.GetVersion())},
}

// NewExternal creates a new external strategy instance what allows to authenticate users based on public external APIs.
func NewExternal(conf *config.External, logger logging.Logger) (Strategy, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	send, err := sender.New(senderconfig.Default, logger)
	if err != nil {
		return nil, err
	}

	return &external{conf: conf, logger: logger, send: send}, nil
}

type external struct {
	conf   *config.External
	logger logging.Logger
	send   sender.Send

	mu     sync.Mutex
	status int
}

func (instance *external) Connect(ctx context.Context) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if instance.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	instance.status = patterns.StatusConnected
	return nil
}

func (instance *external) Readiness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	url := fmt.Sprintf("%s/healthz/readiness", instance.conf.Uri)
	return sender.Check(instance.send, url)
}

func (instance *external) Liveness() error {
	if instance.status == patterns.StatusDisconnected {
		return nil
	}
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	url := fmt.Sprintf("%s/healthz/liveness", instance.conf.Uri)
	return sender.Check(instance.send, url)
}

func (instance *external) Disconnect(ctx context.Context) error {
	if instance.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	instance.status = patterns.StatusDisconnected
	return nil
}

func (instance *external) Register(ctx context.Context, acc entities.Account) error {
	body, err := json.Marshal(acc)
	if err != nil {
		return err
	}

	req := &senderentities.Request{
		Method:  http.MethodPost,
		Headers: ExternalDefaultHeaders,
		Uri:     fmt.Sprintf("%s/authn/register", instance.conf.Uri),
		Body:    body,
	}
	req.Headers.Set("Idempotency-Key", idx.New("ik"))
	res, err := instance.send(ctx, req)
	if err != nil {
		return err
	}

	if !res.Ok() {
		return errors.New(res.StatusText())
	}

	return nil
}

func (instance *external) Login(ctx context.Context, creds entities.Credentials) (*entities.Tokens, error) {
	body, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	req := &senderentities.Request{
		Method:  http.MethodPost,
		Headers: ExternalDefaultHeaders,
		Uri:     fmt.Sprintf("%s/authn/login", instance.conf.Uri),
		Body:    body,
	}
	req.Headers.Set("Idempotency-Key", idx.New("ik"))
	res, err := instance.send(ctx, req)
	if err != nil {
		return nil, err
	}

	if !res.Ok() {
		return nil, errors.New(res.StatusText())
	}

	var tokens entities.Tokens
	if err := json.Unmarshal(res.Body, &tokens); err != nil {
		return nil, err
	}

	if err := tokens.Validate(); err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (instance *external) Logout(ctx context.Context, tokens entities.Tokens) error {
	req := &senderentities.Request{
		Method:  http.MethodPost,
		Headers: ExternalDefaultHeaders,
		Uri:     fmt.Sprintf("%s/account/logout", instance.conf.Uri),
		Body:    nil,
	}
	req.Headers.Set("Authorization", tokens.Access)
	req.Headers.Set("Idempotency-Key", idx.New("ik"))
	res, err := instance.send(ctx, req)
	if err != nil {
		return err
	}

	if !res.Ok() {
		return errors.New(res.StatusText())
	}

	return nil
}

func (instance *external) Verify(ctx context.Context, tokens entities.Tokens) (*entities.Account, error) {
	req := &senderentities.Request{
		Method:  http.MethodGet,
		Headers: ExternalDefaultHeaders,
		Uri:     fmt.Sprintf("%s/account/me", instance.conf.Uri),
		Body:    nil,
	}
	req.Headers.Set("Authorization", tokens.Access)
	req.Headers.Set("Idempotency-Key", idx.New("ik"))
	res, err := instance.send(ctx, req)
	if err != nil {
		return nil, err
	}

	if !res.Ok() {
		return nil, errors.New(res.StatusText())
	}

	var account entities.Account
	if err := json.Unmarshal(res.Body, &account); err != nil {
		return nil, err
	}

	if err := account.Validate(); err != nil {
		return nil, err
	}

	return &account, nil
}

func (instance *external) Deactivate(ctx context.Context, username string, at int64) error {
	return errors.New("PASSPORT.EXTERNAL.DEACTIVATE.UNIMPLEMENT.ERROR")
}

func (instance *external) List(ctx context.Context, usernames []string) ([]*entities.Account, error) {
	return nil, errors.New("PASSPORT.EXTERNAL.LIST.UNIMPLEMENT.ERROR")
}

func (instance *external) Update(ctx context.Context, account entities.Account) error {
	return errors.New("PASSPORT.EXTERNAL.UPDATE.UNIMPLEMENT.ERROR")
}
