package password

import (
	"testing"

	"github.com/jaswdr/faker"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	t.Run("OK", func(st *testing.T) {
		pass := faker.New().Internet().Password()
		hash, err := Hash(pass)
		require.NoError(st, err)

		require.NoError(st, Compare(pass, hash))

	})

	t.Run("KO", func(st *testing.T) {
		_, err := Hash(testdata.Fake.Lorem().Sentence(100))
		require.NotNil(t, err)
	})
}
