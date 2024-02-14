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

func TestClient(t *testing.T) {
	t.Run("ko in init", func(st *testing.T) {

		_, err := NewClient(&config.Config{})
		require.NotNil(st, err)
	})

	t.Run("readiness", func(st *testing.T) {
		name := testdata.Fake.UUID().V4()
		conf := config.Default(name, 100)
		client, _ := NewClient(conf)

		st.Run("ko because of read status", func(sst *testing.T) {
			require.ErrorIs(sst, client.Readiness(), os.ErrNotExist)
		})
	})

	t.Run("liveness", func(st *testing.T) {
		name := testdata.Fake.UUID().V4()
		conf := config.Default(name, 100)
		client, _ := NewClient(conf)

		st.Run("ko because of read status", func(sst *testing.T) {
			require.ErrorIs(sst, client.Liveness(), os.ErrNotExist)
		})

		st.Run("ko because of incorrect data", func(sst *testing.T) {
			filepath := fmt.Sprintf("%s.liveness", conf.Dest)
			data := time.Now().Format(time.RFC3339)

			err := os.WriteFile(filepath, []byte(data), os.ModePerm)
			require.Nil(sst, err)

			require.ErrorContains(sst, client.Liveness(), "strconv.ParseInt")
		})
		st.Run("ko because of delay", func(sst *testing.T) {
			filepath := fmt.Sprintf("%s.liveness", conf.Dest)
			data := fmt.Sprintf("%d", time.Now().Add(-time.Hour).UnixMilli())

			err := os.WriteFile(filepath, []byte(data), os.ModePerm)
			require.Nil(sst, err)

			require.ErrorContains(sst, client.Liveness(), "HEALTHCHECK.BACKGROUND.CLIENT.LIVENESS.TIMEOUT.ERROR")
		})
	})
}
