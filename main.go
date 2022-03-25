package main

import (
	controllers "Tools/controllers"
	"fmt"
)

func main() {
	var data controllers.DataBorrowed
	data.UserName = "Maycelline"
	data.CourierName = "Dadang Sudrajat"
	data.OrderDate = "20 Mei 2021"
	data.Time = "19.00"

	var book1 controllers.Book
	book1.Title = "Daun yang jatuh tak pernah membenci angin"
	book1.Author = "Tere Liye"

	var book2 controllers.Book
	book2.Title = "Siksa Kubur"
	book2.Author = "Testing"

	var book3 controllers.Book
	book3.Title = "Dear Nathan"
	book3.Author = "Rintiksedu"

	var books []controllers.Book
	books = append(books, book1)
	books = append(books, book2)
	books = append(books, book3)
	data.Books = books

	var branch controllers.Branch
	branch.Name = "Cikutra"
	branch.Address = "Jalan cikutra no 19"

	data.Branch = branch

	fmt.Println(data.Books)

	controllers.SendBorrowAcceptedEmail("maycelinesudarsono@gmail.com", data)

}
