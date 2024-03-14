package encryption

import (
	"crypto/aes"
	"encoding/base64"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/kanthorlabs/common/testdata"
	"github.com/stretchr/testify/assert"
)

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

func TestEncryption_DecryptAny(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		count := testdata.Fake.IntBetween(5, 10)
		var keys []string
		for i := 0; i < count; i++ {
			keys = append(keys, genkey(32))
		}

		key := keys[testdata.Fake.IntBetween(0, count-1)]
		data := faker.New().Lorem().Sentence(256)
		encrypted, err := Encrypt(key, data)
		assert.NoError(t, err)

		decrypted, err := DecryptAny(keys, encrypted)
		assert.NoError(t, err)
		assert.Equal(t, data, decrypted)
	})

	t.Run("KO - all keys failed", func(t *testing.T) {
		encrypted := faker.New().Lorem().Sentence(256)
		_, err := DecryptAny([]string{genkey(32), genkey(32)}, encrypted)
		assert.ErrorContains(t, err, "ENCRIPTION.DECRYPT")
	})
}
