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

func SendRetunAcceptedEmail(destinationAddress string, data model.BorrowDataHTML) {
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

	//popular book diinisiasi disini agar user mendapat popular book terbaru
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
	var users []model.User = GetAllUsers()

	for i := 0; i < len(users); i++ {
		SetScheduler(users[i].Email)
	}
}
