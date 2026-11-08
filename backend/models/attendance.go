package models

import (
	"time"

	"gorm.io/gorm"
)

type Attendance struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	EmployeeID string         `gorm:"type:text" json:"employee_id"`
	Employee   string         `gorm:"not null" json:"employee"`
	Department string         `json:"department"`
	CheckIn    string         `json:"check_in"`
	CheckOut   string         `json:"check_out"`
	Status     string         `json:"status"`
}
