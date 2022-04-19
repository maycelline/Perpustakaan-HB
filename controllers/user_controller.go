package controllers

import (
	"Perpustakaan-HB/model"
	"crypto/md5"
	"encoding/hex"
	_ "encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/dlclark/regexp2"
)

func CheckUserLogin(w http.ResponseWriter, r *http.Request) {
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		sendServerErrorResponse(w, "Internal Server Error")
		return
	}

	password := encodePassword(r.Form.Get("password"))
	userName := r.Form.Get("userName")
	// fmt.Println(password)
	// fmt.Println(userName)

	if password != "" && userName != "" {
		query := "SELECT * FROM users WHERE password = ? AND userName = ?"

		var user model.User

		rows := db.QueryRow(query, password, userName)
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.AdditionalAddress, &user.Password, &user.UserType); err != nil {
			fmt.Println(err)
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
					log.Println(err)
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
	db := Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
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

	birthDateTime, _ := time.Parse("YYYY-MM-DD", birthDate)
	var user model.User = model.User{FullName: fullName, UserName: userName, BirthDate: birthDateTime, PhoneNumber: phone, Email: email, Address: address, AdditionalAddress: additionalAddress, Password: password, UserType: "MEMBER"}
	var checkPass = checkPasswordValidation(password, w)
	var checkUname = checkUsernameValidation(userName, w)
	var checkMail = checkMailValidation(email, w)
	var memberId int
	if password == confirmPass && checkPass && checkUname && checkMail {
		if fullName != "" && phone != "" && address != "" {
			result1, errQuery1 := db.Exec("INSERT INTO users(fullName, userName, birthDate, phoneNumber, email, address, additionalAddress, password, userType) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
				fullName,
				userName,
				birthDate,
				phone,
				email,
				address,
				additionalAddress,
				encodePassword(password),
				"MEMBER",
			)
			// memberId, _ = result1.LastInsertId()
			// memberId,_ = strconv.Atoi(memberid)
			if errQuery1 != nil {
				log.Println(errQuery1)
				sendBadRequestResponse(w, "Error Can Not Register, error query 1")
				return
			}

			tempId, _ := result1.LastInsertId()
			memberId = int(tempId)
			_, errQuery2 := db.Exec("INSERT INTO members(memberId, balance) values (?,?)", tempId, 0)

			if errQuery2 != nil {
				log.Println(errQuery2)
				sendBadRequestResponse(w, "Error Can Not Register, error query 2")
				return
			}
		} else {
			sendBadRequestResponse(w, "Error Missing Values")
			return
		}
	} else {
		sendBadRequestResponse(w, "Your input not valid")
		return
	}
	var member model.Member

	member.User = user
	member.User.ID = memberId
	generateMemberToken(w, member)
	sendSuccessResponse(w, "Register Success", nil)
	go SendRegisterEmail(user)
	SetEmailWeeklyScheduler(email)
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
	db := Connect()
	defer db.Close()

	query := "SELECT userId, fullName, userName, birthDate, phoneNumber, email, address, additionalAddress, password, userType FROM users WHERE userType = 'MEMBER'"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return nil
	}

	var users []model.User
	var user model.User

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.FullName, &user.UserName, &user.BirthDate, &user.PhoneNumber, &user.Email, &user.Address, &user.AdditionalAddress, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			return nil
		} else {
			users = append(users, user)
		}
	}
	return users
}

func checkPasswordValidation(pass string, w http.ResponseWriter) bool {
	regex, err := regexp2.Compile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)([^\@$!%*?&/^\s]){8,10}$`, 0)
	if err != nil {
		sendBadRequestResponse(w, "Regex Not Correct")
		return false
	}
	checkPass, err := regex.MatchString(pass)
	if err != nil {
		sendBadRequestResponse(w, "Password Not Match Criteria")
	}
	return checkPass
}

func checkUsernameValidation(username string, w http.ResponseWriter) bool {
	regex, err := regexp2.Compile(`^(?=.*[a-zA-Z])(?=.*\d)([^\@$!%*?&/^\s]){4,16}$`, 0)
	if err != nil {
		sendBadRequestResponse(w, "Regex Not Correct") //ini perlu ?
		return false
	}
	checkUname, err := regex.MatchString(username)
	if err != nil {
		sendBadRequestResponse(w, "Password Not Match Criteria") //ini perlu ?
	}
	return checkUname
}

func checkMailValidation(email string, w http.ResponseWriter) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		sendBadRequestResponse(w, "Mail Not Correct") //ini perlu ?
		return false
	}
	return true
}
