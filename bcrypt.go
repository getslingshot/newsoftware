package main

import (
	"golang.org/x/crypto/bcrypt"
)

// Bcrypt Clear
func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

// Bcrypt Crypt
func Crypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// Bcrypt Compare Function
func Compare(encrypted_password []byte, plain_password []byte) error {
	return bcrypt.CompareHashAndPassword(encrypted_password, plain_password)
}
