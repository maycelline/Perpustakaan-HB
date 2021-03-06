package main

import (
	"fmt"
	"log"
	"net/http"

	controllers "Perpustakaan-HB/controllers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// scheduler initiation for all users before API start
	controllers.WeeklyEmailScheduler()
	controllers.SetEmailBorrowingInfoScheduler()

	router := mux.NewRouter()

	// General
	router.HandleFunc("/login", controllers.CheckUserLogin).Methods("POST")
	router.HandleFunc("/register", controllers.CreateUserRegister).Methods("POST")
	router.HandleFunc("/book/popular", controllers.GetPopularBooks).Methods("GET")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// Member (1)
	router.HandleFunc("/book/list", controllers.Authenticate(controllers.GetAllBooks, 1)).Methods("GET")
	router.HandleFunc("/member/cart", controllers.Authenticate(controllers.GetMemberCart, 1)).Methods("GET")
	router.HandleFunc("/member/cart/add", controllers.Authenticate(controllers.AddBookToCart, 1)).Methods("POST")
	router.HandleFunc("/member/cart/remove/{branch_name}", controllers.Authenticate(controllers.RemoveBookFromCart, 1)).Methods("DELETE")
	router.HandleFunc("/member/borrowing/checkout", controllers.Authenticate(controllers.CheckoutBorrowing, 1)).Methods("POST")
	router.HandleFunc("/member/borrowing/return", controllers.Authenticate(controllers.ReturnBorrowing, 1)).Methods("GET")
	router.HandleFunc("/member/profile", controllers.Authenticate(controllers.GetUserData, 1)).Methods("GET")
	router.HandleFunc("/member/profile/edit", controllers.Authenticate(controllers.EditUserProfile, 1)).Methods("PUT")
	router.HandleFunc("/member/password/edit", controllers.Authenticate(controllers.EditUserPassword, 1)).Methods("PUT")
	router.HandleFunc("/member/topup", controllers.Authenticate(controllers.TopupUserBalance, 1)).Methods("PUT")
	router.HandleFunc("/member/borrowHistory", controllers.Authenticate(controllers.GetMemberHistory, 1)).Methods("GET")

	// Admin (2)
	router.HandleFunc("/admin/home", controllers.Authenticate(controllers.GetAdminData, 2)).Methods("GET")
	router.HandleFunc("/admin/borrowApprove", controllers.Authenticate(controllers.GetUnapprovedBorrowing, 2)).Methods("GET")
	router.HandleFunc("/admin/returnApprove", controllers.Authenticate(controllers.GetUnapprovedReturn, 2)).Methods("GET")
	router.HandleFunc("/admin/chooseCourier/{borrow_id}", controllers.Authenticate(controllers.ChangeBorrowingState, 2)).Methods("PUT")
	router.HandleFunc("/admin/addBook", controllers.Authenticate(controllers.AddNewBook, 2)).Methods("POST")

	// Owner (3)
	router.HandleFunc("/owner/home", controllers.Authenticate(controllers.GetOwnerData, 3)).Methods("GET")
	router.HandleFunc("/owner/branchIncome", controllers.Authenticate(controllers.GetBranchIncome, 3)).Methods("GET")
	router.HandleFunc("/owner/income", controllers.Authenticate(controllers.GetAllIncome, 3)).Methods("GET")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://phb.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
