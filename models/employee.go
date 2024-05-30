package models

type Employee struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Surname     string `json:"surname"`
    Phone       string `json:"phone"`
    CompanyID   int    `json:"company_id"`
    DepartmentID int   `json:"department_id"`
}