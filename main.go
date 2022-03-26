package main

import (
	"fmt"

	"github.com/jasonlvhit/gocron"
	_ "github.com/jasonlvhit/gocron"
)

func main() {
	count := 0
	gocron.Every(5).Do(func() {
		fmt.Println("monyet ", count)
		count++
	})
	<-gocron.Start()
}
