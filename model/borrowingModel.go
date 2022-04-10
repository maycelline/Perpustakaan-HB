package model

import "time"

type Borrowing struct {
	ID         int       `json:"idBorrowing,omitempty"`
	Book       []Book    `json:"book,omitempty"`
	BorrowDate time.Time `json:"borrowDate,omitempty"`
	ReturnDate time.Time `json:"returnDate,omitempty"`
	Price      int       `json:"borrowPrice,omitempty"`
}
