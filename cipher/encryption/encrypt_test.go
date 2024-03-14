package encryption

import (
	"testing"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestEncryption_Encrypt(t *testing.T) {
	t.Run("KO - empty data", func(t *testing.T) {
		key := genkey(32)
		s, err := Encrypt(key, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, s)
	})

	t.Run("KO - empty key error", func(t *testing.T) {
		data := faker.New().Lorem().Sentence(256)
		_, err := Encrypt("", data)
		assert.ErrorContains(t, err, "ENCRIPTION.ENCRYPT.CIPHER_GENERATE")
	})
}
