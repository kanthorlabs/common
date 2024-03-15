package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

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

	// make sure the checksum is correct
	if len(ciphertext) < 32 {
		return "", errors.New("ENCRIPTION.DECRYPT.CHECKSUM.SIZE.ERROR")
	}

	// extract the checksum
	compare := ciphertext[len(ciphertext)-32:]
	data := ciphertext[:len(ciphertext)-32]
	checksum := md5.Sum([]byte(data))
	if hex.EncodeToString(checksum[:]) != string(compare) {
		return "", errors.New("ENCRIPTION.DECRYPT.CHECKSUM.ERROR")
	}

	return string(data), nil
}

// DecryptAny decrypts the encrypted text using the keys provided
// The usecase is you want to rotate the key and still be able to decrypt the old encrypted
// So you you rotate the key to obtain the new key, add it into the beginning of the keys slice
// After that the first key will be used to encrypt the data, and the rest of the keys will be used to decrypt the old data
func DecryptAny(keys []string, encrypted string) (string, error) {
	for _, key := range keys {
		decrypted, err := Decrypt(key, encrypted)
		if err == nil {
			return decrypted, nil
		}
	}
	return "", errors.New("ENCRIPTION.DECRYPT.ERROR")
}
