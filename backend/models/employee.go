package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	EmployeeID string `gorm:"unique;not null" json:"employee_id"`
	FullName   string `gorm:"column:name;not null" json:"full_name"`
	Email      string `gorm:"unique;not null" json:"email"`
	Password   string `gorm:"not null" json:"password,omitempty"`
	Department string `json:"department"`
	Position   string `gorm:"column:designation" json:"position"`
	Status     string `gorm:"default:'Active'" json:"status"`
}

// TableName specifies the table name for Employee
func (Employee) TableName() string {
	return "employees"
}
