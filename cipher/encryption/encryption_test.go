package encryption

import (
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
	assert.NoError(t, err)

	original, err := Decrypt(key, ciphertext)
	assert.NoError(t, err)

	assert.Equal(t, data, original)
}

func genkey(n int) string {
	var str string
	count := n / 32
	for i := 0; i <= count; i++ {
		str += strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	return str[:n]
}
