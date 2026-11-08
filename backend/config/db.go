package config

import (
	"fmt"
	"hrm-backend/models"
	"hrm-backend/utils"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" {
		host = "postgres_db"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "password"
	}
	if dbname == "" {
		dbname = "hrm_db"
	}
	if port == "" {
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Colombo", host, user, password, dbname, port)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:        time.Second,
			LogLevel:             logger.Warn,
			ParameterizedQueries: true,
		}),
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database! %v", err))
	}

	fmt.Println("Database connection successfully established!")

	// Avoid altering existing legacy `employees` table (created via init.sql).
	// Ensure the attendance and leave request tables exist without forcing changes to employees.
	if !database.Migrator().HasTable(&models.Attendance{}) {
		if err := database.Migrator().CreateTable(&models.Attendance{}); err != nil {
			panic(fmt.Sprintf("Failed to create attendances table: %v", err))
		}
		fmt.Println("Created attendances table")
	}

	if !database.Migrator().HasTable(&models.LeaveRequest{}) {
		if err := database.Migrator().CreateTable(&models.LeaveRequest{}); err != nil {
			panic(fmt.Sprintf("Failed to create leave_requests table: %v", err))
		}
		fmt.Println("Created leave_requests table")
	}

	if !database.Migrator().HasTable(&models.Department{}) {
		if err := database.Migrator().CreateTable(&models.Department{}); err != nil {
			panic(fmt.Sprintf("Failed to create departments table: %v", err))
		}
		fmt.Println("Created departments table")
	}

	if !database.Migrator().HasTable(&models.Payroll{}) {
		if err := database.Migrator().CreateTable(&models.Payroll{}); err != nil {
			panic(fmt.Sprintf("Failed to create payroll table: %v", err))
		}
	}
	if err := ensureEmployeeSchema(database); err != nil {
		panic(fmt.Sprintf("Failed to migrate employees table: %v", err))
	}
	if err := ensureEmployeeRecordLinks(database); err != nil {
		panic(fmt.Sprintf("Failed to migrate employee record links: %v", err))
	}

	if err := database.AutoMigrate(&models.AppSetting{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate app_settings table: %v", err))
	}

	DB = database
}

func ensureEmployeeRecordLinks(database *gorm.DB) error {
	for _, table := range []string{"attendances", "leave_requests"} {
		if !database.Migrator().HasColumn(table, "employee_id") {
			if err := database.Exec("ALTER TABLE " + table + " ADD COLUMN employee_id TEXT").Error; err != nil {
				return err
			}
		}
		update := fmt.Sprintf(`UPDATE %s AS record
			SET employee_id = (
				SELECT e.employee_id FROM employees e
				WHERE e.deleted_at IS NULL
					AND LOWER(TRIM(e.name)) = LOWER(TRIM(record.employee))
				ORDER BY e.id
				LIMIT 1
			)
			WHERE (record.employee_id IS NULL OR record.employee_id = '')
				AND EXISTS (
					SELECT 1 FROM employees e
					WHERE e.deleted_at IS NULL
						AND LOWER(TRIM(e.name)) = LOWER(TRIM(record.employee))
				)`, table)
		if err := database.Exec(update).Error; err != nil {
			return err
		}
	}
	return nil
}

func ensureEmployeeSchema(database *gorm.DB) error {
	if !database.Migrator().HasTable(&models.Employee{}) {
		return database.Migrator().CreateTable(&models.Employee{})
	}

	if !database.Migrator().HasColumn(&models.Employee{}, "EmployeeID") {
		if err := database.Exec("ALTER TABLE employees ADD COLUMN employee_id VARCHAR(32)").Error; err != nil {
			return err
		}
	}
	if err := database.Exec("UPDATE employees SET employee_id = 'EMP' || LPAD(id::text, 3, '0') WHERE employee_id IS NULL OR employee_id = ''").Error; err != nil {
		return err
	}
	if !database.Migrator().HasColumn(&models.Employee{}, "Status") {
		if err := database.Exec("ALTER TABLE employees ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'Active'").Error; err != nil {
			return err
		}
	} else if err := database.Exec("UPDATE employees SET status = 'Active' WHERE status IS NULL OR status = ''").Error; err != nil {
		return err
	}
	if !database.Migrator().HasColumn(&models.Employee{}, "Role") {
		if err := database.Exec("ALTER TABLE employees ADD COLUMN role TEXT").Error; err != nil {
			return err
		}
	}
	if err := database.Exec("UPDATE employees SET role = 'Admin' WHERE employee_id = 'EMP001' AND (role IS NULL OR role = '')").Error; err != nil {
		return err
	}
	if err := database.Exec("UPDATE employees SET role = 'Employee' WHERE role IS NULL OR role = ''").Error; err != nil {
		return err
	}
	if err := database.AutoMigrate(&models.Employee{}); err != nil {
		return err
	}
	if err := migrateEmployeePasswords(database); err != nil {
		return err
	}
	if err := database.Exec("ALTER TABLE employees ALTER COLUMN employee_id SET NOT NULL").Error; err != nil {
		return err
	}
	if err := database.Exec("ALTER TABLE employees ALTER COLUMN role SET DEFAULT 'Employee'").Error; err != nil {
		return err
	}
	return database.Exec("ALTER TABLE employees ALTER COLUMN role SET NOT NULL").Error
}

func migrateEmployeePasswords(database *gorm.DB) error {
	var employees []struct {
		ID       uint
		Password string
	}
	if err := database.Model(&models.Employee{}).Select("id", "password").Find(&employees).Error; err != nil {
		return err
	}
	for _, employee := range employees {
		if strings.HasPrefix(employee.Password, "$2a$") ||
			strings.HasPrefix(employee.Password, "$2b$") ||
			strings.HasPrefix(employee.Password, "$2y$") {
			continue
		}
		hashed, err := utils.HashPassword(employee.Password)
		if err != nil {
			return fmt.Errorf("unable to migrate employee password hashes: %w", err)
		}
		if err := database.Model(&models.Employee{}).Where("id = ?", employee.ID).Update("password", hashed).Error; err != nil {
			return fmt.Errorf("unable to migrate employee password hashes: %w", err)
		}
	}
	return nil
}
