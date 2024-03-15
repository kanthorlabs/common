package signature

import (
	"testing"

	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		sign := Sign(key, data)

		require.NoError(st, Verify(key, data, sign))
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

func TestVerifyAny(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		sign := Sign(key, data)

		require.NoError(st, VerifyAny([]string{key}, data, sign))
	})

	t.Run("OK - rotate key", func(st *testing.T) {
		sign := Sign(key, data)

		require.NoError(st, VerifyAny([]string{testdata.Fake.Lorem().Sentence(1), key}, data, sign))
	})

	t.Run("KO - missmatch signature", func(st *testing.T) {
		sign := Sign(key, data)

		require.ErrorContains(st, VerifyAny([]string{testdata.Fake.Lorem().Sentence(1)}, data, sign), "SIGNATURE.VERIFY.NOT_MATCH.")
	})
}
