package entities

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	req := &Request{}
	require.ErrorContains(t, req.Validate(), "SENDER.REQUEST")
}
