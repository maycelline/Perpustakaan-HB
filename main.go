package main

import (
	controllers "Tools/controllers"
	"fmt"
)

func main() {

	fmt.Println("asu")
	var data controllers.DataBorrowed
	data.UserName = "Maycelline"
	data.CourierName = "Dadang Sudrajat"
	data.OrderDate = "20 Mei 2021"
	data.Time = "19.00"

	controllers.SendBorrowAcceptedEmail("maycelinesudarsono@gmail.com", data)

}
