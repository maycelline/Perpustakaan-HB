package controllers

import (
	"Perpustakaan-HB/model"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	memberId := getIdFromCookies(r)

	query := "SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?"

	rows := db.QueryRow(query, memberId)

	var member model.Member
	if err := rows.Scan(&member.User.ID, &member.User.FullName, &member.User.UserName, &member.User.BirthDate, &member.User.PhoneNumber, &member.User.Email, &member.User.Address, &member.User.Password, &member.Balance); err != nil {
		sendBadRequestResponse(w, "Error Field Undefined")
		return
	} else {
		sendSuccessResponse(w, "Get Success", member)
	}

	db.Close()
}

func GetMemberCart(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	memberId := getIdFromCookies(r)

	query := "SELECT d.cartId, a.bookTitle, a.author, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE e.memberId = ? ORDER BY d.cartId ASC"

	rows, err := db.Query(query, memberId)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var cart model.Cart
	var carts []model.Cart
	for rows.Next() {
		if err := rows.Scan(&cart.ID, &cart.Book.Title, &cart.Book.Author, &cart.Book.RentPrice, &cart.Book.Stock, &cart.Book.BranchName); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			carts = append(carts, cart)
		}
	}

	if len(carts) != 0 {
		sendSuccessResponse(w, "Get Success", carts)
	} else {
		sendBadRequestResponse(w, "Error Array Size Not Correct")
	}

	db.Close()
}

func AddBookToCart(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	memberId := getIdFromCookies(r)
	booksId := r.URL.Query()["bookId"]
	branchName := r.Form.Get("branchName")

	var book model.Book
	var books []model.Book
	var stocks = make([]int, len(booksId))

	for i, bookId := range booksId {
		query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName, b.stockId FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId WHERE a.bookId = ? AND c.branchName = ?"
		row := db.QueryRow(query, bookId, branchName)
		if err := row.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName, &stocks[i]); err != nil {
			log.Println(err)
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}

		if books[i].Stock <= 0 {
			sendBadRequestResponse(w, "Error Stocks Unavailable")
			return
		}

		_, errQuery := db.Exec("INSERT INTO carts(memberId, stockId) VALUES (?, ?)", memberId, stocks[i])

		if errQuery != nil {
			log.Println(errQuery)
			sendBadRequestResponse(w, "Error Can Not Add to Cart")
			return
		}
	}
	sendSuccessResponse(w, "Add to Cart Success", books)

	db.Close()
}

