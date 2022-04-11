package controllers

import (
	"bytes"
	"fmt"

	"Perpustakaan-HB/model"
	"text/template"

	"gopkg.in/gomail.v2"
)

func SendRegisterEmail(destinationAddress string, user model.User) {
	mail := gomail.NewMessage()

	template := "assets/email_template/register.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", destinationAddress)
	mail.SetHeader("Subject", "Register Success!")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	}
}

func SendBorrowAcceptedEmail(destinationAddress string, data model.BorrowDataHTML) {
	mail := gomail.NewMessage()

	template := "assets/email_template/accept_borrow_book.html"

	result, _ := parseTemplate(template, data)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", destinationAddress)
	mail.SetHeader("Subject", "Accepted Order")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	}
}

func SendBorrowRejectedEmail(destinationAddress string, data model.BorrowDataHTML) {
	mail := gomail.NewMessage()

	template := "assets/email_template/reject_borrow_book.html"

	result, _ := parseTemplate(template, data)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", destinationAddress)
	mail.SetHeader("Subject", "Rejected Order")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	}
}

//ubah text plain ke dalam bentuk buffer
func parseTemplate(templateFileName string, data interface{}) (string, error) {
	// mengubah text html ke dalam bentuk byte
	t, err := template.ParseFiles(templateFileName)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	buff := new(bytes.Buffer)

	//render struct ke dalam file html (td ada var name di dalam htmlnya)
	if err = t.Execute(buff, data); err != nil {
		fmt.Println(err)
		return "", err
	}

	return buff.String(), nil
}

func SendWeeklyEmail(destinationEmailAddress string) {
	mail := gomail.NewMessage()
	template := "assets/email_template/weekly_news.html"

	var popularBooks model.PopularBooksEmail
	popularBooks.Books = PopularBooks()

	result, _ := parseTemplate(template, popularBooks)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", destinationEmailAddress)
	mail.SetHeader("Subject", "Weekly Popular Books")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to " + destinationEmailAddress)
	}
}

//function ini hanya dijalankan saat API baru dinyalakan untuk membuat scheduler bagi semua user
func WeeklyEmailScheduler() {
	//disini perlu ada panggil function untuk dapetin email dari semua address
	var user1 model.User
	user1.FullName = "Maycelline Selvyanti"
	user1.Email = "maycelinesudarsono@gmail.com"

	var user2 model.User
	user2.FullName = "Silvi Prisillia"
	user2.Email = "if-20019@students.ithb.ac.id"

	var user3 model.User
	user3.FullName = "Feliciana Gunadi"
	user3.Email = "if-20009@students.ithb.ac.id"

	var users []model.User
	users = append(users, user1)
	users = append(users, user2)
	users = append(users, user3)

	//ganti line 116 - 131 sama function buat get all user data

	//INI MASIH SEMENTARA//
	var book1 model.Book
	book1.Title = "Daun yang jatuh tak pernah membenci angin"
	book1.Author = "Tere Liye"

	var book2 model.Book
	book2.Title = "Laut Bercerita"
	book2.Author = "Leila S Chudori"

	var book3 model.Book
	book3.Title = "Please Look After Mom"
	book3.Author = "Orang korea tp lupa namanya"

	var books []model.Book
	books = append(books, book1)
	books = append(books, book2)
	books = append(books, book3)

	//Kode dari line 136 sampe 151 harus diganti sama function

	// scheduler := gocron.NewScheduler()
	// scheduler.Every(1).Week().Do(func() {
	// 	for i := 0; i < len(users); i++ {
	// 		go SendWeeklyEmail(users[i].Email)
	// 	}
	// })
	// <-scheduler.Start()

	for i := 0; i < len(users); i++ {
		SetScheduler(users[i].Email)
	}
}
