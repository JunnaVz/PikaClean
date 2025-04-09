package models

import "github.com/google/uuid"

type User struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Address     string
	PhoneNumber string
	Email       string
	Password    string
}
