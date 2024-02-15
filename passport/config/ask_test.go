package config

import (
	"testing"

	"github.com/kanthorlabs/common/passport/entities"
	"github.com/stretchr/testify/require"
)

func TestAsk(t *testing.T) {
	t.Run(".Validate", func(st *testing.T) {
		st.Run("KO - no account", func(sst *testing.T) {
			conf := &Ask{Accounts: make([]entities.Account, 0)}
			require.ErrorContains(sst, conf.Validate(), "PASSPORT.STRATEGY.ASK.CONFIG.ACCOUNTS")
		})

		st.Run("KO - account error", func(sst *testing.T) {
			conf := &Ask{Accounts: make([]entities.Account, 1)}
			require.ErrorContains(sst, conf.Validate(), "PASSPORT.ACCOUNT")
		})
	})
}
