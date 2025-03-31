package password_hash

type PasswordHash interface {
	GetHash(stringToHash string) (string, error)
	CompareHashAndPassword(hashedPassword, plainPassword string) bool
}
