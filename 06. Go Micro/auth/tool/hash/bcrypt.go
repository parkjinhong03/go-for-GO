package hash

import "golang.org/x/crypto/bcrypt"

func BcryptGenerate(pwd string, cost int) (hash string, err error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), cost)
	if err != nil {
		return
	}
	hash = string(b)
	return
}