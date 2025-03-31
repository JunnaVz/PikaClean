package password_hash

import (
	"golang.org/x/crypto/bcrypt"
)

type bcryptHash struct {
}

func NewPasswordHash() PasswordHash {
	return &bcryptHash{}
}

func (b *bcryptHash) GetHash(stringToHash string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(stringToHash), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func (b *bcryptHash) CompareHashAndPassword(hashedPassword, plainPassword string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return res == nil
}
