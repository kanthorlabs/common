package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func Encrypt(key, raw string) (string, error) {
	b := []byte(raw)
	ciphertext := make([]byte, aes.BlockSize+len(b))
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
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key, encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", errors.New("ENCRIPTION.DECRYPT.DECODE.ERROR")
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ENCRIPTION.DECRYPT.CIPHERTEXT.SIZE.ERROR")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("ENCRIPTION.DECRYPT.CIPHER_GENERATE.ERROR: %v", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
