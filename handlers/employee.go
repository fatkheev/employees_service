package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"employees_service/db"
	"employees_service/models"

	"github.com/gorilla/mux"
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

    // Insert passport information
    if employee.Passport != nil {
        query = `
            INSERT INTO new_schema.passports (type, number, employee_id)
            VALUES ($1, $2, $3)
        `
        _, err = db.DB.Exec(query,
            employee.Passport.Type, employee.Passport.Number, id,
        )
        if err != nil {
            log.Println("Error executing query for passport:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    log.Printf("Successfully created employee with ID: %d", id)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Println("Invalid employee ID:", vars["id"])
        http.Error(w, "Invalid employee ID", http.StatusBadRequest)
        return
    }

    log.Printf("Received request to delete employee with ID: %d", id)

    query := `
        DELETE FROM new_schema.employees
        WHERE id = $1
    `
    log.Printf("Executing query: %s", query)

    result, err := db.DB.Exec(query, id)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Println("Error getting rows affected:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        log.Printf("No employee found with ID: %d", id)
        http.Error(w, "Employee not found", http.StatusNotFound)
        return
    }

    log.Printf("Successfully deleted employee with ID: %d", id)
    w.WriteHeader(http.StatusNoContent)
}

func GetEmployeesByCompany(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    companyID, err := strconv.Atoi(vars["company_id"])
    if err != nil {
        log.Println("Invalid company ID:", vars["company_id"])
        http.Error(w, "Invalid company ID", http.StatusBadRequest)
        return
    }

    log.Printf("Received request to get employees for company ID: %d", companyID)

    query := `
        SELECT e.id, e.name, e.surname, e.phone, e.company_id, e.department_id,
               p.type, p.number,
               d.name, d.phone
        FROM new_schema.employees e
        LEFT JOIN new_schema.passports p ON e.id = p.employee_id
        JOIN new_schema.departments d ON e.department_id = d.id
        WHERE e.company_id = $1
    `
    rows, err := db.DB.Query(query, companyID)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    employees := []models.Employee{}
    for rows.Next() {
        var employee models.Employee
        var passport models.Passport
        var department models.Department
        if err := rows.Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.Phone, &employee.CompanyID, &employee.DepartmentID,
            &passport.Type, &passport.Number,
            &department.Name, &department.Phone); err != nil {
            log.Println("Error scanning row:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employee.Passport = &passport
        employee.Department = &department
        employees = append(employees, employee)
    }

    if err := rows.Err(); err != nil {
        log.Println("Error after scanning rows:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Found %d employees for company ID: %d", len(employees), companyID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(employees)
}

func GetEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    departmentID, err := strconv.Atoi(vars["department_id"])
    if err != nil {
        log.Println("Invalid department ID:", vars["department_id"])
        http.Error(w, "Invalid department ID", http.StatusBadRequest)
        return
    }

    log.Printf("Received request to get employees for department ID: %d", departmentID)

    query := `
        SELECT e.id, e.name, e.surname, e.phone, e.company_id, e.department_id,
               p.type, p.number,
               d.name, d.phone
        FROM new_schema.employees e
        LEFT JOIN new_schema.passports p ON e.id = p.employee_id
        JOIN new_schema.departments d ON e.department_id = d.id
        WHERE e.department_id = $1
    `
    rows, err := db.DB.Query(query, departmentID)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    employees := []models.Employee{}
    for rows.Next() {
        var employee models.Employee
        var passport models.Passport
        var department models.Department
        if err := rows.Scan(&employee.ID, &employee.Name, &employee.Surname, &employee.Phone, &employee.CompanyID, &employee.DepartmentID,
            &passport.Type, &passport.Number,
            &department.Name, &department.Phone); err != nil {
            log.Println("Error scanning row:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employee.Passport = &passport
        employee.Department = &department
        employees = append(employees, employee)
    }

    if err := rows.Err(); err != nil {
        log.Println("Error after scanning rows:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Found %d employees for department ID: %d", len(employees), departmentID)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(employees)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Println("Invalid employee ID:", vars["id"])
        http.Error(w, "Invalid employee ID", http.StatusBadRequest)
        return
    }

    var employee models.Employee
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received request to update employee with ID: %d, data: %+v", id, employee)

    fields := []string{}
    values := []interface{}{}
    i := 1

    if employee.Name != "" {
        fields = append(fields, "name = $"+strconv.Itoa(i))
        values = append(values, employee.Name)
        i++
    }
    if employee.Surname != "" {
        fields = append(fields, "surname = $"+strconv.Itoa(i))
        values = append(values, employee.Surname)
        i++
    }
    if employee.Phone != "" {
        fields = append(fields, "phone = $"+strconv.Itoa(i))
        values = append(values, employee.Phone)
        i++
    }
    if employee.CompanyID != 0 {
        fields = append(fields, "company_id = $"+strconv.Itoa(i))
        values = append(values, employee.CompanyID)
        i++
    }
    if employee.DepartmentID != 0 {
        fields = append(fields, "department_id = $"+strconv.Itoa(i))
        values = append(values, employee.DepartmentID)
        i++
    }

    if len(fields) == 0 {
        log.Println("No fields to update")
        http.Error(w, "No fields to update", http.StatusBadRequest)
        return
    }

    values = append(values, id)
    query := `
        UPDATE new_schema.employees
        SET ` + strings.Join(fields, ", ") + `
        WHERE id = $` + strconv.Itoa(i)

    log.Printf("Executing query: %s with values %v", query, values)

    _, err = db.DB.Exec(query, values...)
    if err != nil {
        log.Println("Error executing query:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Successfully updated employee with ID: %d", id)
    w.WriteHeader(http.StatusNoContent)
}
