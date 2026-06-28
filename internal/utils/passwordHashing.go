package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// bcrypt.DefaultCost is currently 10. You can increase this (up to 31)
	// as servers get faster over the years.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain text password with a bcrypt hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// If err is nil, the passwords match perfectly
	return err == nil
}
