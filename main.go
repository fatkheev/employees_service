package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "your_project/handlers"
    "your_project/db"
)

func main() {
    db.InitDB()

    r := mux.NewRouter()
    r.HandleFunc("/employees", handlers.CreateEmployee).Methods("POST")
    r.HandleFunc("/employees/{id}", handlers.DeleteEmployee).Methods("DELETE")
    r.HandleFunc("/companies/{id}/employees", handlers.GetEmployeesByCompany).Methods("GET")
    r.HandleFunc("/departments/{id}/employees", handlers.GetEmployeesByDepartment).Methods("GET")
    r.HandleFunc("/employees/{id}", handlers.UpdateEmployee).Methods("PUT")

    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}