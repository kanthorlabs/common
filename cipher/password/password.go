package password

import "golang.org/x/crypto/bcrypt"

// HashString returns the bcrypt hash of the provided string with the default cost.
func HashString(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

// CompareString compares the provided string with the provided hash.
func CompareString(pass, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}
