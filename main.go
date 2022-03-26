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
	controllers.WeeklyEmailScheduler()

	router := mux.NewRouter()

	// EndPoint
	router.HandleFunc("/login", controllers.CheckUserLogin).Methods("POST")
	router.HandleFunc("/register", controllers.CreateUserRegister).Methods("POST")
	router.HandleFunc("/book/popular", controllers.GetPopularBook).Methods("GET")
	router.HandleFunc("/logout", controllers.UserLogout).Methods("POST")

	// Member
	// router.HandleFunc("/book/list",).Methods("GET")
	router.HandleFunc("/member/cart/{member_id}", controllers.GetMemberCart).Methods("GET")
	router.HandleFunc("/member/borrowing/checkout/{member_id}", controllers.CreateBorrowingList).Methods("POST")
	router.HandleFunc("/member/return/{member_id}", controllers.GetOngoingBorrowing).Methods("GET")
	router.HandleFunc("/member/profile/{member_id}", controllers.GetAUser).Methods("GET")
	router.HandleFunc("/member/profile/edit/{member_id}", controllers.EditUserProfile).Methods("PUT")
	router.HandleFunc("/member/password/edit/{member_id}", controllers.EditUserPassword).Methods("PUT")
	router.HandleFunc("/member/topup/{member_id}", controllers.TopupUserBalance).Methods("POST")
	router.HandleFunc("/member/delete/{member_id}", controllers.DeleteAccount).Methods("DELETE")
	router.HandleFunc("/member/borrowHistory/{member_id}", controllers.GetMemberHistory).Methods("GET")
	// ADMIN
	router.HandleFunc("/admin/home", controllers.GetAdminData).Methods("GET")
	router.HandleFunc("/admin/borrowApprove", controllers.ApproveBorrowing).Methods("GET")
	router.HandleFunc("/admin/returnApprove", controllers.ApproveUserReturn).Methods("GET")
	router.HandleFunc("/admin/chooseCourier/{borrow_id}", controllers.ChangeBorrowingState).Methods("PUT")
	router.HandleFunc("/admin/caddBook", controllers.CreateNewBook).Methods("POST")

	// OWNER
	// router.HandleFunc("/owner/home").Methods("GET")
	// router.HandleFunc("/owner/branchIncome").Methods("GET")
	// router.HandleFunc("owner/income").Methods("GET")

	// CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true, //kalo ga nanti ga bisa ngakses  karena cookies dkk
	})

	handler := corsHandler.Handler(router)

	http.Handle("/", handler)
	fmt.Println("Connected to port 8080")
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
