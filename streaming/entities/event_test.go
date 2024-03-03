package entities

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	e := &Event{Id: uuid.NewString()}
	require.ErrorContains(t, e.Validate(), "STREAMING.EVENT")

	data, err := json.Marshal(e)
	require.NoError(t, err)
	require.Equal(t, string(data), e.String())
}
