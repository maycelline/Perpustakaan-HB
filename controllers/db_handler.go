package controllers

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ = godotenv.Load()
var db_port = os.Getenv("DB_PORT")
var dataSource = "admin:admin@tcp(localhost:" + db_port + ")/perpustakaanhb?parseTime=true&loc=Asia%2FJakarta"

func Connect() *sql.DB {
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectGorm() *gorm.DB {
	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
