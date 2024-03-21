package streaming

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/kanthorlabs/common/logging"
	"github.com/kanthorlabs/common/patterns"
	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/streaming/config"
	"github.com/kanthorlabs/common/validator"
	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func NewNats(conf *config.Config, logger logging.Logger) (Stream, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &nats{conf: conf, logger: logger}, nil
}

type nats struct {
	conf   *config.Config
	logger logging.Logger

	conn *natsio.Conn
	js   jetstream.JetStream

	mu          sync.Mutex
	status      int
	publishers  map[string]Publisher
	subscribers map[string]Subscriber
}

func (streaming *nats) Connect(ctx context.Context) error {
	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	conn, err := streaming.connect()
	if err != nil {
		return err
	}
	streaming.conn = conn

	js, err := jetstream.New(conn)
	if err != nil {
		return err
	}
	streaming.js = js

	stream, err := streaming.stream(ctx)
	if err != nil {
		return err
	}

	streaming.status = patterns.StatusConnected
	streaming.logger.Infow(
		"connected",
		"stream_name", stream.CachedInfo().Config.Name,
		"stream_created_at", stream.CachedInfo().Created.Format(time.RFC3339),
	)
	return nil
}

func (streaming *nats) Readiness() error {
	if streaming.status == patterns.StatusDisconnected {
		return nil
	}
	if streaming.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if !streaming.conn.IsConnected() {
		return ErrNotConnected
	}

	return nil
}

func (streaming *nats) Liveness() error {
	if streaming.status == patterns.StatusDisconnected {
		return nil
	}
	if streaming.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	if !streaming.conn.IsConnected() {
		return ErrNotConnected
	}

	return nil
}

func (streaming *nats) connect() (*natsio.Conn, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	opts := []natsio.Option{
		natsio.Name(hostname),
		natsio.ReconnectWait(3 * time.Second),
		natsio.Timeout(3 * time.Second),
		natsio.MaxReconnects(9),
		natsio.DisconnectErrHandler(func(c *natsio.Conn, err error) {
			if err != nil {
				streaming.logger.Errorw("STREAMING.NATS.DISCONNECTED.ERROR", "error", err)
				return
			}
		}),
		natsio.ReconnectHandler(func(conn *natsio.Conn) {
			streaming.logger.Warnw("STREAMING.NATS.RECONNECT", "url", conn.ConnectedUrl())
		}),
		natsio.ErrorHandler(func(c *natsio.Conn, s *natsio.Subscription, err error) {
			if err == natsio.ErrSlowConsumer {
				count, bytes, err := s.Pending()
				streaming.logger.Errorw("STREAMING.NATS.SLOW_CONSUMER.ERROR", "error", err, "subject", s.Subject, "count", count, "bytes", bytes)
				return
			}

			streaming.logger.Errorw("STREAMING.NATS.ERROR", "error", err)
		}),
	}

	return natsio.Connect(streaming.conf.Uri, opts...)
}

func (streaming *nats) stream(ctx context.Context) (jetstream.Stream, error) {
	return streaming.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		// non-editable
		Name: streaming.conf.Name,
		// editable
		Replicas: streaming.conf.Nats.Replicas,
		// project.Subject(">") mean we accept all subjects that belong to the configured project and tier
		Subjects:          []string{project.Subject(">")},
		MaxBytes:          streaming.conf.Nats.Limits.Bytes,
		MaxMsgSize:        streaming.conf.Nats.Limits.MsgSize,
		MaxMsgsPerSubject: streaming.conf.Nats.Limits.MsgCount,
		MaxAge:            time.Millisecond * time.Duration(streaming.conf.Nats.Limits.MsgAge),
		Duplicates:        time.Millisecond * time.Duration(streaming.conf.Nats.Limits.MsgAge),
		Retention:         jetstream.LimitsPolicy,
		Discard:           jetstream.DiscardOld,
	})
}

func (streaming *nats) Disconnect(ctx context.Context) error {
	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	streaming.status = patterns.StatusDisconnected
	streaming.logger.Info("disconnected")

	var retruning error
	if len(streaming.publishers) > 0 {
		streaming.publishers = nil
	}
	if len(streaming.subscribers) > 0 {
		streaming.subscribers = nil
	}
	if err := streaming.conn.Drain(); err != nil {
		retruning = errors.Join(retruning, err)
	}
	streaming.conn = nil

	streaming.js = nil

	return retruning
}

func (streaming *nats) Publisher(name string) (Publisher, error) {
	if streaming.status != patterns.StatusConnected {
		return nil, ErrNotConnected
	}

	err := validator.StringAlphaNumericUnderscore("STREAMING.PUBLISHER.NAME", name)()
	if err != nil {
		return nil, err
	}

	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.publishers == nil {
		streaming.publishers = map[string]Publisher{}
	}

	if publisher, exist := streaming.publishers[name]; exist {
		return publisher, nil
	}

	publisher := &NatsPublisher{
		name:   name,
		conf:   streaming.conf,
		logger: streaming.logger.With("streaming.publisher", name),
		js:     streaming.js,
	}
	streaming.publishers[name] = publisher

	return publisher, nil
}

func (streaming *nats) Subscriber(name string) (Subscriber, error) {
	if streaming.status != patterns.StatusConnected {
		return nil, ErrNotConnected
	}

	err := validator.StringAlphaNumericUnderscore("STREAMING.SUBSCRIBER.NAME", name)()
	if err != nil {
		return nil, err
	}

	streaming.mu.Lock()
	defer streaming.mu.Unlock()

	if streaming.subscribers == nil {
		streaming.subscribers = map[string]Subscriber{}
	}

	if subscriber, exist := streaming.subscribers[name]; exist {
		return subscriber, nil
	}

	subscriber := &NatsSubscriber{
		name:   name,
		conf:   streaming.conf,
		logger: streaming.logger.With("streaming.subscriber", name),
		js:     streaming.js,
	}
	streaming.subscribers[name] = subscriber

	return subscriber, nil
}
