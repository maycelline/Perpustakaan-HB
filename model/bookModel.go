package model

type Book struct {
	ID         int    `json:"idBook"`
	CoverPath  string `json:"coverPath"`
	Title      string `json:"bookTitle"`
	Author     string `json:"author"`
	Genre      string `json:"genre"`
	Year       int    `json:"year"`
	Page       int    `json:"page"`
	RentPrice  int    `json:"rentPrice"`
	Stock      int    `json:"bookStock"`
	BranchName string `json:"branchName"`
}
