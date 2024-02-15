package password

import "golang.org/x/crypto/bcrypt"

func HashString(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

func CompareString(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
