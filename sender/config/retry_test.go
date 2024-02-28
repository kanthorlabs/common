package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	conf := Retry{}
	require.ErrorContains(t, conf.Validate(), "SENDER.CONFIG.RETRY")
}
