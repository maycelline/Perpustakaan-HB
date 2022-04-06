package model

import "time"

type Cart struct {
	ID         int       `json:"idCart,omitempty"`
	BorrowDate time.Time `json:"borrowDate,omitempty"`
	Book       Book      `json:"book,omitempty"`
}

type Book struct {
	ID         int    `json:"idBook,omitempty"`
	Title      string `json:"title,omitempty"`
	CoverPath  string `json:"coverPath,omitempty"`
	Author     string `json:"author,omitempty"`
	Genre      string `json:"genre,omitempty"`
	Year       int    `json:"year,omitempty"`
	Page       int    `json:"page,omitempty"`
	RentPrice  int    `json:"rentPrice,omitempty"`
	Stock      int    `json:"bookStock,omitempty"`
	BranchName string `json:"branchName,omitempty"`
}

type PopularBooksEmail struct {
	Books    []Book
	FullName string
}
