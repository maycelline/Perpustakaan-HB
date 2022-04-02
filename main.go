package main

import (
	"fmt"
	"log"
	"net/http"

	"Perpustakaan-HB/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// controllers.WeeklyEmailScheduler()

	router := mux.NewRouter()

	// General
	router.HandleFunc("/login", controllers.CheckUserLogin).Methods("POST")
	router.HandleFunc("/register", controllers.CreateUserRegister).Methods("POST")
	router.HandleFunc("/book/popular", controllers.GetPopularBook).Methods("GET")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// Member
	router.HandleFunc("/book/list", controllers.GetAllBooks).Methods("GET")
	router.HandleFunc("/member/cart", controllers.GetMemberCart).Methods("GET")
	router.HandleFunc("/member/borrowing/checkout", controllers.CreateBorrowingList).Methods("POST")
	router.HandleFunc("/member/return", controllers.GetOngoingBorrowing).Methods("GET")
	router.HandleFunc("/member/profile", controllers.GetAUser).Methods("GET")
	router.HandleFunc("/member/profile/edit", controllers.EditUserProfile).Methods("PUT")
	router.HandleFunc("/member/password/edit", controllers.EditUserPassword).Methods("PUT")
	router.HandleFunc("/member/topup", controllers.TopupUserBalance).Methods("POST")
	router.HandleFunc("/member/delete", controllers.DeleteAccount).Methods("DELETE")
	router.HandleFunc("/member/borrowHistory", controllers.GetMemberHistory).Methods("GET")
	// ADMIN
	router.HandleFunc("/admin/home", controllers.GetAdminData).Methods("GET")
	router.HandleFunc("/admin/borrowApprove", controllers.ApproveBorrowing).Methods("GET")
	router.HandleFunc("/admin/returnApprove", controllers.ApproveUserReturn).Methods("GET")
	router.HandleFunc("/admin/chooseCourier/{borrow_id}", controllers.ChangeBorrowingState).Methods("PUT")
	router.HandleFunc("/admin/addBook", controllers.CreateNewBook).Methods("POST")

	// OWNER
	router.HandleFunc("/owner/home", controllers.GetOwnerData).Methods("GET")
	router.HandleFunc("/owner/branchIncome", controllers.GetBranchIncome).Methods("GET")
	router.HandleFunc("/owner/income", controllers.GetAllIncome).Methods("GET")

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
