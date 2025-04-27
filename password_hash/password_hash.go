// Package password_hash provides functionality for securely hashing and verifying passwords.
// It implements bcrypt hashing algorithm for password storage.
package password_hash

import (
	"golang.org/x/crypto/bcrypt"
)

// bcryptHash implements the PasswordHash interface using bcrypt algorithm.
type bcryptHash struct {
}

// NewPasswordHash creates and returns a new PasswordHash implementation.
// This is the entry point for creating password hashers.
func NewPasswordHash() PasswordHash {
	return &bcryptHash{}
}

// GetHash generates a secure bcrypt hash from a plaintext string.
// It uses the default cost for the hash computation.
func (b *bcryptHash) GetHash(stringToHash string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(stringToHash), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// CompareHashAndPassword verifies if a plaintext password matches a previously hashed password.
// Returns true if the passwords match, false otherwise.
func (b *bcryptHash) CompareHashAndPassword(hashedPassword, plainPassword string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return res == nil
}
