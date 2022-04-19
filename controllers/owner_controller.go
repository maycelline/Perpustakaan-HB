package controllers

import (
	"Perpustakaan-HB/model"
	"net/http"
	"strconv"
)

func GetOwnerData(w http.ResponseWriter, r *http.Request) {
	// db := connectGorm()
	// defer db.Close()
	var user model.User
	user.ID, user.FullName, user.UserName, user.BirthDate, user.PhoneNumber, user.UserType, _, _, _, _ = getDataFromCookies(r)
	// row := db.Table("users").Select("fullName").Where("userId = ?", userId).Row()

	// var user model.User

	// if err := row.Scan(&user.FullName, &user.BirthDate, &user.PhoneNumber, &user.); err != nil {
	// 	log.Println(err)
	// 	sendBadRequestResponse(w, "Bad Query")
	// } else {
	sendSuccessResponse(w, "Success", user)
	// }
}

func GetBranchIncome(w http.ResponseWriter, r *http.Request) {
	branchId, _ := strconv.Atoi(r.URL.Query().Get("branch_id"))
	incomes, err := GetIncome(branchId, w)
	if err == nil {
		sendSuccessResponse(w, "Success Get Income Data", incomes)
	}
}

func GetIncome(branchId int, w http.ResponseWriter) ([]model.MonthIncome, error) {
	db := Connect()
	defer db.Close()

	query := "SELECT MONTHNAME(borrows.borrowDate), COUNT(borrowslist.borrowId), SUM(DISTINCT(borrows.borrowPrice)) FROM borrows JOIN borrowslist ON borrowslist.borrowId = borrows.borrowId JOIN stocks ON stocks.stockId = borrowslist.stockId WHERE stocks.branchId = ? GROUP BY borrowslist.borrowId, stocks.branchId ORDER BY MONTH(borrows.borrowDate), stocks.branchId ASC"

	rows, err := db.Query(query, branchId)
	if err != nil {
		sendBadRequestResponse(w, "Error Query Get Month Income")
		return nil, err
	}

	var income model.MonthIncome
	var incomes []model.MonthIncome
	var tempMonth string
	tempMonth = ""
	i := -1

	for rows.Next() {
		if err := rows.Scan(&income.MonthName, &income.SumBorrows, &income.Income); err != nil {
			sendBadRequestResponse(w, "Error Can't Fit Query Result")
			return nil, err
		} else {
			if i != -1 {
				if income.MonthName == tempMonth {
					incomes[i].Income = incomes[i].Income + income.Income
					incomes[i].SumBorrows = incomes[i].SumBorrows + income.SumBorrows
				} else {
					incomes = append(incomes, income)
					tempMonth = income.MonthName
					i++
				}
			} else {
				incomes = append(incomes, income)
				tempMonth = income.MonthName
				i++
			}
		}
	}
	return incomes, nil
}
func GetAllIncome(w http.ResponseWriter, r *http.Request) {
	db := connectGorm()

	rows, err := db.Table("branches").Select("branchId", "branchName").Rows()

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
			return
		} else {
			incomeBranch.BranchName = branch.Name
			incomeBranch.IncomeMonth, errTemp = GetIncome(branch.ID, w)
			incomesBranch = append(incomesBranch, incomeBranch)
			if errTemp != nil {
				sendBadRequestResponse(w, "Failed Fit Data Income")
				return
			}
		}
	}
	if errTemp == nil {
		sendSuccessResponse(w, "Success Get All Income", incomesBranch)
	}
}
