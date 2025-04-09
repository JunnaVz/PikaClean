package models

import "github.com/google/uuid"

type Worker struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Address     string
	PhoneNumber string
	Email       string
	Role        int
	Password    string
}

const ManagerRole = 1
const MasterRole = 2

var WorkerRole = map[int]string{
	ManagerRole: "Менеджер",
	MasterRole:  "Мастер",
}

func (w Worker) DisplayRole() string {
	return WorkerRole[w.Role]
}

func (w Worker) FullName() string {
	return w.Name + " " + w.Surname
}
