package model

type Book struct {
	ID         int    `json:"idBook"`
	Title      string `json:"title"`
	CoverPath  string `json:"coverPath"`
	Author     string `json:"author"`
	Genre      string `json:"genre"`
	Year       int    `json:"year"`
	Page       int    `json:"page,omitempty"`
	RentPrice  int    `json:"rentPrice,omitempty"`
	Stock      int    `json:"bookStock,omitempty"`
	BranchName string `json:"branchName,omitempty"`
}

// type PopularBook struct {
// 	ID        int    `json:"idBook"`
// 	Title     string `json:"title"`
// 	CoverPath string `json:"coverPath"`
// 	Author    string `json:"author"`
// 	Genre     string `json:"genre"`
// 	Year      int    `json:"year"`
// }

type PopularBooksEmail struct {
	Books    []Book
	FullName string
}
