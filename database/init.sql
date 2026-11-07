-- 1. සේවකයින්ගේ තොරතුරු තබා ගන්නා Table එක (Employees)
CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(32) UNIQUE,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    department VARCHAR(50) NOT NULL,
    designation VARCHAR(50) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'Active',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    joined_date DATE DEFAULT CURRENT_DATE
);

-- 2. නිවාඩු විස්තර තබා ගන්නා Table එක (Leaves)
CREATE TABLE IF NOT EXISTS leaves (
    id SERIAL PRIMARY KEY,
    employee_id INT REFERENCES employees(id) ON DELETE CASCADE, -- Shared Identifier [cite: 37]
    leave_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'Pending' -- Pending, Approved, Rejected
);

-- 3. පඩි විස්තර තබා ගන්නා Table එක (Payroll)
CREATE TABLE IF NOT EXISTS payroll (
    id SERIAL PRIMARY KEY,
    employee_id INT REFERENCES employees(id) ON DELETE CASCADE, -- Shared Identifier [cite: 37]
    basic_salary NUMERIC(10, 2) NOT NULL,
    allowances NUMERIC(10, 2) DEFAULT 0.00,
    deductions NUMERIC(10, 2) DEFAULT 0.00,
    net_salary NUMERIC(10, 2) NOT NULL,
    paid_date DATE DEFAULT CURRENT_DATE
);

-- 4. Assignment එකේ ඉල්ලලා තියෙන විදිහට පරීක්ෂා කරලා බලන්න Sample Data (Mock Data) ටිකක් ඇතුළත් කිරීම [cite: 69]
INSERT INTO employees (name, email, password, department, designation) VALUES
('Kavindu Perera', 'kavindu@gmail.com', 'password123', 'IT', 'Software Engineer'),
('Nimesha Fernando', 'nimesha@gmail.com', 'password123', 'HR', 'HR Manager'),
('Tharushi Silva', 'tharushi@gmail.com', 'password123', 'Finance', 'Accountant');

INSERT INTO leaves (employee_id, leave_type, start_date, end_date, status) VALUES
(1, 'Sick Leave', '2026-05-20', '2026-05-22', 'Pending'),
(2, 'Casual Leave', '2026-06-01', '2026-06-03', 'Approved');

INSERT INTO payroll (employee_id, basic_salary, allowances, deductions, net_salary) VALUES
(1, 85000.00, 15000.00, 5000.00, 95000.00),
(2, 95000.00, 20000.00, 7000.00, 108000.00);