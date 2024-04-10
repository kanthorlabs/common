package sender

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/kanthorlabs/common/project"
	"github.com/kanthorlabs/common/sender/entities"
)

func Check(send Send, url string) error {
	req := &entities.Request{
		Method: http.MethodGet,
		Headers: http.Header{
			"Accept":     []string{"*/*; charset=utf-8"},
			"User-Agent": []string{fmt.Sprintf("Kanthor/%s", project.GetVersion())},
		},
		Uri: url,
	}
	res, err := send(context.Background(), req)
	if err != nil {
		return err
	}

	if !res.Ok() {
		return errors.New(res.StatusText())
	}

	return nil
}
