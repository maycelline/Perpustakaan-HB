package model

type Book struct {
	ID         int    `json:"idBook"`
	CoverPath  string `json:"coverPath"`
	Author     string `json:"author"`
	Genre      string `json:"genre"`
	Year       int    `json:"year"`
	Page       int    `json:"page"`
	RentPrice  int    `json:"rentPrice"`
	BookStock  int    `json:"bookStock"`
	BranchName string `json:"branchName"`
}
