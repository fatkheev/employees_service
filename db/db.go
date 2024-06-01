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
    log.Printf("Подключение к базе данных с использованием строки подключения: %s", connStr)
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Ошибка при открытии подключения к базе данных:", err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatal("Ошибка при проверке соединения с базой данных:", err)
    }

    log.Println("Успешное подключение к базе данных")
}

func SetDB(mockDB *sql.DB) {
    DB = mockDB
}
