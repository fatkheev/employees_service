package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "employees_service/db"
    "employees_service/models"
)

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
    var employee models.Employee
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received request to create employee: %+v", employee)

    query := `
        INSERT INTO new_schema.employees (name, surname, phone, company_id, department_id)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    log.Printf("Executing query: %s", query)

    var id int
    err := db.DB.QueryRow(query,
        employee.Name, employee.Surname, employee.Phone, employee.CompanyID, employee.DepartmentID,
    ).Scan(&id)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Successfully created employee with ID: %d", id)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}
