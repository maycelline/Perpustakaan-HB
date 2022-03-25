package model

import "time"

type MonthIncome struct {
	Date       time.Time `json:"date"`
	SumBorrows int       `json:"sumBorrow"`
	Income     int       `json:"income"`
}

type Income struct {
	Branch      Branch        `json:"branch"`
	IncomeMonth []MonthIncome `json:"incomeMonth"`
}
