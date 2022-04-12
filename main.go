package main

import (
	"fmt"
	"log"
	"net/http"

	"Perpustakaan-HB/controllers"
	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// controllers.GetBranchIncome()
	// controllers.WeeklyEmailScheduler()

	//inisiasi scheduler untuk user yang telah terdaftar sebelum API dinyalakan
	controllers.WeeklyEmailScheduler()

	router := mux.NewRouter()

	// General
	router.HandleFunc("/login", controllers.CheckUserLogin).Methods("POST")
	router.HandleFunc("/register", controllers.CreateUserRegister).Methods("POST")
	router.HandleFunc("/book/popular", controllers.GetPopularBooks).Methods("GET")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// Member (1)
	router.HandleFunc("/book/list", controllers.Authenticate(controllers.GetAllBooks, 1)).Methods("GET")
	router.HandleFunc("/member/cart", controllers.Authenticate(controllers.GetMemberCart, 1)).Methods("GET")
	router.HandleFunc("/member/cart", controllers.Authenticate(controllers.AddBookToCart, 1)).Methods("POST")
	router.HandleFunc("/member/borrowing/checkout", controllers.Authenticate(controllers.CreateBorrowingList, 1)).Methods("POST")
	router.HandleFunc("/member/return", controllers.Authenticate(controllers.GetOngoingBorrowing, 1)).Methods("GET")
	router.HandleFunc("/member/profile", controllers.Authenticate(controllers.GetAUser, 1)).Methods("GET")
	router.HandleFunc("/member/profile/edit", controllers.Authenticate(controllers.EditUserProfile, 1)).Methods("PUT")
	router.HandleFunc("/member/password/edit", controllers.Authenticate(controllers.EditUserPassword, 1)).Methods("PUT")
	router.HandleFunc("/member/topup", controllers.Authenticate(controllers.TopupUserBalance, 1)).Methods("POST")
	router.HandleFunc("/member/delete", controllers.Authenticate(controllers.DeleteAccount, 1)).Methods("DELETE")
	router.HandleFunc("/member/borrowHistory", controllers.Authenticate(controllers.GetMemberHistory, 1)).Methods("GET")

	// Admin (2)
	router.HandleFunc("/admin/home", controllers.Authenticate(controllers.GetAdminData, 2)).Methods("GET")
	router.HandleFunc("/admin/borrowApprove", controllers.Authenticate(controllers.GetUnapprovedBorrowing, 2)).Methods("GET")
	router.HandleFunc("/admin/returnApprove", controllers.Authenticate(controllers.GetUnapprovedReturn, 2)).Methods("GET")
	router.HandleFunc("/admin/chooseCourier/{borrow_id}", controllers.Authenticate(controllers.ChangeBorrowingState, 2)).Methods("PUT")
	router.HandleFunc("/admin/addBook", controllers.Authenticate(controllers.CreateNewBook, 2)).Methods("POST")

	// Owner (3)
	// router.HandleFunc("/owner/home", controllers.Authenticate(controllers.GetOwnerData, 3)).Methods("GET")
	// router.HandleFunc("/owner/branchIncome", controllers.Authenticate(controllers.GetBranchIncome, 3)).Methods("GET")
	// router.HandleFunc("/owner/income", controllers.Authenticate(controllers.GetAllIncome, 3)).Methods("GET")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
