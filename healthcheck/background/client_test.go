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
	t.Run("NewClient - KO", func(st *testing.T) {

		_, err := NewClient(&config.Config{})
		require.NotNil(st, err)
	})

	t.Run(".Readiness", func(st *testing.T) {
		name := testdata.Fake.UUID().V4()
		conf := config.Default(name, 100)
		client, _ := NewClient(conf)

		st.Run("KO - read status", func(sst *testing.T) {
			require.ErrorIs(sst, client.Readiness(), os.ErrNotExist)
		})
	})

	t.Run(".Liveness", func(st *testing.T) {
		name := testdata.Fake.UUID().V4()
		conf := config.Default(name, 100)
		client, _ := NewClient(conf)

		st.Run("KO - read status", func(sst *testing.T) {
			require.ErrorIs(sst, client.Liveness(), os.ErrNotExist)
		})

		st.Run("KO - incorrect data", func(sst *testing.T) {
			filepath := fmt.Sprintf("%s.liveness", conf.Dest)
			data := time.Now().Format(time.RFC3339)

			err := os.WriteFile(filepath, []byte(data), os.ModePerm)
			require.NoError(sst, err)

			require.ErrorContains(sst, client.Liveness(), "strconv.ParseInt")
		})
		st.Run("KO - no signal", func(sst *testing.T) {
			filepath := fmt.Sprintf("%s.liveness", conf.Dest)
			data := fmt.Sprintf("%d", time.Now().Add(-time.Hour).UnixMilli())

			err := os.WriteFile(filepath, []byte(data), os.ModePerm)
			require.NoError(sst, err)

			require.ErrorContains(sst, client.Liveness(), "HEALTHCHECK.BACKGROUND.CLIENT.LIVENESS.TIMEOUT.ERROR")
		})
	})
}
