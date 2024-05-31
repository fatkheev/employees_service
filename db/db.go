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
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), "postgres")
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

    // Проверка наличия схемы и установка search_path
    var schemaExists bool
    err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.schemata WHERE schema_name = 'new_schema')").Scan(&schemaExists)
    if err != nil {
        log.Fatal("Error checking if schema exists:", err)
    }

    if !schemaExists {
        log.Fatal("Schema 'new_schema' does not exist")
    }

    _, err = DB.Exec("SET search_path TO new_schema")
    if err != nil {
        log.Fatal("Error setting search_path:", err)
    }

    log.Println("Search path set to new_schema")

    // Проверка текущей базы данных и схемы
    var dbName, schemaName string
    err = DB.QueryRow("SELECT current_database(), current_schema()").Scan(&dbName, &schemaName)
    if err != nil {
        log.Fatal("Error querying current database and schema:", err)
    }

    log.Printf("Connected to database: %s, schema: %s", dbName, schemaName)
}

func SetDB(mockDB *sql.DB) {
    DB = mockDB
}
