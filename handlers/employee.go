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
		log.Println("Ошибка при декодировании тела запроса:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на создание сотрудника: %+v", employee)

	query := `
		INSERT INTO new_schema.employees (name, surname, phone, company_id, department_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	log.Printf("Выполнение запроса: %s", query)

	var id int
	err := db.DB.QueryRow(query,
		employee.Name, employee.Surname, employee.Phone, employee.CompanyID, employee.DepartmentID,
	).Scan(&id)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Вставка информации о паспорте
	if employee.Passport != nil {
		query = `
			INSERT INTO new_schema.passports (type, number, employee_id)
			VALUES ($1, $2, $3)
		`
		_, err = db.DB.Exec(query,
			employee.Passport.Type, employee.Passport.Number, id,
		)
		if err != nil {
			log.Println("Ошибка при выполнении запроса для паспорта:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Сотрудник успешно создан с ID: %d", id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Недействительный ID сотрудника:", vars["id"])
		http.Error(w, "Недействительный ID сотрудника", http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на удаление сотрудника с ID: %d", id)

	query := `
		DELETE FROM new_schema.employees
		WHERE id = $1
	`
	log.Printf("Выполнение запроса: %s", query)

	result, err := db.DB.Exec(query, id)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Ошибка при получении затронутых строк:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Printf("Сотрудник с ID %d не найден", id)
		http.Error(w, "Сотрудник не найден", http.StatusNotFound)
		return
	}

	log.Printf("Сотрудник с ID %d успешно удален", id)
	w.WriteHeader(http.StatusNoContent)
}

func GetEmployeesByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID, err := strconv.Atoi(vars["company_id"])
	if err != nil {
		log.Println("Недействительный ID компании:", vars["company_id"])
		http.Error(w, "Недействительный ID компании", http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на получение сотрудников для компании с ID: %d", companyID)

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
		log.Println("Ошибка при выполнении запроса:", err)
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
			log.Println("Ошибка при сканировании строки:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employee.Passport = &passport
		employee.Department = &department
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка после сканирования строк:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено %d сотрудников для компании с ID: %d", len(employees), companyID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func GetEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	departmentID, err := strconv.Atoi(vars["department_id"])
	if err != nil {
		log.Println("Недействительный ID отдела:", vars["department_id"])
		http.Error(w, "Недействительный ID отдела", http.StatusBadRequest)
		return
	}

	log.Printf("Получен запрос на получение сотрудников для отдела с ID: %d", departmentID)

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
		log.Println("Ошибка при выполнении запроса:", err)
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
			log.Println("Ошибка при сканировании строки:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employee.Passport = &passport
		employee.Department = &department
		employees = append(employees, employee)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка после сканирования строк:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Найдено %d сотрудников для отдела с ID: %d", len(employees), departmentID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Println("Недействительный ID сотрудника:", vars["id"])
        http.Error(w, "Недействительный ID сотрудника", http.StatusBadRequest)
        return
    }

    var employee models.Employee
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        log.Println("Ошибка при декодировании тела запроса:", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Получен запрос на обновление сотрудника с ID: %d, данные: %+v", id, employee)

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
        log.Println("Нет полей для обновления")
        http.Error(w, "Нет полей для обновления", http.StatusBadRequest)
        return
    }

    values = append(values, id)
    query := `
        UPDATE new_schema.employees
        SET ` + strings.Join(fields, ", ") + `
        WHERE id = $` + strconv.Itoa(i)

    log.Printf("Выполнение запроса: %s с значениями %v", query, values)

    _, err = db.DB.Exec(query, values...)
    if err != nil {
        log.Println("Ошибка при выполнении запроса:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    log.Printf("Сотрудник с ID %d успешно обновлен", id)
    w.WriteHeader(http.StatusNoContent)
}
