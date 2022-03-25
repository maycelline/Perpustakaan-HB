package main

import (
	controllers "Tools/controllers"
	// "fmt"
	"Tools/model"
	"time"
)

func main() {
	var data model.BorrowDataHTML
	var borrows []model.Borrowing
	var borrowing1 model.Borrowing
	var borrowing2 model.Borrowing

	borrowing1.Book.Title = "Daun yang jatuh tidak pernah membenci angin"
	borrowing1.Book.Author = "Tere Liye"
	// borrowing1.Book.BranchName = "Cikutra"
	borrowing1.BorrowDate = time.Now()

	borrowing2.Book.Title = "Please, Look after Mom"
	borrowing2.Book.Author = "Tere Liye"
	// borrowing2.Book.BranchName = "Cikutra"
	borrowing2.BorrowDate = time.Now()

	borrows = append(borrows, borrowing1)
	borrows = append(borrows, borrowing2)

	var courier model.Courier
	courier.CourierName = "Dadang Sudrajat"

	// var Branch model.Branch

	data.Borrows = borrows
	data.Courier = courier
	data.User.FullName = "Maycelline Selvyanti"
	data.Branch.Name = "Cikutra"
	data.Branch.Address = "Jl Cikutra no 19"
	data.CourierCome = time.Now().Add(time.Minute * 30)

	controllers.SendBorrowAcceptedEmail("maycelinesudarsono@gmail.com", data)

	// data.Time = "19.00"

	// var book1 controllers.Book
	// book1.Title = "Daun yang jatuh tak pernah membenci angin"
	// book1.Author = "Tere Liye"

	// var book2 controllers.Book
	// book2.Title = "Siksa Kubur"
	// book2.Author = "Testing"

	// var book3 controllers.Book
	// book3.Title = "Dear Nathan"
	// book3.Author = "Rintiksedu"

	// var books []controllers.Book
	// books = append(books, book1)
	// books = append(books, book2)
	// books = append(books, book3)
	// data.Books = books

	// var branch controllers.Branch
	// branch.Name = "Cikutra"
	// branch.Address = "Jalan cikutra no 19"

	// data.Branch = branch

	// fmt.Println(data.Books)

	// controllers.SendBorrowAcceptedEmail("maycelinesudarsono@gmail.com", data)

}
