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
		SetEmailWeeklyScheduler(users[i].Email)
	}
}

func CheckUserBorrowing() ([]model.UserBorrowing, bool) {
	db := connect()
	defer db.Close()

	var users []model.User = GetAllUsers()
	var userBorrow model.UserBorrowing
	var usersBorrow []model.UserBorrowing
	for i := 0; i < len(users); i++ {
		query := "SELECT borrowId FROM borrows WHERE borrows.memberId = ? HAVING (TIMESTAMPDIFF(WEEK, borrows.borrowDate, CURDATE()) >= 2)"

		var book model.Book
		var books []model.Book
		var borrow model.Borrowing
		var borrows []model.Borrowing

		rows, err := db.Query(query, users[i].ID)

		if err != nil {
			return nil, false
		}

		for rows.Next() {
			if err := rows.Scan(); err != nil {
				return nil, false
			} else {
				var stockId int
				query := "SELECT * FROM borrowslist WHERE borrowslist.borrowId = ? AND borrowslist.borrowState = 'BORROWED' OR borrowslist.borrowState = 'OVERDUE'"
				rows2, err := db.Query(query, borrow.ID)

				if err != nil {
					return nil, false
				}

				for rows2.Next() {
					if err := rows2.Scan(&borrow.ID, &stockId); err != nil {
						return nil, false
					} else {
						query := "SELECT books.* FROM books JOIN stocks WHERE stocks.stockId = ?"
						rows3 := db.QueryRow(query, stockId)
						rows3.Scan()
						_, err := db.Exec("UPDATE borrowslist SET borrowState = 'OVERDUE' WHERE borrowId = ? AND bookId = ?", borrow.ID, book.ID)
						if err != nil {
							return nil, false
						}

						books = append(books, book)
					}
				}
				borrow.Book = books
				borrows = append(borrows, borrow)
			}
		}
		userBorrow.UserData = users[i]
		userBorrow.UserBorrowingData = borrows
		usersBorrow = append(usersBorrow, userBorrow)
	}

	if len(usersBorrow) < 1 {
		return nil, false
	} else {
		return usersBorrow, true
	}
}
