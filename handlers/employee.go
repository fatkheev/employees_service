package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "your_project/db"
    "your_project/models"
)

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
    var employee models.Employee
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var id int
    err := db.DB.QueryRow(
        "INSERT INTO employees (name, surname, phone, company_id, department_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
        employee.Name, employee.Surname, employee.Phone, employee.CompanyID, employee.DepartmentID,
    ).Scan(&id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Implement other handlers similarly...
