package controllers

import (
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func SetEmailWeeklyScheduler(email string) {
	scheduler := gocron.NewScheduler()
	scheduler.Every(1).Week().Do(func() {
		SendWeeklyEmail(email)
	})
	scheduler.Start()
}

func SetEmailBorrowingInfoScheduler() {
	var schedulerOverdue = gocron.NewScheduler()
	schedulerOverdue.Every(1).Day().Do(func() {
		fmt.Println("Masuk sini 1")
		userBorrowsData, status := CheckUserBorrowing()
		if status {
			for i := 0; i < len(userBorrowsData); i++ {
				fmt.Println("Masuk sini 2")
				SendOverdueEmail(userBorrowsData[i])
			}
		}
	})
	schedulerOverdue.Start()
}
