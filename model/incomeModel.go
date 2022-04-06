package model

import "time"

type MonthIncome struct {
	Date       time.Time `json:"date"`
	SumBorrows int       `json:"sumBorrow"`
	Income     int       `json:"income"`
}

type Income struct {
	BranchName  string        `json:"branchName"`
	IncomeMonth []MonthIncome `json:"incomeMonth"`
}
