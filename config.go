package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("root"), os.Getenv("abdomysql2001"), os.Getenv("localhost"), os.Getenv("abdulrahman"))
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}
	if err := database.Ping(); err != nil {
		log.Fatal("Could not ping database: ", err)
	}
	DB = database
	log.Println("Database connected successfully!")
}
