package background

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kanthorlabs/common/healthcheck"
	"github.com/kanthorlabs/common/healthcheck/config"
)

func NewClient(conf *config.Config) (healthcheck.Client, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &client{conf: conf}, nil
}

type client struct {
	conf *config.Config
}

func (client *client) Readiness() (err error) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(client.conf.Readiness.Timeout))
	defer ticker.Stop()

	var count int
	for range ticker.C {
		count++
		if count >= client.conf.Readiness.MaxTry {
			return
		}

		_, err = client.read(Readiness)
		if err != nil {
			continue
		}

		err = nil
		return
	}

	return
}

func (client *client) Liveness() (err error) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(client.conf.Liveness.Timeout))
	defer ticker.Stop()

	var diff, delta int64
	var count int
	for range ticker.C {
		count++
		if count >= client.conf.Liveness.MaxTry {
			return
		}

		diff, err = client.read(Liveness)
		if err != nil {
			continue
		}

		delta = int64(client.conf.Liveness.Timeout * client.conf.Liveness.MaxTry)
		if diff > delta {
			err = fmt.Errorf("HEALTHCHECK.BACKGROUND.CLIENT.LIVENESS.TIMEOUT.ERROR: diff(%d) > delta(:%d)", diff, delta)
			continue
		}

		err = nil
		return
	}

	return
}

func (client *client) read(name string) (int64, error) {
	file := fmt.Sprintf("%s.%s", client.conf.Dest, name)
	data, err := os.ReadFile(file)
	if err != nil {
		return 0, err
	}

	prev, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return 0, err
	}

	now := time.Now().UnixMilli()
	return now - prev, nil
}
