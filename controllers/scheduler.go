package controllers

import "github.com/jasonlvhit/gocron"

func SetScheduler(email string) {
	scheduler := gocron.NewScheduler()
	scheduler.Every(1).Week().Do(func() {
		SendWeeklyEmail(email)
	})
	<-scheduler.Start()
}
