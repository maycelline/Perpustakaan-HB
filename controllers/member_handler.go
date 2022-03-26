package controllers

import (
	"net/http"

	"Perpustakaan-HB/model"
	// "github.com/gorilla/mux"
)

func GetAUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT users.userId, users.fullName, users.userName, users.password, members.balance FROM users JOIN members ON users.userId = members.memberId WHERE members.Id = ?"

	rows := db.QueryRow(query, 1)

	var user model.User
	if err := rows.Scan(&user.ID, &user.FullName, &user.Username, &user.Password, &user.Balance); err != nil {
		// log.Println(err.Error())
		// response := errorTableField()
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode((response))
		return
	}

	// var response UsersResponse
	// 	response.Status = 200
	// 	response.Message = "Success"
	// 	response.Data = users
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode((response))
}

func GetMemberCart(w http.ResponseWriter, r *http.Request) {
	return
}

func CreateBorrowingList(w http.ResponseWriter, r *http.Request) {
	return
}

func GetOngoingBorrowing(w http.ResponseWriter, r *http.Request) {

}

func EditUserProfile(w http.ResponseWriter, r *http.Request) {

}

func EditUserPassword(w http.ResponseWriter, r *http.Request) {

}

func TopupUserBalance(w http.ResponseWriter, r *http.Request) {

}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {

}

func GetMemberHistory(w http.ResponseWriter, r *http.Request) {

}
