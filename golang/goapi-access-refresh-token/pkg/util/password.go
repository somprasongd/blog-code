package util

import (
	"github.com/matthewhartstonge/argon2"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	encoded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encoded), err
}

func HashPasswordArgon(password string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))
	return string(encoded), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordHashArgon(password, hash string) bool {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(hash))
	if err != nil {
		return false
	}
	return ok
}
