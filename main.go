package main

import (
	"fmt"
	"log"
	"net/http"

	"Tools/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	// EndPoint
	// router.HandleFunc("/login").Methods("POST")
	// router.HandleFunc("/register").Methods("POST")
	// router.HandleFunc("/book/popular").Methods("GET")
	// router.HandleFunc("/logout").Methods("POST")

	// Member
	router.HandleFunc("/book/list", controllers.GetAllBooks).Methods("GET")
	router.HandleFunc("/member/cart/{member_id}", controllers.GetMemberCart).Methods("GET")

	// router.HandleFunc("/book/list", controllers.Authenticate(controllers.GetAllBooks, 0)).Methods("GET")
	// router.HandleFunc("/member/cart/{member_id}").Methods("GET")
	// router.HandleFunc("/member/borrowing/checkout/{member_id}").Methods("POST")
	// router.HandleFunc("/member/return/{member_id}").Methods("GET")
	// router.HandleFunc("/member/profile/{member_id}", controllers.GetAUsers).Methods("GET")
	// router.HandleFunc("/member/profile/edit/{member_id}").Methods("PUT")
	// router.HandleFunc("/member/password/edit/{member_id}").Methods("PUT")
	// router.HandleFunc("/member/topup/{member_id}").Methods("POST")
	// router.HandleFunc("/member/delete/{member_id}").Methods("DELETE")

	// ADMIN
	// router.HandleFunc("/admin/home").Methods("GET")
	// router.HandleFunc("/admin/borrowApprove").Methods("GET")
	// router.HandleFunc("/admin/returnApprove").Methods("GET")
	// router.HandleFunc("/admin/chooseCourier/{borrow_id}").Methods("PUT")
	// router.HandleFunc("/admin/caddBook").Methods("POST")

	// OWNER
	// router.HandleFunc("/owner/home").Methods("GET")
	// router.HandleFunc("/owner/branchIncome").Methods("GET")
	// router.HandleFunc("owner/income").Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(router)

	http.Handle("/", router)
	fmt.Println("Connected to port 8080")
	log.Println(http.ListenAndServe(":8080", handler))

	// var data model.BorrowDataHTML
	// var borrows []model.Borrowing
	// var borrowing1 model.Borrowing
	// var borrowing2 model.Borrowing

	// borrowing1.Book.Title = "Daun yang jatuh tidak pernah membenci angin"
	// borrowing1.Book.Author = "Tere Liye"
	// // borrowing1.Book.BranchName = "Cikutra"
	// borrowing1.BorrowDate = time.Now()

	// borrowing2.Book.Title = "Please, Look after Mom"
	// borrowing2.Book.Author = "Tere Liye"
	// // borrowing2.Book.BranchName = "Cikutra"
	// borrowing2.BorrowDate = time.Now()

	// borrows = append(borrows, borrowing1)
	// borrows = append(borrows, borrowing2)

	// var courier model.Courier
	// courier.CourierName = "Dadang Sudrajat"

	// // var Branch model.Branch

	// data.Borrows = borrows
	// data.Courier = courier
	// data.User.FullName = "Maycelline Selvyanti"
	// data.Branch.Name = "Cikutra"
	// data.Branch.Address = "Jl Cikutra no 19"
	// data.CourierCome = time.Now().Add(time.Minute * 30)

	// controllers.SendBorrowAcceptedEmail("maycelinesudarsono@gmail.com", data)

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

	// 	key := "key-1"
	//     data := "Hello Redis"
	//     ttl := time.Duration(3) * time.Second

	//     // store data using SET command
	//     op1 := rdb.Set(context.Background(), key, data, ttl)
	//     if err := op1.Err(); err != nil {
	//         fmt.Printf("unable to SET data. error: %v", err)
	//         return
	//     }
	//     log.Println("set operation success")

	//     // ...
	// }
}
