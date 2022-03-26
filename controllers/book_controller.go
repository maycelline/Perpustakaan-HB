package controllers

import (
	// "crypto/md5"
	// "encoding/hex"
	"Perpustakaan-HB/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	// "github.com/gorilla/mux"
)

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	// bookId, _ := strconv.Atoi(r.URL.Query().Get("bookId"))
	bookTitle := r.URL.Query().Get("bookTitle")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	rentPrice, _ := strconv.Atoi(r.URL.Query().Get("rentPrice"))
	branchName := r.URL.Query().Get("branchName")

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c WHERE b.branchId = c.branchId AND c.branchName = '" + branchName + "'"

	if bookTitle != "" {
		query += " AND a.bookTitle = '" + bookTitle + "'"
	} else if author != "" {
		query += " AND a.author = '" + author + "'"
	} else if genre != "" {
		query += " AND a.genre = '" + genre + "'"
	} else if year != 0 {
		query += " AND a.year > " + strconv.Itoa(year)
	} else if rentPrice != 0 {
		query += " AND a.rentPrice > " + strconv.Itoa(rentPrice)
	} else if branchName != "" {
		query += " AND c.branchName = '" + branchName + "'"
	}

	log.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		log.Println(err)
		return
	}

	var book model.Book
	var books []model.Book
	for rows.Next() {
		if err := rows.Scan(&book.ID, &book.CoverPath, &book.Title, &book.Author, &book.Genre, &book.Year, &book.Page, &book.RentPrice, &book.Stock); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			log.Println(err)
			return
		} else {
			log.Println(books)
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

func GetMemberCart(w http.ResponseWriter, r *http.Request) {

}

func GetPopularBook(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT a.bookTitle, a.coverPath, a.author, a.genre, a.year, a.page FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN borrows c ON b.stockId = c.stockId GROUP BY b.bookId"

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		log.Println(err)
		return
	}

	var book model.PopularBook
	var books []model.PopularBook
	for rows.Next() {
		if err := rows.Scan(&book.Title, &book.CoverPath, &book.Author, &book.Genre, &book.Year, &book.Page); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			log.Println(err)
			return
		} else {
			books = append(books, book)
		}
	}

	json, err := json.Marshal(books)
	if err != nil {
		fmt.Println(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err = client.Set("popular", json, 0).Err()
	if err != nil {
		return
	}

	val, err := client.Get("popular").Result()
	if err != nil {
		return
	}
	sendSuccessResponse(w, "Get Success", val)
}
