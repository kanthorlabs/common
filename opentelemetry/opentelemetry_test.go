package opentelemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupAndTeardown(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		require.NoError(st, Setup(context.Background()))
		require.NoError(st, Teardown(context.Background()))
	})
}
