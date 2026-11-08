package models

import "gorm.io/gorm"

type LeaveRequest struct {
	gorm.Model
	EmployeeID string `gorm:"type:text" json:"employee_id"`
	Employee   string `json:"employee"`
	LeaveType  string `json:"leave_type"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Days       int    `json:"days"`
	Reason     string `json:"reason"`
	Status     string `json:"status"`
	AppliedOn  string `json:"applied_on"`
}
