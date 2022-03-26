package model

import "time"

type BorrowData struct {
	Borrows  []Borrowing `json:"borrowing"`
	Couriers []Courier   `json:"courier"`
}

type BorrowDataHTML struct {
	Branch      Branch
	User        User
	Borrows     []Borrowing
	Courier     Courier
	CourierCome time.Time
}
