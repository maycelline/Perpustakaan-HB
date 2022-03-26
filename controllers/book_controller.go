package controllers

import (
	// "crypto/md5"
	// "encoding/hex"
	model "Tools/model"
	"log"
	"net/http"
	"strconv"
	// "github.com/gorilla/mux"
)

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c WHERE b.branchId = c.branchId, "

	bookId, _ := strconv.Atoi(r.URL.Query().Get("bookId"))
	bookTitle := r.URL.Query().Get("bookTitle")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	rentPrice, _ := strconv.Atoi(r.URL.Query().Get("rentPrice"))
	branchName := r.URL.Query().Get("branchName")

	if bookId != 0 {
		query += "a.bookId = " + strconv.Itoa(bookId)
	} else if bookTitle != "" {
		query += "a.bookTitle = '" + bookTitle + "'"
	} else if author != "" {
		query += "a.author = '" + author + "'"
	} else if genre != "" {
		query += "a.genre = '" + genre + "'"
	} else if year != 0 {
		query += "a.year > " + strconv.Itoa(year)
	} else if rentPrice != 0 {
		query += "a.rentPrice > " + strconv.Itoa(rentPrice)
	} else if branchName != "" {
		query += "c.branchName = '" + branchName + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		log.Println(err)
		return
	}

	// query = "SELECT bookId, coverPath, bookTitle, author, genre, year, page, rentPrice, bookStock FROM books WHERE "
	var book model.Book
	var books []model.Book
	for rows.Next() {
		if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			log.Println(err)
			return
		} else {
			books = append(books, book)
		}
	}

	if len(books) != 0 {
		sendSuccessResponse(w, "Get Success")
	} else {
		sendBadRequestResponse(w, "Error Array Size Not Correct")
	}

	db.Close()
}
