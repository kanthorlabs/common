package signature

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

var key = uuid.NewString()
var data = fmt.Sprintf(
	"%s.%s.%d",
	testdata.Fake.Lorem().Sentence(1),
	testdata.Fake.Lorem().Sentence(1),
	time.Now().UTC().UnixMilli(),
)

func TestSign(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		sign := Sign(key, data)

		signatures := strings.Split(sign, SignatureDivider)
		require.Equal(st, len(versions), len(signatures))
	})
}

func TestVerify(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		sign := Sign(key, data)

		require.Nil(st, Verify(key, data, sign))
	})

	t.Run("KO - malformed signature", func(st *testing.T) {
		sign := "v1"
		require.ErrorContains(st, Verify(key, data, sign), "SIGNATURE.VERIFY.NOT_MATCH.")
	})

	t.Run("KO - missmatch version", func(st *testing.T) {
		sign := "v0="
		require.ErrorContains(st, Verify(key, data, sign), "SIGNATURE.VERIFY.NOT_MATCH.")
	})

	t.Run("KO - missmatch signature", func(st *testing.T) {
		sign := "v1=" + testdata.Fake.Lorem().Sentence(1)
		require.ErrorContains(st, Verify(key, data, sign), "SIGNATURE.VERIFY.NOT_MATCH.")
	})
}
