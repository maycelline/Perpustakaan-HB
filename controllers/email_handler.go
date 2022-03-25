package controllers

import (
	"bytes"
	"fmt"

	"gopkg.in/gomail.v2"

	"text/template"
)

func SendRegisterEmail(destinationAddress string, user User) {
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

func SendBorrowAcceptedEmail(destinationAddress string, data DataBorrowed) {
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

func SendBorrowRejectedEmail(destinationAddress string, data DataBorrowed) {
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
