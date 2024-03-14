package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

func Encrypt(key, raw string) (string, error) {
	// generate md5 checksum of the data so we can verify in decryption later
	checksum := md5.Sum([]byte(raw))
	data := []byte(raw + hex.EncodeToString(checksum[:]))

	ciphertext := make([]byte, aes.BlockSize+len(data))
	// fill random nonce
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.New("ENCRIPTION.ENCRYPT.NONCE_GENERATE.ERROR")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ENCRIPTION.ENCRYPT.CIPHER_GENERATE.ERROR: %v", err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
