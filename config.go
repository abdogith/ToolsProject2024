package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Global DB variable
var DB *sql.DB

func ConnectDB() {
	dsn := "root:abdomysql2001@tcp(mysql-container:3306)/userdb"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
		return
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("Database connection failed: ", err)
		return
	}

	log.Println("Database connected successfully!")
}
