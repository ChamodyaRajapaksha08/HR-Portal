-- දැනට තියෙන DB එකේ mock data අයින් කරන script එක.
-- කලින් BACKUP එකක් ගන්න: docker exec postgres_db pg_dump -U postgres hrm_db > backup.sql
BEGIN;
TRUNCATE TABLE payroll, leaves, leave_requests, attendances, departments, employees RESTART IDENTITY CASCADE;
-- Real admin user (password වෙනස් කරන්න!)
INSERT INTO employees (employee_id, name, email, password, department, designation, status, created_at, updated_at)
VALUES ('EMP001', 'Admin123', 'admin123@gmail.com', 'CHANGE_ME_NOW', 'HR', 'HR Manager', 'Active', NOW(), NOW());
COMMIT;


hrmadmin123-aluth pw ek
