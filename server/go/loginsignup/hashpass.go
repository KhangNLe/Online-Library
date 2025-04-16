package loginsignup

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass(text string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func ComparePass(hashPass string, pass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
	if err != nil {
		return false
	} else {
		return true
	}
}
