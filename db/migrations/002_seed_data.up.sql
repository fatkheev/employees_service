INSERT INTO new_schema.companies (name) VALUES ('Company A'), ('Company B');

INSERT INTO new_schema.departments (name, phone, company_id) VALUES 
('HR', '123-456-7890', 1),
('IT', '987-654-3210', 1),
('HR', '123-456-7890', 2),
('IT', '987-654-3210', 2);

INSERT INTO new_schema.employees (name, surname, phone, company_id, department_id) VALUES 
('John', 'Doe', '123456789', 1, 1),
('Jane', 'Doe', '987654321', 1, 2),
('Jim', 'Beam', '555555555', 2, 3),
('Jack', 'Daniels', '444444444', 2, 4);

INSERT INTO new_schema.passports (type, number, employee_id) VALUES 
('Type1', 'A1234567', 1),
('Type2', 'B2345678', 2),
('Type1', 'C3456789', 3),
('Type2', 'D4567890', 4);
