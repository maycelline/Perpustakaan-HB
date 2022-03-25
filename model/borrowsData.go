package model

type BorrowData struct {
	Borrows  []Borrowing `json:"borrowing"`
	Couriers []Courier   `json:"courier"`
}
