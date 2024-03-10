package password

import "golang.org/x/crypto/bcrypt"

// Hash returns the bcrypt hash of the provided string with the default cost.
func Hash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

// Compare compares the provided string with the provided hash.
func Compare(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
