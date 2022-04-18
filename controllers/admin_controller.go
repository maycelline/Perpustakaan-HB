package controllers

import (
	"Perpustakaan-HB/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	// _ "github.com/lib/pq"
)

func GetAdminData(w http.ResponseWriter, r *http.Request) {
	db := Connect()
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
	db := Connect()
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

	queryBorrow := "SELECT borrows.borrowId, borrows.returnDate FROM borrows JOIN borrowslist ON borrows.borrowId = borrowsList.borrowId WHERE borrowState = 'BORROW_PROCESS'"
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
	db := Connect()
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

	queryBorrow := "SELECT borrows.borrowId, borrows.returnDate FROM borrows JOIN borrowslist ON borrows.borrowId = borrowslist.borrowId WHERE borrowState = 'RETURN_PROCESS'"
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
	// fmt.Print("masok")
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	stateType := r.Form.Get("stateType")
	couriedId := r.Form.Get("courierId")
	stockId := r.Form.Get("stockId")
	vars := mux.Vars(r)
	borrowId := vars["borrow_id"]
	stockIds := strings.Split(stockId, ",")
	deliveryFee := r.Form.Get("deliveryFee")
	// fmt.Println("Deliv Fee: " + deliveryFee)

	fmt.Println(stateType)
	fmt.Println(couriedId)

	// state := true

	var count int64 = 0

	if len(stateType) <= 0 || len(couriedId) <= 0 {
		sendBadRequestResponse(w, "Please input all fields")
		return
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	result, _ := tx.ExecContext(ctx, "UPDATE borrows SET borrowPrice = borrowPrice + ?", deliveryFee)

	for i := 0; i < len(stockIds); i++ {

		result, err = tx.ExecContext(ctx, "UPDATE borrowslist SET borrowState = ? WHERE borrowId=? AND stockId=?", stateType, borrowId, stockIds[i])
		if err != nil {
			tx.Rollback()
			fmt.Println(err)
			sendBadRequestResponse(w, "Error Can Not Change")
			return
		} else {
			num, _ := result.RowsAffected()
			if num != 0 {
				if stateType == "RETURNED" {
					_, err = tx.ExecContext(ctx, "UPDATE stocks SET stock = stock + 1 WHERE stockId=?", stockIds[i])
					if err != nil {
						tx.Rollback()
						log.Fatal(err)
						sendBadRequestResponse(w, "Error: stock cannot be updated")
						return
					}
				}
				count += num
			}
		}
	}

	if count == int64(len(stockIds)) {
		tx.Commit()
		sendSuccessResponse(w, "State changed to "+strings.ToLower(stateType), nil)
		if stateType == "RETURNED" {
			//Get User Balance
			// var user model.User
			var balance int
			var memberId int
			query := "SELECT members.balance, members.memberId from borrows JOIN members ON borrows.memberId = members.memberId WHERE borrowId = ?"

			row := db.QueryRow(query, borrowId)
			if err := row.Scan(&balance, &memberId); err != nil {
				fmt.Println("Balance error: ")
				fmt.Println(err)
				return
			} else {
				_, err := db.Exec("UPDATE members SET balance = balance - ? WHERE memberId = ?", deliveryFee, memberId)

				if err != nil {
					fmt.Println("Update balance error: ")
					fmt.Println(err)
					return
				}
			}

			fmt.Println(balance)
		}
		getAllDataForTransactionEmail(borrowId, stateType, couriedId, stockIds)
	} else {
		tx.Rollback()
		sendServerErrorResponse(w, "Your stock id is not in the list ")
	}

}

func AddNewBook(w http.ResponseWriter, r *http.Request) {

	db := Connect()
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

func getAllDataForTransactionEmail(borrowId string, borrowState string, courierId string, booksId []string) {
	fmt.Println("This function is for send email")
	db := Connect()
	defer db.Close()

	//Get Books Data
	var book model.Book
	var books []model.Book

	for i := 0; i < len(booksId); i++ {
		query := "SELECT books.bookTitle, books.author FROM stocks JOIN books ON books.bookId = stocks.bookId WHERE stocks.stockId = ?"
		row := db.QueryRow(query, booksId[i])
		if err := row.Scan(&book.Title, &book.Author); err != nil {
			fmt.Println("Book error: ")
			fmt.Println(err)
		} else {
			books = append(books, book)
		}
	}

	//Get User Name and Email
	var user model.User
	query := "SELECT users.fullName,  users.email FROM borrows JOIN users ON borrows.memberId = users.userId WHERE borrows.borrowId = ?"
	row := db.QueryRow(query, borrowId)
	if err := row.Scan(&user.FullName, &user.Email); err != nil {
		fmt.Println("user error: ")
		fmt.Println(err)
	}

	//Get Courier
	var courier model.Courier
	query = "SELECT courierName FROM couriers WHERE courierId = ?"
	row = db.QueryRow(query, courierId)
	if err := row.Scan(&courier.CourierName); err != nil {
		fmt.Println("courier error: ")
		fmt.Println(err)
	}

	//Get Branch
	var branch model.Branch
	query = "SELECT branches.branchName, branches.branchAddress FROM borrows JOIN borrowslist ON borrows.borrowId = borrowslist.borrowId JOIN stocks ON borrowslist.stockId = stocks.stockId JOIN branches ON stocks.branchId = branches.branchId WHERE borrows.borrowId = ?"
	row = db.QueryRow(query, borrowId)
	if err := row.Scan(&branch.Name, &branch.Address); err != nil {
		fmt.Println("branch error: ")
		fmt.Println(err)
	}

	//Prepare data for email
	var data model.BorrowDataHTML
	data.Books = books
	data.Branch = branch
	data.Courier = courier
	data.User = user

	date := time.Now()
	y, m, d := date.Date()

	data.CourierCome = strconv.Itoa(d)
	data.CourierCome += " "
	data.CourierCome += m.String()
	data.CourierCome += " "
	data.CourierCome += strconv.Itoa(y)

	if borrowState == "BORROWED" {
		SendBorrowAcceptedEmail(data)
	} else {
		SendReturnAcceptedEmail(data)
	}

}
