package password

import (
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	pass := faker.New().Internet().Password()
	hash, err := HashString(pass)
	require.Nil(t, err)

	require.Nil(t, CompareString(pass, hash))
}
