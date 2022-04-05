package controllers

import (
	"Perpustakaan-HB/model"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	branchName := r.URL.Query().Get("branchName")
	bookTitle := r.URL.Query().Get("bookTitle")
	author := r.URL.Query().Get("author")
	genre := r.URL.Query().Get("genre")
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	rentPrice, _ := strconv.Atoi(r.URL.Query().Get("rentPrice"))

	query := "SELECT a.bookId, a.coverPath, a.bookTitle, a.author, a.genre, a.year, a.page, a.rentPrice, b.stock, c.branchName FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN branches c WHERE b.branchId = c.branchId AND c.branchName = '" + branchName + "'"

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

func GetPopularBook(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT a.bookId, a.bookTitle, a.coverPath, a.author, a.genre, a.year FROM books a JOIN stocks b ON a.bookId = b.bookId JOIN borrows c ON b.stockId = c.stockId GROUP BY b.bookId"

	rows, err := db.Query(query)
	if err != nil {
		sendNotFoundResponse(w, "Table Not Found")
		return
	}

	var book model.Book
	var books []model.Book
	for rows.Next() {
		if err := rows.Scan(&book.ID, &book.Title, &book.CoverPath, &book.Author, &book.Genre, &book.Year); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		} else {
			books = append(books, book)
		}
	}

	converted, err := json.Marshal(books)
	if err != nil {
		sendBadRequestResponse(w, "Error JSON Undefined")
		return
	}

	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err = client.Set(ctx, "popular", converted, 0).Err()
	if err != nil {
		sendBadRequestResponse(w, "Error Redis Undefined")
		return
	}

	value, err := client.Get(ctx, "popular").Result()
	if err != nil {
		sendBadRequestResponse(w, "Error Redis Undefined")
		return
	}

	_ = json.Unmarshal([]byte(value), &books)

	sendSuccessResponse(w, "Get Success", books)
}
