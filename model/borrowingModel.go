package model

import "time"

type Borrowing struct {
	ID         int       `json:"idBorrowing"`
	Book       Book      `json:"book"`
	BorrowDate time.Time `json:"borrowDate"`
	ReturnDate time.Time `json:"returnDate"`
}
