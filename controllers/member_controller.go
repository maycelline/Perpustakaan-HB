package controllers

import (
	"Perpustakaan-HB/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT users.userId, users.fullName, users.userName, users.password, users.userType, members.balance FROM users JOIN members ON users.userId = members.memberId WHERE members.Id = ?"

	vars := mux.Vars(r)
	memberId, _ := strconv.Atoi(vars["member_id"])

	rows := db.QueryRow(query, memberId)

	var member model.Member
	if err := rows.Scan(&member.User.ID, &member.User.FullName, &member.User.Username, &member.User.Password, &member.User.UserType, &member.Balance); err != nil {
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
