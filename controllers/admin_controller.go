package controllers

import (
	"Perpustakaan-HB/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetAdminData(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	adminId := getIdFromCookies(r)
	query := "SELECT u.FullName, b.branchId, b.branchName ,b.branchAddress FROM users u JOIN admins a ON u.userId = a.adminId JOIN branches b ON a.branchId = b.branchId WHERE a.adminId = " + strconv.Itoa(adminId) + "; "

	row := db.QueryRow(query)

	var user model.User
	var branch model.Branch

	if err := row.Scan(&user.FullName, &branch.ID, &branch.Name, &branch.Address); err != nil {
		log.Println(err)
		sendBadRequestResponse(w, "Error Field Undefined")
	} else {
		var admin model.Admin
		admin.User = user
		admin.Branch = branch
		sendSuccessResponse(w, "Get Success", admin)
	}
}

func GetUnapprovedBorrowing(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	adminId := getIdFromCookies(r)
	var branchId int

	//Get Admin Branch
	query := "SELECT branchId FROM admins WHERE adminId = " + strconv.Itoa(adminId) + "; "
	row := db.QueryRow(query)

	if err := row.Scan(&branchId); err != nil {
		log.Println(err)
		sendBadRequestResponse(w, "Error Field Undefined")
		return
	}

	fmt.Println(branchId, "branch")

	queryBorrow := "SELECT borrowId, returnDate FROM borrows WHERE borrowState = 'BORROW_PROCESS'"
	rowsBorrow, err := db.Query(queryBorrow)

	if err != nil {
		fmt.Println(err)
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var borrowing model.Borrowing
	var borrowings []model.Borrowing
	var courier model.Courier
	var couriers []model.Courier
	for rowsBorrow.Next() {
		if err := rowsBorrow.Scan(&borrowing.ID, &borrowing.ReturnDate); err != nil {
			fmt.Println(err)
			sendBadRequestResponse(w, "Error Field Undifined")
			return
		} else {
			borrowings = append(borrowings, borrowing)
		}
	}

	queryCourier := "SELECT courierId, courierName FROM couriers WHERE courierState = 'AVAILABLE'"
	rowsCourier, err := db.Query(queryCourier)

	if err != nil {
		fmt.Println(err)
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	for rowsCourier.Next() {
		if err := rowsCourier.Scan(&courier.ID, &courier.CourierName); err != nil {
			fmt.Println(err)
			sendBadRequestResponse(w, "Error Field Undifined")
			return
		} else {
			couriers = append(couriers, courier)
		}
	}

	var borrowData model.BorrowData
	borrowData.Borrows = borrowings
	borrowData.Couriers = couriers

	sendSuccessResponse(w, "Approve Success", borrowData)

}

func GetUnapprovedReturn(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	adminId := getIdFromCookies(r)
	var branchId int

	//Get Admin Branch -- Rencananya mau ditampilin per cabang
	query := "SELECT branchId FROM admins WHERE adminId = " + strconv.Itoa(adminId) + "; "
	row := db.QueryRow(query)

	if err := row.Scan(&branchId); err != nil {
		log.Println(err)
		sendBadRequestResponse(w, "Bad Query")
		return
	}

	queryBorrow := "SELECT borrowId, returnDate FROM borrows WHERE borrowState = 'RETURN_PROCESS'"
	rowsBorrow, err := db.Query(queryBorrow)

	if err != nil {
		fmt.Println(err)
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var borrowing model.Borrowing
	var borrowings []model.Borrowing
	var courier model.Courier
	var couriers []model.Courier
	for rowsBorrow.Next() {
		if err := rowsBorrow.Scan(&borrowing.ID, &borrowing.ReturnDate); err != nil {
			fmt.Println(err)
			sendBadRequestResponse(w, "Error Field Undifined")
			return
		} else {
			borrowings = append(borrowings, borrowing)
		}
	}

	queryCourier := "SELECT courierId, courierName FROM couriers WHERE courierState = 'AVAILABLE'"
	rowsCourier, err := db.Query(queryCourier)

	if err != nil {
		fmt.Println(err)
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	for rowsCourier.Next() {
		if err := rowsCourier.Scan(&courier.ID, &courier.CourierName); err != nil {
			fmt.Println(err)
			sendBadRequestResponse(w, "Error Field Undifined")
			return
		} else {
			couriers = append(couriers, courier)
		}
	}

	var borrowData model.BorrowData
	borrowData.Borrows = borrowings
	borrowData.Couriers = couriers

	sendSuccessResponse(w, "Approve Success", borrowData)
}

func ChangeBorrowingState(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	stateType := r.Form.Get("stateType")
	couriedId := r.Form.Get("courierId")
	vars := mux.Vars(r)
	borrowId := vars["borrow_id"]

	if len(stateType) <= 0 || len(couriedId) <= 0 {
		sendBadRequestResponse(w, "Value Not Found")
		return
	}

	result, err := db.Exec("UPDATE borrows SET borrowState = ? WHERE borrowId=?", stateType, borrowId)

	if err != nil {
		sendBadRequestResponse(w, "Error Can Not Change")
		return
	} else {
		num, _ := result.RowsAffected()
		if num == 0 {
			sendServerErrorResponse(w, "This order's state already in "+strings.ToLower(stateType))
		} else {
			sendSuccessResponse(w, "State changed to "+strings.ToLower(stateType), nil)
			getAllDataForTransactionEmail(borrowId, stateType, couriedId)
		}
	}
}

func AddNewBook(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	title := r.Form.Get("title")
	coverPath := r.Form.Get("coverPath")
	author := r.Form.Get("author")
	genre := r.Form.Get("genre")
	year, _ := strconv.Atoi(r.Form.Get("year"))
	page, _ := strconv.Atoi(r.Form.Get("page"))
	rentPrice, _ := strconv.Atoi(r.Form.Get("rentPrice"))
	branchId, _ := strconv.Atoi(r.Form.Get("branchId"))
	stock, _ := strconv.Atoi(r.Form.Get("stock"))

	query1 := "INSERT INTO books(bookTitle, author, genre, year, page, rentPrice, coverPath) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, errQuery1 := db.Exec(query1, title, author, genre, year, page, rentPrice, coverPath)

	if errQuery1 != nil {
		sendBadRequestResponse(w, "Error cannot add new book")
		return
	}

	bookId, _ := result.LastInsertId()
	query2 := "Insert into stocks(bookId, brandId, stock) values(?,?,?)"

	_, errQuery2 := db.Exec(query2, bookId, branchId, stock)

	if errQuery2 != nil {
		sendBadRequestResponse(w, "Error cannot add new book")
		return
	}

	sendSuccessResponse(w, "Add book successfully", nil)
}

func getAllDataForTransactionEmail(borrowId string, borrowState string, courierId string) {
	db := connect()
	defer db.Close()

	queryBooks := "SELECT books.bookTitle, books.author FROM borrowslist JOIN stocks ON borrowslist.stockId = stocks.stockId JOIN books ON books.bookId = stocks.bookId WHERE borrowslist.borrowId = ?"

	rowsBook, err := db.Query(queryBooks, borrowId)

	if err != nil {
		fmt.Println(err)
		return
	}

	var book model.Book
	var books []model.Book

	for rowsBook.Next() {
		if err := rowsBook.Scan(&book.Title, &book.Author); err != nil {
			log.Println(err)
			return
		} else {
			books = append(books, book)
		}
	}

	fmt.Println(books)

	// var member model.Member
	// var courier model.Courier

}
