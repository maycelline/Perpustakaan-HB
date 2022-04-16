package model

import "time"

type BorrowData struct {
	Borrows  []Borrowing `json:"borrowings,omitempty"`
	Couriers []Courier   `json:"couriers,omitempty"`
}

type BorrowDataHTML struct {
	Branch      Branch
	User        User
	Books       []Book
	Courier     Courier
	CourierCome time.Time
}
