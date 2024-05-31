INSERT INTO new_schema.companies (name) VALUES ('Company A'), ('Company B');

INSERT INTO new_schema.departments (name, phone, company_id) VALUES 
('HR', '123-456-7890', 1), 
('Engineering', '234-567-8901', 1),
('Sales', '345-678-9012', 2);

INSERT INTO new_schema.employees (name, surname, phone, company_id, department_id) VALUES 
('John', 'Doe', '123456789', 1, 1),
('Jane', 'Smith', '987654321', 1, 2),
('Robert', 'Brown', '555555555', 2, 3);

INSERT INTO new_schema.passports (type, number, employee_id) VALUES 
('Passport', 'A12345678', 1),
('Passport', 'B98765432', 2),
('Passport', 'C55555555', 3);
