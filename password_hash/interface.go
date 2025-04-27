package password_hash

// PasswordHash defines the interface for password hashing operations.
type PasswordHash interface {
	// GetHash generates a secure hash from a plaintext string.
	// Returns the hashed password as a string and any error that occurred.
	GetHash(stringToHash string) (string, error)

	// CompareHashAndPassword verifies if a plaintext password matches a hashed password.
	// Returns true if the passwords match, false otherwise.
	CompareHashAndPassword(hashedPassword, plainPassword string) bool
}
