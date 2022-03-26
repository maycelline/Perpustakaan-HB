package main

import (
	_ "log"
	_ "net/http"

	controllers "Perpustakaan-HB/controllers"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
	_ "github.com/rs/cors"
)

func main() {
	controllers.WeeklyEmailScheduler()
	// use
	// router := mux.NewRouter()

	// EndPoint
	// router.HandleFunc("/login",).Methods("POST")
	// router.HandleFunc("/register",).Methods("POST")
	// router.HandleFunc("/book/popular",).Methods("GET")
	// router.HandleFunc("/logout",).Methods("POST")

	// Member
	// router.HandleFunc("/book/list",).Methods("GET")
	// router.HandleFunc("/member/cart/{member_id}", ).Methods("GET")
	// router.HandleFunc("/member/borrowing/checkout/{member_id}", ).Methods("POST")
	// router.HandleFunc("/member/return/{member_id}",).Methods("GET")
	// router.HandleFunc("/member/profile/{member_id}", controllers.GetAUsers).Methods("GET")
	// router.HandleFunc("/member/profile/edit/{member_id}",).Methods("PUT")
	// router.HandleFunc("/member/password/edit/{member_id}",).Methods("PUT")
	// router.HandleFunc("/member/topup/{member_id}",).Methods("POST")
	// router.HandleFunc("/member/delete/{member_id}",).Methods("DELETE")

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

	// CORS
	// 	corsHandler := cors.New(cors.Options{
	// 		AllowedOrigins:   []string{"*"},
	// 		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	// 		AllowCredentials: true, //kalo ga nanti ga bisa ngakses  karena cookies dkk
	// 	})

	// 	handler := corsHandler.Handler(router)

	// 	http.Handle("/", handler)
	// 	fmt.Println("Connected to port 8080")
	// 	log.Println("Connected to port 8080")
	// 	log.Fatal(http.ListenAndServe(":8080", handler))
}
