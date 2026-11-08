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

-- 4. පළමු Admin user (HR Manager). Mock data නෑ.
-- !! password එක දාන්න කලින් වෙනස් කරන්න. Login එකෙන් පස්සේ system එකෙන්ම මාරු කරගන්න.
INSERT INTO employees (employee_id, name, email, password, department, designation, status, created_at, updated_at)
VALUES ('EMP001', 'Admin123', 'admin123@gmail.com', 'CHANGE_ME_NOW', 'HR', 'HR Manager', 'Active', NOW(), NOW())
ON CONFLICT DO NOTHING;
