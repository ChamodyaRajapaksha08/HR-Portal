package models

import "time"

type Payroll struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	EmployeeID  uint       `gorm:"column:employee_id" json:"-"`
	BasicSalary float64    `gorm:"column:basic_salary" json:"basic_salary"`
	Allowances  float64    `json:"allowances"`
	Deductions  float64    `json:"deductions"`
	NetSalary   float64    `gorm:"column:net_salary" json:"net_salary"`
	PaidDate    *time.Time `gorm:"column:paid_date" json:"paid_date"`
}

func (Payroll) TableName() string { return "payroll" }
