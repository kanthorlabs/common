package background

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kanthorlabs/common/healthcheck"
	"github.com/kanthorlabs/common/healthcheck/config"
)

var Readiness = "readiness"
var Liveness = "liveness"

func NewServer(conf *config.Config) (healthcheck.Server, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	return &server{
		conf:       conf,
		terminated: make(chan int64, 1),
	}, nil
}

type server struct {
	conf *config.Config

	terminated chan int64
}

func (server *server) Connect(ctx context.Context) error {
	return nil
}

func (server *server) Disconnect(ctx context.Context) error {
	server.terminated <- time.Now().UnixMilli()
	return nil
}

func (server *server) Readiness(check func() error) error {
	if err := check(); err != nil {
		return err
	}
	if err := server.write(Readiness); err != nil {
		return err
	}

	return nil
}

func (server *server) Liveness(check func() error) error {
	ticker := time.NewTicker(time.Millisecond * time.Duration(server.conf.Liveness.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-server.terminated:
			return nil
		case <-ticker.C:
			if err := check(); err != nil {
				return err
			}
			if err := server.write(Liveness); err != nil {
				return err
			}
		}
	}
}

func (server *server) write(name string) error {
	data := fmt.Sprintf("%d", time.Now().UnixMilli())

	file := fmt.Sprintf("%s.%s", server.conf.Dest, name)
	return os.WriteFile(file, []byte(data), os.ModePerm)
}
