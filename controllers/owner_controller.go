package controllers

import (
	"Perpustakaan-HB/model"
	"log"
	"net/http"
	"strconv"
)

func GetOwnerData(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userId := getIdFromCookies(r)
	query := "SELECT fullName FROM users WHERE userId = ?"

	row := db.QueryRow(query, userId)

	var user model.User

	if err := row.Scan(&user.FullName); err != nil {
		log.Println(err)
		sendBadRequestResponse(w, "Bad Query")
	} else {
		sendSuccessResponse(w, "Success", user)
	}
}

func GetBranchIncome(w http.ResponseWriter, r *http.Request) {
	branchId, _ := strconv.Atoi(r.URL.Query().Get("branch_id"))
	incomes, err := GetIncome(branchId, w)
	if err == nil {
		sendSuccessResponse(w, "Success Get Income Data", incomes)
	}
}

func GetIncome(branchId int, w http.ResponseWriter) ([]model.MonthIncome, error) {
	db := connect()
	defer db.Close()

	query := "SELECT MONTHNAME(borrows.borrowDate), COUNT(borrowslist.borrowId), SUM(borrows.borrowPrice)FROM borrows JOIN borrowslist ON borrowslist.borrowId = borrows.borrowId JOIN stocks ON stocks.stockId = borrowslist.stockId WHERE stocks.branchId = ? GROUP BY stocks.branchId ORDER BY MONTH(borrows.borrowDate), stocks.branchId ASC"

	rows, err := db.Query(query, branchId)
	if err != nil {
		sendBadRequestResponse(w, "Error Query Get Month Income")
		return nil, err
	}

	var income model.MonthIncome
	var incomes []model.MonthIncome

	for rows.Next() {
		if err := rows.Scan(&income.MonthName, &income.SumBorrows, &income.Income); err != nil {
			sendBadRequestResponse(w, "Error Can't Fit Query Result")
			return nil, err
		} else {
			incomes = append(incomes, income)
		}
	}
	return incomes, nil
}
func GetAllIncome(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT branchId, branchName from branches"

	rows, err := db.Query(query)

	if err != nil {
		sendBadRequestResponse(w, "Failed Get Branch Data")
	}

	var branch model.Branch
	var incomesBranch []model.Income
	var incomeBranch model.Income
	var errTemp error

	for rows.Next() {
		if err := rows.Scan(&branch.ID, &branch.Name); err != nil {
			sendBadRequestResponse(w, "Failed Fit Branch Query Result")
		} else {
			incomeBranch.BranchName = branch.Name
			incomeBranch.IncomeMonth, errTemp = GetIncome(branch.ID, w)
			incomesBranch = append(incomesBranch, incomeBranch)
			if errTemp != nil {
				break
			}
		}
	}
	if errTemp == nil {
		sendSuccessResponse(w, "Success Get All Income", incomesBranch)
	}
}
