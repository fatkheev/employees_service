package models

type Passport struct {
    Type   string `json:"type"`
    Number string `json:"number"`
}

type Department struct {
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

type Employee struct {
    ID          int        `json:"id"`
    Name        string     `json:"name"`
    Surname     string     `json:"surname"`
    Phone       string     `json:"phone"`
    CompanyID   int        `json:"company_id"`
    DepartmentID int       `json:"department_id"`
    Passport    *Passport  `json:"passport,omitempty"`
    Department  *Department `json:"department,omitempty"`
}
