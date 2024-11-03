package middlewares

import "golang.org/x/crypto/bcrypt"

func BcryptPassword(s string) (h string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	h = string(hash)
	return
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
