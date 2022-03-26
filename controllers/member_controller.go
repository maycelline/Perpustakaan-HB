package controllers

import (
	"net/http"
	"strconv"

	"Perpustakaan-HB/model"

	"github.com/gorilla/mux"
)

func GetAUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?"

	// nanti get cookies memberId?
	rows := db.QueryRow(query, 6)

	var user model.User
	if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.Password, &user.Balance); err != nil {
		sendBadRequestResponse(w, "Error Field Undefined")
		return
	} else {
		sendSuccessResponse(w, "Get Success", user)
	}

	db.Close()
}

func GetMemberCart(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	memberId := strconv.Itoa(6)

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE e.memberId = " + memberId

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var book model.Book
	var books []model.Book
	for rows.Next() {
		if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}
	}

	if len(books) != 0 {
		sendSuccessResponse(w, "Get Success", books)
	} else {
		sendBadRequestResponse(w, "Error Array Size Not Correct")
	}

	db.Close()
}

func CreateBorrowingList(w http.ResponseWriter, r *http.Request) {
	// insert new to borrows
	// delete cart yang udah ga kepake
	// return
}

func GetOngoingBorrowing(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	memberId := strconv.Itoa(6)

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrows d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE WHERE d.borrowState = 'ON-GOING' AND  e.memberId = " + memberId

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var book model.Book
	var books []model.Book
	for rows.Next() {
		if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}
	}

	if len(books) != 0 {
		sendSuccessResponse(w, "Get Success", books)
	} else {
		sendBadRequestResponse(w, "Error Array Size Not Correct")
	}

	db.Close()
}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendNotFoundResponse(w, "Value Not Found")
		return
	}

	// dari cookies
	// userID := 6
	fullName := "Ariesta Leevine"
	userName := "ariestacsleevine"
	birthDate := "2020-04-09"
	phone := "0898-1234-1234"
	email := "ariestacsleevine"
	address := "Jalan Peta No. 241"

	vars := mux.Vars(r)
	userID := vars["member_id"]

	if fullName != r.Form.Get("fullName") {
		fullName = r.Form.Get("fullName")
	} else if userName != r.Form.Get("userName") {
		userName = r.Form.Get("userName")
	} else if birthDate != r.Form.Get("birthDate") {
		birthDate = r.Form.Get("birthDate")
	} else if phone != r.Form.Get("phone") {
		phone = r.Form.Get("phone")
	} else if email != r.Form.Get("email") {
		email = r.Form.Get("email")
	} else if address != r.Form.Get("address") {
		address = r.Form.Get("address")
	}

	result, errQuery := db.Exec("UPDATE users SET fullName=?, userName=?, birthDate=?, phone=?, email=?, address=?", fullName, userName, birthDate, phone, email, address)
	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userID)

	num, _ := result.RowsAffected()

	var user model.User
	var users []model.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.Password, &user.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			users = append(users, user)
		}
	}

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Update Success", users)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Update")
	}

	db.Close()
}

func EditUserPassword(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendNotFoundResponse(w, "Value Not Found")
		return
	}

	// dari cookies
	// userID := 6

	vars := mux.Vars(r)
	userID := vars["member_id"]
	password := r.Form.Get("password")

	result, errQuery := db.Exec("UPDATE users SET password=? WHERE memberId=?", password, userID)
	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userID)

	num, _ := result.RowsAffected()

	var user model.User
	var users []model.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.Password, &user.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			users = append(users, user)
		}
	}

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Update Success", users)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Update")
	}

	db.Close()
}

func TopupUserBalance(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendNotFoundResponse(w, "Value Not Found")
		return
	}

	// dari cookies
	// userID := 6
	//balance get dari cookies

	vars := mux.Vars(r)
	userID := vars["member_id"]
	newBalance, _ := strconv.Atoi(r.Form.Get("balance"))
	balance := 50000 + newBalance

	result, errQuery := db.Exec("UPDATE users SET balance=? WHERE memberId=?", balance, userID)
	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userID)

	num, _ := result.RowsAffected()

	var user model.User
	var users []model.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.Password, &user.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			users = append(users, user)
		}
	}

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Update Success", users)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Update")
	}

	db.Close()
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	vars := mux.Vars(r)
	userID := vars["userID"]

	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userID)

	var user model.User
	var users []model.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.Password, &user.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			users = append(users, user)
		}
	}

	result, errQuery := db.Exec("DELETE FROM users WHERE id=?", userID)

	num, _ := result.RowsAffected()

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Delete Success", users)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Delete")
	}

	db.Close()
}

func GetMemberHistory(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	memberId := strconv.Itoa(6)

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrows d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE d.borrowState = 'ON-GOING' OR d.BorrowState = 'FINISHED' AND  e.memberId = " + memberId

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var book model.Book
	var books []model.Book
	for rows.Next() {
		if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}
	}

	if len(books) != 0 {
		sendSuccessResponse(w, "Get Success", books)
	} else {
		sendBadRequestResponse(w, "Error Array Size Not Correct")
	}

	db.Close()
}
