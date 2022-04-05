package controllers

import (
	"crypto/md5"
	"encoding/hex"
	_ "encoding/json"
	"fmt"
	_ "log"
	"net/http"

	"Perpustakaan-HB/model"
	// "github.com/gorilla/mux"
)

func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	password := encodePassword(r.Form.Get("password"))
	fmt.Println(password)
	userName := r.Form.Get("userName")

	if password != "" && userName != "" {
		query := "SELECT * FROM users WHERE password = ? AND username = ?"

		var user model.User

		rows := db.QueryRow(query, password, userName)
		if err := rows.Scan(&user.ID, &user.FullName, &user.Username, &user.BirthDate, &user.Phone, &user.Email, &user.Address, &user.AdditionalAddress, &user.Password, &user.UserType); err != nil {
			// log.Println(err.Error())
			// response := errorTableField()
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode((response))
			return
		}

		if user.FullName != "" {
			// response := successUserInfoProcess()
			// response.Data = GetAUserInfo(user)
			// generateToken(w, user.ID, user.Name, user.UserType)
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode((response))
			if user.UserType == "MEMBER" {
				var member model.Member
				member.User = user
				query = "SELECT balance FROM members WHERE memberId = ?"
				rows = db.QueryRow(query, member.User.ID)
				if err := rows.Scan(&member.Balance); err != nil {
					return
				}
			} else if user.UserType == "ADMIN" {
				var admin model.Admin
				admin.User = user
				query = "SELECT branches.branchId, branches.branchName, branches.branchAddress FROM admins JOIN branches WHERE admins.adminId = ? AND admins.branchId = branches.branchId"
				rows = db.QueryRow(query, admin.User.ID)
				if err := rows.Scan(&admin.Branch.ID, &admin.Branch.Name, &admin.Branch.Address); err != nil {
					return
				}
			} else {

			}
		} else {
			// response := errorUserNotFound()
			// w.Header().Set("Content-Type", "application/json")
			// json.NewEncoder(w).Encode((response))
		}
	} else {
		// response := errorEmptyForm()
		// w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode((response))
	}
}

func CreateUserRegister(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {

	}

	fullName := r.Form.Get("fullName")
	userName := r.Form.Get("userName")
	phone := r.Form.Get("PhoneNumber")
	address := r.Form.Get("Address")
	additionalAddress := r.Form.Get("Additional Address")
	password := r.Form.Get("password")
	confirmPass := r.Form.Get("password")

	if password == confirmPass {
		if fullName != "" && userName != "" && phone != "" && address != "" && password != "" {
			result1, errQuery1 := db.Exec("INSERT INTO users(fullName, userName, birthDate, phoneNumber, email, address, additionalAddress, password, userType) values (?,?,?,?,?,?,?,?,?)", fullName, userName, "", phone, address, additionalAddress, encodePassword(password))
			tempId, _ := result1.LastInsertId()
			_, errQuery2 := db.Exec("INSERT INTO members(id, balance) values (?,?)", tempId, 0)

			if errQuery1 != nil && errQuery2 != nil {

			} else {

			}
		} else {

		}
	} else {

	}

}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	//ngurusin cookie nanti pokonya

	//nampilin success info

	//udah

}

func encodePassword(pass string) string {
	encodePass := md5.Sum([]byte(pass))
	return hex.EncodeToString(encodePass[:])
}