func RemoveBookFromCart(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	memberId := getIdFromCookies(r)
	booksId := r.URL.Query()["bookId"]

	vars := mux.Vars(r)
	branchName := vars["branch_name"]

	var stocks = make([]int, len(booksId))

	for i, bookId := range booksId {
		query := "SELECT b.stockId FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId WHERE a.bookId = ? AND c.branchName = ?"
		row := db.QueryRow(query, bookId, branchName)
		if err := row.Scan(&stocks[i]); err != nil {
			log.Println(err)
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		}

		_, errQuery := db.Exec("DELETE FROM carts WHERE memberId=? AND stockId=?", memberId, stocks[i])

		if errQuery != nil {
			log.Println(errQuery)
			sendBadRequestResponse(w, "Error Can Not Remove from Cart")
			return
		}
	}
	sendSuccessResponse(w, "Remove from Cart Success", nil)

	db.Close()
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CheckoutBorrowing(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	memberId, _, _, _, _, _, _, _, _, balance := getDataFromCookies(r)
	booksId := r.URL.Query()["bookId"]
	deliveryDate := r.Form.Get("deliveryDate")
	deliveryTime := r.Form.Get("deliveryTime")
	borrowDate := deliveryDate + " " + deliveryTime

	var book model.Book
	var books []model.Book
	var totalBorrowPrice int
	var borrowState int
	var stocks = make([]int, len(booksId))
	var branches = make([]int, len(booksId))

	query := "SELECT COUNT(borrowState) FROM borrows JOIN borrowslist ON borrows.borrowId = borrowslist.borrowId JOIN members ON members.memberId = borrows.memberId WHERE members.memberId = ? AND borrowState = 'OVERDUE'"
	row := db.QueryRow(query, memberId)
	if err := row.Scan(&borrowState); err != nil {
		log.Println(err)
		return
	} else {
		if borrowState != 0 {
			sendBadRequestResponse(w, "Error Can Not Check Out")
			return
		}
	}

	for i, bookId := range booksId {
		query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE a.bookId = ? AND e.memberId = ?"
		row := db.QueryRow(query, bookId, memberId)
		if err := row.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}

		query2 := "SELECT b.stockId, b.branchId FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE a.bookId = ? AND e.memberId = ?"
		row2 := db.QueryRow(query2, bookId, memberId)
		if err := row2.Scan(&stocks[i], &branches[i]); err != nil {
			log.Println(err)
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		}

		totalBorrowPrice = totalBorrowPrice + books[i].RentPrice

		if books[i].Stock <= 0 {
			sendBadRequestResponse(w, "Error Stocks Unavailable")
			return
		}
	}

	if totalBorrowPrice > balance {
		sendBadRequestResponse(w, "Error Low Balance")
	}

	var branchExist []int
	var branchFound = false
	var borrowList = make(map[int]int)

	var borrowId int

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	for i := range stocks {
		branchFound = intInSlice(branches[i], branchExist)
		if !branchFound {
			branchExist = append(branchExist, branches[i])
			result, errQuery := tx.ExecContext(ctx, "INSERT INTO borrows(memberId, borrowDate, returnDate, borrowPrice) VALUES (?, ?, ?, ?)", memberId, borrowDate, "0000-00-00 00:00:00", 0)
			if errQuery != nil {
				tx.Rollback()
				log.Println(errQuery)
				sendBadRequestResponse(w, "Error Can Not Checkout")
				return
			}
			lastBorrowId, _ := result.LastInsertId()
			borrowId = int(lastBorrowId)
			borrowList[branches[i]] = borrowId
		}

		borrowId = borrowList[branches[i]]

		_, err = tx.ExecContext(ctx, "UPDATE stocks SET stock = ? WHERE stockId = ?", books[i].Stock-1, stocks[i])
		if err != nil {
			tx.Rollback()
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		}

		balance := balance - books[i].RentPrice
		_, err2 := tx.ExecContext(ctx, "UPDATE members SET balance = ? WHERE memberId = ?", balance-books[i].RentPrice, memberId)
		if err2 != nil {
			tx.Rollback()
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		}

		_, err3 := tx.ExecContext(ctx, "UPDATE borrows SET borrowPrice = borrowPrice + ? WHERE borrowId = ?", books[i].RentPrice, borrowId)
		if err3 != nil {
			tx.Rollback()
			log.Println(err3)
			sendBadRequestResponse(w, "Error Can Not Checkout")
			return
		}

		_, err4 := tx.ExecContext(ctx, "INSERT INTO borrowslist VALUES (?, ?, ?)", borrowId, stocks[i], "BORROW_PROCESS")
		if err4 != nil {
			tx.Rollback()
			log.Println(err4)
			sendBadRequestResponse(w, "Error Can Not Checkout")
			return
		}

		_, err5 := tx.ExecContext(ctx, "DELETE FROM carts WHERE memberId = ? AND stockId = ?", memberId, stocks[i])
		if err5 != nil {
			tx.Rollback()
			log.Println(err5)
			sendBadRequestResponse(w, "Error Can Not Checkout")
			return
		}
	}
	tx.Commit()
	sendSuccessResponse(w, "Checkout Success", books)
	db.Close()
}

func ReturnBorrowing(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	memberId := getIdFromCookies(r)
	booksId := r.URL.Query()["bookId"]

	var book model.Book
	var books []model.Book
	var stocks = make([]int, len(booksId))
	var borrowId = make([]int, len(booksId))

	for i, bookId := range booksId {
		query := "SELECT b.stockId, d.borrowId, a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrowslist d ON b.stockId = d.stockId JOIN borrows e ON d.borrowId = e.borrowId JOIN members f ON e.memberId = f.memberId WHERE a.bookId = ? AND e.memberId = ?"
		row := db.QueryRow(query, bookId, memberId)
		if err := row.Scan(&stocks[i], &borrowId[i], &book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
			log.Println(err)
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}

		_, err := db.Exec("UPDATE borrowslist SET borrowState = 'RETURN_PROCESS' WHERE borrowId = ? AND stockId = ?", borrowId[i], stocks[i])
		if err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		}
	}
	sendSuccessResponse(w, "Return Success", books)
	db.Close()
}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	userId, fullName, userName, birthDate, phone, email, address, additionalAddress, _, _ := getDataFromCookies(r)

	if fullName != r.Form.Get("fullName") {
		fullName = r.Form.Get("fullName")
	} else if userName != r.Form.Get("userName") {
		userName = r.Form.Get("userName")
		checkUname := checkUsernameValidation(userName, w)
		if !checkUname {
			return
		}
	} else if time.Time.String(birthDate) != r.Form.Get("birthDate") {
		birthDate, _ = time.Parse("2020-01-01", r.Form.Get("birthDate"))
	} else if phone != r.Form.Get("phone") {
		phone = r.Form.Get("phone")
	} else if email != r.Form.Get("email") {
		email = r.Form.Get("email")
		checkMail := checkMailValidation(email, w)
		if !checkMail {
			return
		}
	} else if address != r.Form.Get("address") {
		address = r.Form.Get("address")
	} else if additionalAddress != r.Form.Get("additionalAddress") {
		additionalAddress = r.Form.Get("additionalAddress")
	}

	result, errQuery := db.Exec("UPDATE users SET fullName=?, userName=?, birthDate=?, phone=?, email=?, address=?, additionalAddress=?", fullName, userName, birthDate, phone, email, address, additionalAddress)
	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userId)

	num, _ := result.RowsAffected()

	var member model.Member
	var members []model.Member
	for rows.Next() {
		if err := rows.Scan(&member.User.ID, &member.User.FullName, &member.User.UserName, &member.User.BirthDate, &member.User.PhoneNumber, &member.User.Email, &member.User.Address, &member.User.Password, &member.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			members = append(members, member)
		}
	}

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Update Success", members)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Update")
	}

	db.Close()
}

