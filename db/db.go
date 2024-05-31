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
    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
    log.Printf("Connecting to database with connection string: %s", connStr)
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Error opening database connection:", err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatal("Error pinging database:", err)
    }

    log.Println("Successfully connected to database")
}

func SetDB(mockDB *sql.DB) {
    DB = mockDB
}