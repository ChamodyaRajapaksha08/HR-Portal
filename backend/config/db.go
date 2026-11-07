package config

import (
	"fmt"
	"hrm-backend/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

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

	if err := ensureEmployeeSchema(database); err != nil {
		panic(fmt.Sprintf("Failed to migrate employees table: %v", err))
	}

	DB = database
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

	if err := database.AutoMigrate(&models.Employee{}); err != nil {
		return err
	}
	return database.Exec("ALTER TABLE employees ALTER COLUMN employee_id SET NOT NULL").Error
}
