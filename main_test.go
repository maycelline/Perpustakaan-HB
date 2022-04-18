package main

import (
	"Perpustakaan-HB/controllers"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := controllers.Connect()
	ensureTableExists(db)

	data := url.Values{}
	data.Set("fullName", "Debora Lexandra")
	data.Set("userName", "deboras123")
	data.Set("birthDate", "2000-01-01")
	data.Set("phoneNumber", "0853-1234-1234")
	data.Set("email", "deboras@gmail.com")
	data.Set("address", "Jalan Dago No. 12")
	data.Set("additionalAddress", "Kota Bandung, Jawa Barat 40132")
	data.Set("password", "Deboras1")
	data.Set("confirmPassword", "Deboras1")

	req, err := http.NewRequest("POST", "/register", strings.NewReader(data.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.CreateUserRegister)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(tableUsersCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(tableMembersCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableUsersCreationQuery = `
CREATE TABLE IF NOT EXISTS users
(
    userId INT(11) AUTO_INCREMENT PRIMARY KEY,
    fullName varchar(255) DEFAULT NULL,
	userName varchar(255) DEFAULT NULL,
	birthDate date DEFAULT NULL,
	phoneNumber varchar(255) DEFAULT NULL,
	email varchar(255) DEFAULT NULL,
	address varchar(255) DEFAULT NULL,
	additionalAddress varchar(255) DEFAULT NULL,
	password varchar(255) DEFAULT NULL,
	userType varchar(255) DEFAULT NULL
)`

const tableMembersCreationQuery = `
CREATE TABLE IF NOT EXISTS members
(
    memberId int(11) DEFAULT NULL,
  	balance int(11) DEFAULT NULL,
	FOREIGN KEY (memberId) references users(userId)
)`
