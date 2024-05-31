package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "github.com/gorilla/mux"
    "employees_service/db"
    "employees_service/models"
)

func TestCreateEmployee(t *testing.T) {
    // Создаем мок базы данных
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer mockDB.Close()

    db.SetDB(mockDB) // Используем мок базы данных

    // Настройка ожидаемых запросов и их результатов
    mock.ExpectExec("SET search_path TO new_schema").
        WillReturnResult(sqlmock.NewResult(1, 1))

    mock.ExpectQuery("INSERT INTO new_schema.employees \\(name, surname, phone, company_id, department_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
        WithArgs("John", "Doe", "123456789", 1, 1).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    mock.ExpectExec("INSERT INTO new_schema.passports \\(type, number, employee_id\\) VALUES \\(\\$1, \\$2, \\$3\\)").
        WithArgs("PassportType", "123456", 1).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Инициализируем базу данных, чтобы установить search_path
    if _, err := mockDB.Exec("SET search_path TO new_schema"); err != nil {
        t.Fatalf("failed to set search_path: %s", err)
    }

    // Создание запроса
    employee := models.Employee{
        Name:         "John",
        Surname:      "Doe",
        Phone:        "123456789",
        CompanyID:    1,
        DepartmentID: 1,
        Passport:     models.Passport{Type: "PassportType", Number: "123456"},
        Department:   models.Department{Name: "IT", Phone: "987654321"},
    }
    body, _ := json.Marshal(employee)
    req, err := http.NewRequest("POST", "/employees", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Создание ResponseRecorder для захвата ответа
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(CreateEmployee)

    // Вызов обработчика
    handler.ServeHTTP(rr, req)

    // Проверка статуса ответа
    assert.Equal(t, http.StatusCreated, rr.Code, "Expected status code 201")

    // Проверка тела ответа
    var response map[string]interface{}
    err = json.NewDecoder(rr.Body).Decode(&response)
    if err != nil {
        t.Fatalf("Failed to decode response body: %s", err)
    }

    id, ok := response["id"].(float64)
    if !ok {
        t.Fatalf("Expected ID to be a number but got %T", response["id"])
    }

    assert.Greater(t, int(id), 0, "ID should be greater than 0")

    // Проверка ожиданий моков
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func TestDeleteEmployee(t *testing.T) {
    // Создаем мок базы данных
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer mockDB.Close()

    db.SetDB(mockDB) // Используем мок базы данных

    // Настройка ожидаемых запросов и их результатов
    mock.ExpectExec("SET search_path TO new_schema").
        WillReturnResult(sqlmock.NewResult(1, 1))

    mock.ExpectExec("DELETE FROM new_schema.employees WHERE id = \\$1").
        WithArgs(1).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Инициализируем базу данных, чтобы установить search_path
    if _, err := mockDB.Exec("SET search_path TO new_schema"); err != nil {
        t.Fatalf("failed to set search_path: %s", err)
    }

    // Создание запроса
    req, err := http.NewRequest("DELETE", "/employees/1", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Создание ResponseRecorder для захвата ответа
    rr := httptest.NewRecorder()
    router := mux.NewRouter()
    router.HandleFunc("/employees/{id}", DeleteEmployee).Methods("DELETE")

    // Вызов обработчика
    router.ServeHTTP(rr, req)

    // Проверка статуса ответа
    assert.Equal(t, http.StatusNoContent, rr.Code, "Expected status code 204")

    // Проверка ожиданий моков
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func TestGetEmployeesByCompany(t *testing.T) {
    // Создаем мок базы данных
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer mockDB.Close()

    db.SetDB(mockDB) // Используем мок базы данных

    // Настройка ожидаемых запросов и их результатов
    mock.ExpectExec("SET search_path TO new_schema").
        WillReturnResult(sqlmock.NewResult(1, 1))

    rows := sqlmock.NewRows([]string{"id", "name", "surname", "phone", "company_id", "type", "number", "department_name", "department_phone"}).
        AddRow(1, "John", "Doe", "123456789", 1, "PassportType", "123456", "IT", "987654321")

    mock.ExpectQuery("SELECT e.id, e.name, e.surname, e.phone, e.company_id, p.type, p.number, d.name, d.phone FROM new_schema.employees e JOIN new_schema.passports p ON e.id = p.employee_id JOIN new_schema.departments d ON e.department_id = d.id WHERE e.company_id = \\$1").
        WithArgs(1).
        WillReturnRows(rows)

    // Инициализируем базу данных, чтобы установить search_path
    if _, err := mockDB.Exec("SET search_path TO new_schema"); err != nil {
        t.Fatalf("failed to set search_path: %s", err)
    }

    // Создание запроса
    req, err := http.NewRequest("GET", "/companies/1/employees", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Создание ResponseRecorder для захвата ответа
    rr := httptest.NewRecorder()
    router := mux.NewRouter()
    router.HandleFunc("/companies/{company_id}/employees", GetEmployeesByCompany).Methods("GET")

    // Вызов обработчика
    router.ServeHTTP(rr, req)

    // Проверка статуса ответа
    assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

    // Проверка тела ответа
    var employees []models.Employee
    err = json.NewDecoder(rr.Body).Decode(&employees)
    if err != nil {
        t.Fatalf("Failed to decode response body: %s", err)
    }

    assert.Len(t, employees, 1, "Expected 1 employee")
    assert.Equal(t, employees[0].Name, "John")
    assert.Equal(t, employees[0].Surname, "Doe")
    assert.Equal(t, employees[0].Phone, "123456789")
    assert.Equal(t, employees[0].CompanyID, 1)
    assert.Equal(t, employees[0].Passport.Type, "PassportType")
    assert.Equal(t, employees[0].Passport.Number, "123456")
    assert.Equal(t, employees[0].Department.Name, "IT")
    assert.Equal(t, employees[0].Department.Phone, "987654321")

    // Проверка ожиданий моков
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}

func TestGetEmployeesByDepartment(t *testing.T) {
	// Создаем мок базы данных
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		 t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.SetDB(mockDB) // Используем мок базы данных

	// Настройка ожидаемых запросов и их результатов
	mock.ExpectExec("SET search_path TO new_schema").
		 WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "name", "surname", "phone", "company_id", "type", "number", "department_name", "department_phone"}).
		 AddRow(1, "John", "Doe", "123456789", 1, "PassportType", "123456", "IT", "987654321")

	mock.ExpectQuery("SELECT e.id, e.name, e.surname, e.phone, e.company_id, p.type, p.number, d.name, d.phone FROM new_schema.employees e JOIN new_schema.passports p ON e.id = p.employee_id JOIN new_schema.departments d ON e.department_id = d.id WHERE e.department_id = \\$1").
		 WithArgs(1).
		 WillReturnRows(rows)

	// Инициализируем базу данных, чтобы установить search_path
	if _, err := mockDB.Exec("SET search_path TO new_schema"); err != nil {
		 t.Fatalf("failed to set search_path: %s", err)
	}

	// Создание запроса
	req, err := http.NewRequest("GET", "/departments/1/employees", nil)
	if err != nil {
		 t.Fatal(err)
	}

	// Создание ResponseRecorder для захвата ответа
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/departments/{department_id}/employees", GetEmployeesByDepartment).Methods("GET")

	// Вызов обработчика
	router.ServeHTTP(rr, req)

	// Проверка статуса ответа
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	// Проверка тела ответа
	var employees []models.Employee
	err = json.NewDecoder(rr.Body).Decode(&employees)
	if err != nil {
		 t.Fatalf("Failed to decode response body: %s", err)
	}

	assert.Len(t, employees, 1, "Expected 1 employee")
	assert.Equal(t, employees[0].Name, "John")
	assert.Equal(t, employees[0].Surname, "Doe")
	assert.Equal(t, employees[0].Phone, "123456789")
	assert.Equal(t, employees[0].CompanyID, 1)
	assert.Equal(t, employees[0].Passport.Type, "PassportType")
	assert.Equal(t, employees[0].Passport.Number, "123456")
	assert.Equal(t, employees[0].Department.Name, "IT")
	assert.Equal(t, employees[0].Department.Phone, "987654321")

	// Проверка ожиданий моков
	if err := mock.ExpectationsWereMet(); err != nil {
		 t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateEmployee(t *testing.T) {
	// Создаем мок базы данных
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		 t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db.SetDB(mockDB) // Используем мок базы данных

	// Настройка ожидаемых запросов и их результатов
	mock.ExpectExec("SET search_path TO new_schema").
		 WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE new_schema.employees SET name = \\$1, surname = \\$2, phone = \\$3, company_id = \\$4, department_id = \\$5 WHERE id = \\$6").
		 WithArgs("John", "Doe", "123456789", 1, 1, 1).
		 WillReturnResult(sqlmock.NewResult(1, 1))

	// Инициализируем базу данных, чтобы установить search_path
	if _, err := mockDB.Exec("SET search_path TO new_schema"); err != nil {
		 t.Fatalf("failed to set search_path: %s", err)
	}

	// Создание запроса
	employee := models.Employee{
		 Name:         "John",
		 Surname:      "Doe",
		 Phone:        "123456789",
		 CompanyID:    1,
		 DepartmentID: 1,
	}
	body, _ := json.Marshal(employee)
	req, err := http.NewRequest("PUT", "/employees/1", bytes.NewBuffer(body))
	if err != nil {
		 t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создание ResponseRecorder для захвата ответа
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/employees/{id}", UpdateEmployee).Methods("PUT")

	// Вызов обработчика
	router.ServeHTTP(rr, req)

	// Проверка статуса ответа
	assert.Equal(t, http.StatusNoContent, rr.Code, "Expected status code 204")

	// Проверка ожиданий моков
	if err := mock.ExpectationsWereMet(); err != nil {
		 t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
