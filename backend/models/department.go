package models

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name      string `gorm:"not null" json:"name"`
	Head      string `json:"head"`
	Employees int    `json:"employees"`
	Status    string `gorm:"default:'Active'" json:"status"`
}

func (Department) TableName() string {
	return "departments"
}
