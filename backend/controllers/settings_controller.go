package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"hrm-backend/config"
	"hrm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type companySettings struct {
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
}

type attendanceSettings struct {
	WorkStart     string   `json:"work_start"`
	WorkEnd       string   `json:"work_end"`
	LateAfter     string   `json:"late_after"`
	WorkingDays   []string `json:"working_days"`
	WeeklyOffDays []string `json:"weekly_off_days"`
}

type preferenceSettings struct {
	Currency            string `json:"currency"`
	DateFormat          string `json:"date_format"`
	Timezone            string `json:"timezone"`
	EmailNotifications  bool   `json:"email_notifications"`
}

var settingDefaults = struct {
	Company     companySettings     `json:"company"`
	Attendance  attendanceSettings  `json:"attendance"`
	Preferences preferenceSettings  `json:"preferences"`
}{
	Company: companySettings{},
	Attendance: attendanceSettings{
		WorkingDays:   []string{},
		WeeklyOffDays: []string{},
	},
	Preferences: preferenceSettings{
		Currency: "LKR",
		Timezone: "Asia/Colombo",
	},
}

func readAppSetting(key string, target any) error {
	var setting models.AppSetting
	err := config.DB.First(&setting, "key = ?", key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(setting.Value), target)
}

func writeAppSetting(key string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	setting := models.AppSetting{Key: key, Value: string(encoded), UpdatedAt: time.Now()}
	return config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&setting).Error
}

func GetSettings(c *gin.Context) {
	result := settingDefaults
	if err := readAppSetting("company", &result.Company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to load company settings: %v", err)})
		return
	}
	if err := readAppSetting("attendance", &result.Attendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to load attendance settings: %v", err)})
		return
	}
	if err := readAppSetting("preferences", &result.Preferences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to load preferences: %v", err)})
		return
	}
	if result.Attendance.WorkingDays == nil {
		result.Attendance.WorkingDays = []string{}
	}
	if result.Attendance.WeeklyOffDays == nil {
		result.Attendance.WeeklyOffDays = []string{}
	}
	c.JSON(http.StatusOK, result)
}

func UpdateCompanySettings(c *gin.Context) {
	var input companySettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company settings"})
		return
	}
	input.CompanyName = strings.TrimSpace(input.CompanyName)
	input.Email = strings.TrimSpace(input.Email)
	if input.CompanyName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company name is required"})
		return
	}
	if input.Email != "" {
		if _, err := mail.ParseAddress(input.Email); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid company email"})
			return
		}
	}
	if err := writeAppSetting("company", input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

var clockTimePattern = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

func UpdateAttendanceSettings(c *gin.Context) {
	var input attendanceSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendance settings"})
		return
	}
	for label, value := range map[string]string{
		"Work start": input.WorkStart,
		"Work end":   input.WorkEnd,
		"Late-after": input.LateAfter,
	} {
		if !clockTimePattern.MatchString(value) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%s must be a valid HH:MM time", label)})
			return
		}
	}
	workStart, _ := time.Parse("15:04", input.WorkStart)
	workEnd, _ := time.Parse("15:04", input.WorkEnd)
	if !workEnd.After(workStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Work end must be after work start"})
		return
	}
	if input.WorkingDays == nil {
		input.WorkingDays = []string{}
	}
	if input.WeeklyOffDays == nil {
		input.WeeklyOffDays = []string{}
	}
	if err := writeAppSetting("attendance", input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

func UpdatePreferenceSettings(c *gin.Context) {
	var input preferenceSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid preference settings"})
		return
	}
	input.Currency = strings.TrimSpace(input.Currency)
	if input.Currency == "" {
		input.Currency = "LKR"
	}
	input.Timezone = strings.TrimSpace(input.Timezone)
	if input.Timezone == "" {
		input.Timezone = "Asia/Colombo"
	}
	if _, err := time.LoadLocation(input.Timezone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter a valid timezone"})
		return
	}
	if err := writeAppSetting("preferences", input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}
