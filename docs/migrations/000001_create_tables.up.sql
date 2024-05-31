CREATE SCHEMA IF NOT EXISTS new_schema;

CREATE TABLE new_schema.companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE new_schema.departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    company_id INT REFERENCES new_schema.companies(id) ON DELETE CASCADE
);

CREATE TABLE new_schema.employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    company_id INT REFERENCES new_schema.companies(id) ON DELETE CASCADE,
    department_id INT REFERENCES new_schema.departments(id) ON DELETE CASCADE
);

CREATE TABLE new_schema.passports (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50),
    number VARCHAR(50),
    employee_id INT REFERENCES new_schema.employees(id) ON DELETE CASCADE
);
