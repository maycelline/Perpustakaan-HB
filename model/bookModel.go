package model

type Book struct {
	ID         int    `json:"idBook"`
	Title      string `json:"title"`
	CoverPath  string `json:"coverPath"`
	Author     string `json:"author"`
	Genre      string `json:"genre"`
	Year       int    `json:"year"`
	Page       int    `json:"page"`
	RentPrice  int    `json:"rentPrice"`
	Stock      int    `json:"bookStock"`
	BranchName string `json:"branchName"`
}

type PopularBook struct {
	ID        int    `json:"idBook"`
	Title     string `json:"title"`
	CoverPath string `json:"coverPath"`
	Author    string `json:"author"`
	Genre     string `json:"genre"`
	Year      int    `json:"year"`
	Page      int    `json:"page"`
}
