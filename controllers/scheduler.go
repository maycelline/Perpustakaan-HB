package controllers

import (
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
		userBorrowsData, status := CheckUserBorrowing()
		if status {
			for i := 0; i < len(userBorrowsData); i++ {
				SendOverdueEmail(userBorrowsData[i])
			}
		}
	})
	schedulerOverdue.Start()
}
