package model

import "time"

type User struct {
	ID                int       `json:"idUser"`
	FullName          string    `json:"fullName"`
	Username          string    `json:"userName"`
	BirthDate         time.Time `json:"birthDate"`
	Phone             string    `json:"phone"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	AdditionalAddress string    `json:"additionalAddress"`
	Password          string    `json:"password"`
	UserType          string    `json:"userType"`
}

type Member struct {
	User    User
	Balance int `json:"balance"`
}

type Admin struct {
	User   User
	Branch Branch
}
