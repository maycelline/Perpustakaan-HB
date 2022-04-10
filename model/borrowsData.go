package model

import "time"

type BorrowData struct {
	// User   User   `json:"user"`
	// Branch Branch `json:"branch"`
	// Books  []Book `json:"books"`
	Borrows  []Borrowing `json:"borrowings,omitempty"`
	Couriers []Courier   `json:"couriers,omitempty"`
}

type BorrowDataHTML struct {
	Branch      Branch
	User        User
	Borrows     []Borrowing
	Courier     Courier
	CourierCome time.Time
}
