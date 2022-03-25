package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func ValidatePassword(password, hashed []byte) bool {
	return bcrypt.CompareHashAndPassword(hashed, password) == nil
}
