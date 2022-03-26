package model

import "time"

type User struct {
	ID          int       `json:"idUser"`
	FullName    string    `json:"fullName"`
	UserName    string    `json:"userName"`
	BirthDate   time.Time `json:"birthDate"`
	PhoneNumber string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	Password    string    `json:"password"`
	Balance     int       `json:"balance"`
}
