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

func (client *client) Readiness() error {
	diff, err := client.read(Readiness)
	if err != nil {
		return err
	}

	delta := int64(client.conf.Liveness.Interval)
	if diff > delta {
		return fmt.Errorf("HEALTHCHECK.BACKGROUND.CLIENT.READINESS.TIMEOUT.ERROR: diff(%d) > delta(:%d)", diff, delta)
	}

	return nil
}

func (client *client) Liveness() error {
	diff, err := client.read(Liveness)
	if err != nil {
		return err
	}

	delta := int64(client.conf.Liveness.Interval)
	if diff > delta {
		return fmt.Errorf("HEALTHCHECK.BACKGROUND.CLIENT.LIVENESS.TIMEOUT.ERROR: diff(%d) > delta(:%d)", diff, delta)
	}

	return nil
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
