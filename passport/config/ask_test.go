package config

import (
	"testing"

	"github.com/kanthorlabs/common/passport/entities"
	"github.com/stretchr/testify/require"
)

func TestAsk(t *testing.T) {
	t.Run(".Validate/KO", func(st *testing.T) {
		conf := &Ask{Accounts: make([]entities.Account, 0)}
		require.ErrorContains(st, conf.Validate(), "PASSPORT.STRATEGY.ASK.CONFIG.ACCOUNTS")
	})
}
