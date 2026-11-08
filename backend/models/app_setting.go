package models

import "time"

type AppSetting struct {
	Key       string    `gorm:"type:text;primaryKey"`
	Value     string    `gorm:"type:text;not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (AppSetting) TableName() string { return "app_settings" }
