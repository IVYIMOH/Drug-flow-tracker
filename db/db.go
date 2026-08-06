package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var connStr string

	// 1. Check if Railway provided a full DATABASE_URL
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		// 2. Fall back to individual variables for local dev
		host := os.Getenv("DB_HOST")
		if host == "" {
			log.Fatal("DB Config Error: Neither DATABASE_URL nor DB_HOST is set!")
		}
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host,
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("DB unreachable:", err)
	}
	log.Println("DB connection successful")
}