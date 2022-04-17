package controllers

import (
	"bytes"
	"fmt"

	"Perpustakaan-HB/model"
	"text/template"

	"gopkg.in/gomail.v2"
)

func SendRegisterEmail(user model.User) {
	mail := gomail.NewMessage()

	template := "assets/email_template/register.html"

	result, _ := parseTemplate(template, user)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", user.Email)
	mail.SetHeader("Subject", "Register Success!")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	}
}

func SendBorrowAcceptedEmail(data model.BorrowDataHTML) {
	mail := gomail.NewMessage()

	template := "assets/email_template/accept_borrow_book.html"

	result, _ := parseTemplate(template, data)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", data.User.Email)
	mail.SetHeader("Subject", "Accepted Order")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	}
}

func SendReturnAcceptedEmail(data model.BorrowDataHTML) {
	mail := gomail.NewMessage()

	template := "assets/email_template/accept_return_book.html"

	result, _ := parseTemplate(template, data)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", data.User.Email)
	mail.SetHeader("Subject", "Return Accepted")
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

func SendOverdueEmail(overdueInfo model.UserBorrowing) {
	mail := gomail.NewMessage()
	template := "assets/email_template/overdue_info.html"

	//popular book diinisiasi disini agar user mendapat popular book terbaru
	// var popularBooks model.PopularBooksEmail
	// popularBooks.Books = PopularBooks()

	result, _ := parseTemplate(template, overdueInfo)

	mail.SetHeader("From", "perpushb@gmail.com")
	mail.SetHeader("To", overdueInfo.UserData.Email)
	mail.SetHeader("Subject", "Overdue Borrowed")
	mail.SetBody("text/html", result)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "perpushb@gmail.com", "PerpusHBH1tZ")

	if err := sender.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email sent to " + overdueInfo.UserData.Email)
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
		query := "SELECT borrowId, borrowDate FROM borrows WHERE memberId = ? HAVING (TIMESTAMPDIFF(WEEK, borrowDate, CURDATE()) >= 2)"

		var book model.Book
		var books []model.Book
		var borrow model.Borrowing
		var borrows []model.Borrowing

		rows, err := db.Query(query, users[i].ID)
		if err != nil {
			return nil, false
		}

		for rows.Next() {
			if err := rows.Scan(&borrow.ID, &borrow.BorrowDate); err != nil {
				return nil, false
			} else {
				var stockId int
				query := "SELECT borrowslist.stockId FROM borrowslist WHERE borrowslist.borrowId = ? AND (borrowslist.borrowState = 'BORROWED' OR borrowslist.borrowState = 'OVERDUE')"
				rows2, err := db.Query(query, borrow.ID)

				if err != nil {
					return nil, false
				}

				for rows2.Next() {
					if err := rows2.Scan(&stockId); err != nil {
						return nil, false
					} else {
						fmt.Println(users[i])
						fmt.Println("Masuk sini 5")
						fmt.Println(stockId)
						query := "SELECT books.bookId, books.bookTitle, books.author FROM books JOIN stocks ON books.bookId = stocks.bookId WHERE stocks.stockId = ?"
						rows3 := db.QueryRow(query, stockId)
						rows3.Scan(&book.ID, &book.Title, &book.Author)
						_, err := db.Exec("UPDATE borrowslist SET borrowState = 'OVERDUE' WHERE borrowId = ? AND stockId = ?", borrow.ID, stockId)
						if err != nil {
							fmt.Println(err)
							return nil, false
						}
						books = append(books, book)
						fmt.Println(books)
					}
				}
				if books != nil {
					borrow.Book = books
					borrows = append(borrows, borrow)
				}
			}
		}
		if borrows != nil {
			userBorrow.UserData = users[i]
			userBorrow.UserBorrowingData = borrows
			usersBorrow = append(usersBorrow, userBorrow)
		}
	}

	fmt.Println(usersBorrow)

	if len(usersBorrow) < 1 {
		return nil, false
	} else {
		return usersBorrow, true
	}
}