func EditUserPassword(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	userId, _, _, _, _, _, _, _, password, _ := getDataFromCookies(r)

	checkPass := checkPasswordValidation(password, w)

	if !checkPass {
		return
	}

	password = encodePassword(password)

	result, errQuery := db.Exec("UPDATE users SET password=? WHERE memberId=?", password, userId)
	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userId)

	num, _ := result.RowsAffected()

	var member model.Member
	var members []model.Member
	for rows.Next() {
		if err := rows.Scan(&member.User.ID, &member.User.FullName, &member.User.UserName, &member.User.BirthDate, &member.User.PhoneNumber, &member.User.Email, &member.User.Address, &member.User.Password, &member.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			members = append(members, member)
		}
	}

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Update Success", members)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Update")
	}

	db.Close()
}

func TopupUserBalance(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	userId, _, _, _, _, _, _, _, _, balance := getDataFromCookies(r)

	newBalance, _ := strconv.Atoi(r.Form.Get("balance"))
	balance = balance + newBalance

	result, errQuery := db.Exec("UPDATE members SET balance=? WHERE memberId=?", balance, userId)
	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userId)

	num, _ := result.RowsAffected()

	var member model.Member

	err = rows.Scan(&member.User.ID, &member.User.FullName, &member.User.UserName, &member.User.BirthDate, &member.User.PhoneNumber, &member.User.Email, &member.User.Address, &member.User.Password, &member.Balance)
	if err != nil {
		sendBadRequestResponse(w, "Error Field Undefined")
		return
	}

	if errQuery == nil {
		if num == 0 {
			sendBadRequestResponse(w, "Error 0 Rows Affected")
		} else {
			sendSuccessResponse(w, "Update Success", member)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Update")
	}

	db.Close()
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	userId := getIdFromCookies(r)

	rows, _ := db.Query("SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, password, balance FROM users JOIN members ON users.userId = members.memberId WHERE users.userId=?", userId)

	var member model.Member
	var members []model.Member
	for rows.Next() {
		if err := rows.Scan(&member.User.ID, &member.User.FullName, &member.User.UserName, &member.User.BirthDate, &member.User.PhoneNumber, &member.User.Email, &member.User.Address, &member.User.Password, &member.Balance); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			members = append(members, member)
		}
	}

	_, errQuery := db.Exec("DELETE FROM members WHERE userId=?", userId)
	if errQuery == nil {
		result, errQuery := db.Exec("DELETE FROM users WHERE memberId=?", userId)
		num, _ := result.RowsAffected()
		if errQuery == nil {
			if num == 0 {
				sendBadRequestResponse(w, "Error 0 Rows Affected")
			} else {
				resetUserToken(w)
				sendSuccessResponse(w, "Delete Success", members)
			}
		} else {
			sendSuccessResponse(w, "Delete Success", members)
		}
	} else {
		sendBadRequestResponse(w, "Error Can Not Delete")
	}

	db.Close()
}

func GetMemberHistory(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	memberId := getIdFromCookies(r)

	query := "SELECT e.borrowId, e.borrowDate, e.returnDate FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrowslist d ON b.stockId = d.stockId JOIN borrows e ON d.borrowId = e.borrowId JOIN members f ON e.memberId = f.memberId WHERE d.borrowState = 'BORROWED' OR d.borrowState = 'OVERDUE' OR d.borrowState = 'RETURNED' AND  f.memberId = ?"

	rows, err := db.Query(query, memberId)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var borrowing model.Borrowing
	var borrowings []model.Borrowing
	for rows.Next() {
		if err := rows.Scan(&borrowing.ID, &borrowing.BorrowDate, &borrowing.ReturnDate); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			query2 := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrowslist d ON b.stockId = d.stockId JOIN borrows e ON d.borrowId = e.borrowId JOIN members f ON e.memberId = f.memberId WHERE d.borrowState = 'BORROWED' OR d.borrowState = 'OVERDUE' OR d.borrowState = 'FINISHED' AND  f.memberId = ?"

			rows2, err := db.Query(query2, memberId)
			if err != nil {
				sendNotFoundResponse(w, "Table Not Found")
				return
			}

			var book model.Book
			var books []model.Book
			for rows2.Next() {
				if err := rows2.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
					sendBadRequestResponse(w, "Error Field Undefined")
					return
				} else {
					books = append(books, book)
				}
			}
			borrowing.Book = books
			borrowings = append(borrowings, borrowing)
		}
	}

	if len(borrowings) != 0 {
		sendSuccessResponse(w, "Get Success", borrowings)
	} else {
		sendBadRequestResponse(w, "You have no borrowing history")
	}

	db.Close()
}
