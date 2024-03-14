package encryption

import (
	"crypto/aes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	key := genkey(32)
	data := faker.New().Lorem().Sentence(256)

	ciphertext, err := Encrypt(key, data)
	assert.Nil(t, err)

	original, err := Decrypt(key, ciphertext)
	assert.Nil(t, err)

	assert.Equal(t, data, original)
}

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

func TestEncryption_Decrypt(t *testing.T) {
	t.Run("KO - encrypted base64 decode error", func(t *testing.T) {
		encrypted := faker.New().Lorem().Sentence(256)
		_, err := Decrypt("", encrypted)
		assert.ErrorContains(t, err, "ENCRIPTION.DECRYPT.DECODE")
	})

	t.Run("KO - encrypted text size error", func(t *testing.T) {
		encrypted := base64.StdEncoding.EncodeToString([]byte(genkey(aes.BlockSize - 1)))
		_, err := Decrypt(genkey(32), encrypted)
		assert.ErrorContains(t, err, "ENCRIPTION.DECRYPT.CIPHERTEXT.SIZE")
	})
}

func genkey(n int) string {
	var str string
	count := n / 32
	for i := 0; i <= count; i++ {
		str += strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	return str[:n]
}
