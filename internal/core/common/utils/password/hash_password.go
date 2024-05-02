package password

import "golang.org/x/crypto/bcrypt"

func HasPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashPassword := string(hash)
	return hashPassword, err
}
