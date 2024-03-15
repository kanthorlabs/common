package signature

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSign(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		sign := Sign(key, data)

		signatures := strings.Split(sign, SignaturesDivider)
		require.Equal(st, len(versions), len(signatures))
	})
}
