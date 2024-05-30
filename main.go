package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "employees_service/handlers"
    "employees_service/db"
)

func main() {
    // Загрузка переменных окружения из .env файла
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Логирование загруженных переменных окружения
    log.Printf("DB_HOST: %s, DB_PORT: %s, DB_USER: %s, DB_PASSWORD: %s, DB_NAME: %s",
        os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

    db.InitDB()

    r := mux.NewRouter()
    r.HandleFunc("/employees", handlers.CreateEmployee).Methods("POST")

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
