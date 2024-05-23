package testify

import (
	"context"
	"testing"

	"github.com/kanthorlabs/common/patterns"
	"github.com/stretchr/testify/require"
)

func AssertConnect(t *testing.T, instance patterns.Connectable, err error) {
	require.NoError(t, instance.Connect(context.Background()))
	require.ErrorIs(t, instance.Connect(context.Background()), err)
}

func AssertReadiness(t *testing.T, instance patterns.Connectable, err error) {
	require.ErrorIs(t, instance.Readiness(), err)
	require.NoError(t, instance.Connect(context.Background()))
	require.NoError(t, instance.Readiness())
	require.NoError(t, instance.Disconnect(context.Background()))
	require.NoError(t, instance.Readiness())
}

func AssertLiveness(t *testing.T, instance patterns.Connectable, err error) {
	require.ErrorIs(t, instance.Liveness(), err)
	require.NoError(t, instance.Connect(context.Background()))
	require.NoError(t, instance.Liveness())
	require.NoError(t, instance.Disconnect(context.Background()))
	require.NoError(t, instance.Liveness())
}

func AssertDisconnect(t *testing.T, instance patterns.Connectable, err error) {
	require.ErrorIs(t, instance.Disconnect(context.Background()), err)
	require.NoError(t, instance.Connect(context.Background()))
	require.NoError(t, instance.Disconnect(context.Background()))
}
