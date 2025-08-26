package common

import "golang.org/x/crypto/bcrypt"

// hashing the password

// the function expects password as string, and returns the hashed password as string, and also returns an error
func HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// function to compare passwords, returns a boolean
func ComparePasswords(password, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
