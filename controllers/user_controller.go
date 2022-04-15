package controllers

import (
	"crypto/md5"
	"encoding/hex"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Perpustakaan-HB/model"
)

func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	password := encodePassword(r.Form.Get("password"))
	userName := r.Form.Get("userName")
	fmt.Println(password)
	fmt.Println(userName)

	if password != "" && userName != "" {
		query := "SELECT * FROM users WHERE password = ? AND username = ?"

		var user model.User

		rows := db.QueryRow(query, password, userName)
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.AdditionalAddress, &user.Password, &user.UserType); err != nil {
			sendBadRequestResponse(w, "Error Field Undefined")
			return
		}

		if user.FullName != "" {
			if user.UserType == "MEMBER" {
				var member model.Member
				member.User = user
				query = "SELECT balance FROM members WHERE memberId = ?"
				rows = db.QueryRow(query, member.User.ID)
				if err := rows.Scan(&member.Balance); err != nil {
					return
				}
				generateMemberToken(w, member)
				sendSuccessResponse(w, "Login Success", member)
			} else if user.UserType == "ADMIN" {
				var admin model.Admin
				admin.User = user
				query = "SELECT branches.branchId, branches.branchName, branches.branchAddress FROM admins JOIN branches WHERE admins.adminId = ? AND admins.branchId = branches.branchId"
				rows = db.QueryRow(query, admin.User.ID)
				if err := rows.Scan(&admin.Branch.ID, &admin.Branch.Name, &admin.Branch.Address); err != nil {
					return
				}
				generateAdminToken(w, admin)
				sendSuccessResponse(w, "Login Success", admin)
			} else {
				generateOwnerToken(w, user)
				sendSuccessResponse(w, "Login Success", user)
			}
		} else {
			sendNotFoundResponse(w, "User Not Found")
		}
	} else {
		sendBadRequestResponse(w, "Error Field Undefined")
	}
}

func CreateUserRegister(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	fullName := r.Form.Get("fullName")
	userName := r.Form.Get("userName")
	birthDate := r.Form.Get("birthDate")
	phone := r.Form.Get("phoneNumber")
	email := r.Form.Get("email")
	address := r.Form.Get("address")
	additionalAddress := r.Form.Get("additionalAddress")
	password := r.Form.Get("password")
	confirmPass := r.Form.Get("confirmPassword")

	passwordLength := len(password)

	if passwordLength < 8 {
		sendBadRequestResponse(w, "Need more character")
		return
	} else if passwordLength > 10 {
		sendBadRequestResponse(w, "Too many character")
		return
	}

	containsNumber := 0
	for i := 0; i < 10; i++ {
		number := strconv.Itoa(i)
		if strings.Contains(password, number) {
			containsNumber = containsNumber + 1
		}
	}

	passwordCheck := strings.ToLower(password)
	arrayPassword := []rune(passwordCheck)

	containsLowerCase := 0
	for i := 0; i < passwordLength; i++ {
		char := string(arrayPassword)
		if strings.Contains(password, char) {
			containsLowerCase = containsLowerCase + 1
		}
	}

	if containsNumber == 0 || containsLowerCase == 0 || containsLowerCase == containsNumber {
		sendBadRequestResponse(w, "Bad password")
		return
	}

	if password == confirmPass {
		if fullName != "" && userName != "" && phone != "" && address != "" && password != "" {
			result1, errQuery1 := db.Exec("INSERT INTO users(fullName, userName, birthDate, phoneNumber, email, address, additionalAddress, password, userType) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
				fullName,
				userName,
				birthDate, // belum beres kayanya
				phone,
				email,
				address,
				additionalAddress,
				encodePassword(password),
				"MEMBER",
			)
			tempId, _ := result1.LastInsertId()
			_, errQuery2 := db.Exec("INSERT INTO members(id, balance) values (?,?)", tempId, 0)

			if errQuery1 != nil && errQuery2 != nil {
				sendBadRequestResponse(w, "Error Can Not Register")
			} else {
				sendSuccessResponse(w, "Register Success", nil)
			}
		} else {
			sendBadRequestResponse(w, "Error Missing Values")
		}
	} else {
		sendBadRequestResponse(w, "Error Password Does Not Match")
	}

	go SetScheduler(email)
	SetUsersCache(nil)
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	sendSuccessResponse(w, "Logout Success", nil)
}

func encodePassword(pass string) string {
	encodePass := md5.Sum([]byte(pass))
	return hex.EncodeToString(encodePass[:])
}

func GetAllUsers() []model.User {
	var users []model.User
	users = GetUsersFromCache()

	if users == nil {
		db := connect()
		defer db.Close()

		query := "SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, additionalAddress, password, userType from users"

		rows, err := db.Query(query)
		if err != nil {
			log.Println(err)
			return nil
		}

		var user model.User
		for rows.Next() {
			if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.AdditionalAddress, &user.Password, &user.UserType); err != nil {
				log.Println(err)
				return nil
			} else {
				users = append(users, user)
			}
		}
		SetUsersCache(users)
	}

	return users
}
