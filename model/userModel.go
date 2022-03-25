package model

import "time"

type User struct {
	ID        int       `json:"idUser"`
	FullName  string    `json:"fullName"`
	Username  string    `json:"userName"`
	BirthDate time.Time `json:"birthDate"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Balance   int       `json:"balance"`
}
