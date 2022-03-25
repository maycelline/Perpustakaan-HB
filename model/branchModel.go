package model

type Branch struct {
	ID      int    `json:"branchId"`
	Name    string `json:"branchName"`
	Address string `json:"branchAddress"`
}
