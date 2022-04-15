package controllers

import (
	"Perpustakaan-HB/model"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetUserData(w http.ResponseWriter, r *http.Request) {
	db := connect()
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
	db := connect()
	defer db.Close()

	memberId := getIdFromCookies(r)

	query := "SELECT d.cartId, a.bookTitle, a.author, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE e.memberId = ?"

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
	db := connect()
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
		query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE a.bookId = ? AND e.memberId = ?"
		rows, err := db.Query(query, bookId, memberId)
		if err != nil {
			sendNotFoundResponse(w, "Table Not Found")
			return
		}
		for rows.Next() {
			if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
				sendBadRequestResponse(w, "Error Field Undefined")
				return
			} else {
				books = append(books, book)
			}
		}
		query2 := "SELECT b.stockId FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId WHERE a.bookId = ? AND c.branchName = ?"
		rows2 := db.QueryRow(query2, bookId, branchName)
		if err := rows2.Scan(&stocks[i]); err != nil {
			log.Println(err)
			sendBadRequestResponse(w, "Error Field Undefined")
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

func CreateBorrowingList(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	memberId := getIdFromCookies(r)
	booksId := r.URL.Query()["bookId"]
	deliveryDate := r.Form.Get("deliveryDate")
	deliveryTime := r.Form.Get("deliveryTime")
	borrowDate := deliveryDate + " " + deliveryTime

	var book model.Book
	var books []model.Book
	var stock int
	var stocks = make([]int, len(booksId))
	var borrowPrice = make([]int, len(booksId))
	var borrowId int

	for i, bookId := range booksId {
		// _, err := db.Exec("START TRANSACTION; SAVEPOINT beforeCheckpoint;")
		if err == nil {
			query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE a.bookId = ? AND e.memberId = ?"
			rows, err := db.Query(query, bookId, memberId)
			if err != nil {
				sendNotFoundResponse(w, "Table Not Found")
				return
			}
			for rows.Next() {
				if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock, &book.BranchName); err != nil {
					sendBadRequestResponse(w, "Error Field Undefined")
					return
				} else {
					books = append(books, book)
				}
			}
			query2 := "SELECT b.stock, b.stockId, a.rentPrice FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN carts d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE a.bookId = ? AND e.memberId = ?"
			log.Println(query2)
			rows2 := db.QueryRow(query2, bookId, memberId)
			if err := rows2.Scan(&stock, &stocks[i], &borrowPrice[i]); err != nil {
				log.Println(err)
				sendBadRequestResponse(w, "Error Field Undefined")
				return
			}

			log.Println("Stock: ", stock)
			if stock <= 0 {
				sendBadRequestResponse(w, "Error Stocks Unavailable")
				return
			}

			// BELUM BIKIN KALO MAU MINJEMNYA LEBIH DARI 1 gtw deng mikir dl

			_, err = db.Exec("UPDATE stocks SET stock = ? WHERE stockId = ?", stock-1, stocks[i])
			if err != nil {
				// db.Exec("ROLLBACK TO beforeCheckpoint;")
				sendBadRequestResponse(w, "Error Field Undefined")
				return
			}

			if i == 0 {
				result, errQuery := db.Exec("INSERT INTO borrows(memberId, borrowDate, returnDate, borrowPrice) VALUES (?, ?, ?, ?)", memberId, borrowDate, "0000-00-00 00:00:00", 0)
				if errQuery != nil {
					// db.Exec("ROLLBACK TO beforeCheckpoint;")
					log.Println(errQuery)
					sendBadRequestResponse(w, "Error Can Not Checkout")
					return
				}
				lastBorrowId, _ := result.LastInsertId()
				borrowId = int(lastBorrowId)
			}

			_, errQuery := db.Exec("UPDATE borrows SET borrowPrice = borrowPrice + ? WHERE borrowId = ?", borrowPrice[i], borrowId)
			if errQuery != nil {
				// db.Exec("ROLLBACK TO beforeCheckpoint;")
				log.Println(errQuery)
				sendBadRequestResponse(w, "Error Can Not Checkout")
				return
			} else {
				_, errQuery := db.Exec("INSERT INTO borrowslist VALUES (?, ?, ?)", borrowId, stocks[i], "BORROW_PROCESS")
				if errQuery != nil {
					// db.Exec("ROLLBACK TO beforeCheckpoint;")
					log.Println(errQuery)
					sendBadRequestResponse(w, "Error Can Not Checkout")
					return
				} else {
					_, errQuery := db.Exec("DELETE FROM carts WHERE memberId = ? AND stockId = ?", memberId, stocks[i])
					if errQuery != nil {
						// db.Exec("ROLLBACK TO beforeCheckpoint;")
						log.Println(errQuery)
						sendBadRequestResponse(w, "Error Can Not Checkout")
						return
					}
				}
			}
		}

	}
	sendSuccessResponse(w, "Checkout Success", books)

	db.Close()
}

func GetOngoingBorrowing(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	memberId := getIdFromCookies(r)

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrows d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE d.borrowState = 'ON-GOING' AND  e.memberId = ?"

	rows, err := db.Query(query, memberId)
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
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	userId, fullName, userName, birthDate, phone, email, address, additionalAddress, _, _ := getDataFromCookies(r)

	if fullName != r.Form.Get("fullName") {
		fullName = r.Form.Get("fullName")
	} else if userName != r.Form.Get("userName") {
		userName = r.Form.Get("userName")
	} else if time.Time.String(birthDate) != r.Form.Get("birthDate") {
		birthDate, _ = time.Parse("2020-01-01", r.Form.Get("birthDate"))
	} else if phone != r.Form.Get("phone") {
		phone = r.Form.Get("phone")
	} else if email != r.Form.Get("email") {
		email = r.Form.Get("email")
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
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	userId, _, _, _, _, _, _, _, password, _ := getDataFromCookies(r)

	passwordLength := len(password)

	if passwordLength < 8 {
		sendBadRequestResponse(w, "Need more character")
		return
	} else if passwordLength > 10 {
		sendBadRequestResponse(w, "Too many character")
		return
	}

	containsNumber := 0
	for i := 0; i < 10; i++ {
		number := strconv.Itoa(i)
		if strings.Contains(password, number) {
			containsNumber = containsNumber + 1
		}
	}

	passwordCheck := strings.ToLower(password)
	arrayPassword := []rune(passwordCheck)

	containsLowerCase := 0
	for i := 0; i < passwordLength; i++ {
		char := string(arrayPassword)
		if strings.Contains(password, char) {
			containsLowerCase = containsLowerCase + 1
		}
	}

	if containsNumber == 0 || containsLowerCase == 0 {
		sendBadRequestResponse(w, "Bad password")
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
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	userId, _, _, _, _, _, _, _, _, balance := getDataFromCookies(r)

	newBalance, _ := strconv.Atoi(r.Form.Get("balance"))
	balance = balance + newBalance

	result, errQuery := db.Exec("UPDATE users SET balance=? WHERE memberId=?", balance, userId)
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

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	db := connect()
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

	_, errQuery := db.Exec("DELETE FROM members WHERE id=?", userId)
	if errQuery == nil {
		result, errQuery := db.Exec("DELETE FROM users WHERE id=?", userId)
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
	db := connect()
	defer db.Close()

	memberId := getIdFromCookies(r)

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c ON b.branchId = c.branchId JOIN borrows d ON b.stockId = d.stockId JOIN members e ON d.memberId = e.memberId WHERE d.borrowState = 'ON-GOING' OR d.BorrowState = 'FINISHED' AND  e.memberId = "

	rows, err := db.Query(query, memberId)
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
