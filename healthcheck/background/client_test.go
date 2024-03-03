package background

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kanthorlabs/common/healthcheck/config"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestClient_New(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		name := testdata.Fake.UUID().V4()
		conf := config.Default(name, 100)
		_, err := NewClient(conf)
		require.NoError(st, err)
	})
	t.Run("KO", func(st *testing.T) {
		_, err := NewClient(&config.Config{})
		require.NotNil(st, err)
	})
}

func TestClient_Readiness(t *testing.T) {
	name := testdata.Fake.UUID().V4()
	conf := config.Default(name, 100)
	client, _ := NewClient(conf)

	t.Run("KO - read status", func(st *testing.T) {
		require.ErrorIs(st, client.Readiness(), os.ErrNotExist)
	})
}

func TestClient_Liveness(t *testing.T) {
	name := testdata.Fake.UUID().V4()
	conf := config.Default(name, 100)
	client, _ := NewClient(conf)

	t.Run("KO - read status", func(st *testing.T) {
		require.ErrorIs(st, client.Liveness(), os.ErrNotExist)
	})

	t.Run("KO - incorrect data", func(st *testing.T) {
		filepath := fmt.Sprintf("%s.liveness", conf.Dest)
		data := time.Now().Format(time.RFC3339)

		err := os.WriteFile(filepath, []byte(data), os.ModePerm)
		require.NoError(st, err)

		require.ErrorContains(st, client.Liveness(), "strconv.ParseInt")
	})
	t.Run("KO - no signal", func(st *testing.T) {
		filepath := fmt.Sprintf("%s.liveness", conf.Dest)
		data := fmt.Sprintf("%d", time.Now().Add(-time.Hour).UnixMilli())

		err := os.WriteFile(filepath, []byte(data), os.ModePerm)
		require.NoError(st, err)

		require.ErrorContains(st, client.Liveness(), "HEALTHCHECK.BACKGROUND.CLIENT.LIVENESS.TIMEOUT.ERROR")
	})
}
