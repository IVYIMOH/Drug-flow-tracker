package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := "user=postgres password=yourpassword dbname=drugflow sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("DB unreachable:", err)

	}
	log.Println("DB connection successful")
}
