package database

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Println("LOCAL database")
		connStr = "host=localhost port=5432 user=postgres password=admin dbname=bioskop sslmode=disable"
	} else {
		log.Println("RAILWAY database")

		if !strings.Contains(connStr, "sslmode=") {
			if strings.Contains(connStr, "?") {
				connStr += "&sslmode=require"
			} else {
				connStr += "?sslmode=require"
			}
		}
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database not connected:", err)
	}

	log.Println("Database connected successfully")
}