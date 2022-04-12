package controllers

import (
	"Perpustakaan-HB/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetOwnerData(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userId := r.URL.Query().Get("user_id")
	query := "SELECT fullName FROM users WHERE userId = " + userId + "; "

	row := db.QueryRow(query)

	var user model.User

	if err := row.Scan(&user.FullName); err != nil {
		log.Println(err)
		sendBadRequestResponse(w, "Bad Query")
	} else {
		sendSuccessResponse(w, "Success", user)
	}
}

func GetBranchIncome() /*w http.ResponseWriter, r *http.Request*/ {
	fmt.Println("1")
	db := connect()
	defer db.Close()
	fmt.Println("2")
	// branchId := r.URL.Query().Get("branch_id")
	query := "SELECT borrows.borrowDate, borrows.borrowPrice FROM borrows JOIN borrowslist ON borrowslist.borrowId = borrows.borrowId JOIN stocks ON stock.stockId = borrowslist.stockId WHERE stocks.branchId = 1 ORDER BY borrows.borrowDate ASC; "

	rows, err := db.Query(query)
	if err != nil {

	}

	var borrow model.Borrowing
	var borrows []model.Borrowing

	for rows.Next() {
		if err := rows.Scan(&borrow.BorrowDate, &borrow.Price); err != nil {
			// log.Println(err)
			// sendBadRequestResponse(w, "Bad Query")
		} else {
			borrows = append(borrows, borrow)
		}
	}

	var month []time.Month
	for i := 0; i < len(borrows); i++ {
		checkMonth := 0
		tempMonth := borrows[i].BorrowDate.Month()
		for j := 0; j <= len(month)-1; j++ {
			if month[j] == tempMonth {
				break
			} else {
				checkMonth += 1
			}
		}
		if checkMonth == len(month)-1 {
			month[checkMonth+1] = tempMonth
		}
	}
	fmt.Println(month)
	// var month model.MonthIncome
	// var arrIncome []model.MonthIncome
	// var income model.Income

	// sendSuccessResponse(w, "Success", user)
}

func GetAllIncome(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()

	month, _ := strconv.Atoi(r.Form.Get("month"))

	if err != nil {

		return
	}

	query := "Select sum(borrowPrice) as Income from borrows where extract(MONTH from borrowDate) = ? and borrowState = 'FINISHED'"

	rows, errQuery := db.Query(query, month)

	if errQuery != nil {
		return
	} else {
		rows.Scan()
	}

}
