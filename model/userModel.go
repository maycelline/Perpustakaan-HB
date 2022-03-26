package model

import "time"

type User struct {
	ID                int       `json:"idUser,omitempty"`
	FullName          string    `json:"fullName,omitempty"`
	Username          string    `json:"userName,omitempty"`
	BirthDate         time.Time `json:"birthDate,omitempty"`
	Phone             string    `json:"phone,omitempty"`
	Email             string    `json:"email,omitempty"`
	Address           string    `json:"address,omitempty"`
	AdditionalAddress string    `json:"additionalAddress,omitempty"`
	Password          string    `json:"password,omitempty"`
	UserType          string    `json:"userType,omitempty"`
}

type Member struct {
	User    User
	Balance int `json:"balance"`
}

type Admin struct {
	User   User
	Branch Branch
}
