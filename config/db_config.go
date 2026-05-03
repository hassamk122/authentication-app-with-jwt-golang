package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConnectToDB(dbUrl string) *sql.DB {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error connecting to the database %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Database connection failed. Reason: %v", err)
	}

	log.Println("Connected to database")

	return db
}
