package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	EmployeeID string `gorm:"unique;not null" json:"employee_id"`
	FullName   string `gorm:"column:name;not null" json:"full_name"`
	Email      string `gorm:"unique;not null" json:"email"`
	Password   string `gorm:"not null" json:"-"`
	Department string `json:"department"`
	Position   string `gorm:"column:designation" json:"position"`
	Status     string `gorm:"default:'Active'" json:"status"`
	Role       string `gorm:"type:text;default:'Employee'" json:"role"`
}

// TableName specifies the table name for Employee
func (Employee) TableName() string {
	return "employees"
}
