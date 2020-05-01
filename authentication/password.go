package authentication

import (
	"github.com/rkuris/journey/database"
	"golang.org/x/crypto/bcrypt"
)

// LoginIsCorrect checks the username/password combo and says ok or not
func LoginIsCorrect(name string, password string) bool {
	hashedPassword, err := database.RetrieveHashedPasswordForUser([]byte(name))
	if len(hashedPassword) == 0 || err != nil { // len(hashedPassword) == 0 probably not needed.
		// User name likely doesn't exist
		return false
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return false
	}
	return true
}

// EncryptPassword uses bcrypt to generate a password for a user
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
