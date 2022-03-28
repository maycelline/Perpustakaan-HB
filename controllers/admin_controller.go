package controllers

import (
	"Perpustakaan-HB/model"
	"log"
	"net/http"
)

func GetAdminData(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	adminId := r.URL.Query().Get("admin_id")
	query := "SELECT u.FullName, b.branchId, b.branchName ,b.branchAddress FROM users u JOIN admins a ON u.userId = a.adminId JOIN branches b ON a.branchId = b.branchId WHERE a.adminId = " + adminId + "; "

	row := db.QueryRow(query)

	var user model.User
	var branch model.Branch

	if err := row.Scan(&user.FullName, &branch.ID, &branch.Name, &branch.Address); err != nil {
		log.Println(err)
		sendBadRequestResponse(w, "Bad Query")
	} else {
		var admin model.Admin
		admin.User = user
		admin.Branch = branch
		sendSuccessResponse(w, "Success", admin)
	}
}

func ApproveBorrowing(w http.ResponseWriter, r *http.Request) {
	// return
}

func ApproveUserReturn(w http.ResponseWriter, r *http.Request) {
	// return
}

func ChangeBorrowingState(w http.ResponseWriter, r *http.Request) {
	// return
}

func CreateNewBook(w http.ResponseWriter, r *http.Request) {
	// return
}
