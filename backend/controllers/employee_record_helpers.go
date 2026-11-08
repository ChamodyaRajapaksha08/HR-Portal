package controllers

import (
	"errors"
	"net/http"
	"strings"

	"hrm-backend/config"
	"hrm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func resolveRecordEmployee(employeeID, employeeName string) (models.Employee, error) {
	employeeID = strings.TrimSpace(employeeID)
	employeeName = strings.TrimSpace(employeeName)
	var employee models.Employee
	query := config.DB
	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	} else if employeeName != "" {
		query = query.Where("LOWER(TRIM(name)) = LOWER(TRIM(?))", employeeName).Order("id")
	} else {
		return employee, errors.New("Employee ID is required")
	}
	err := query.First(&employee).Error
	return employee, err
}

func respondEmployeeResolutionError(c *gin.Context, err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return true
	}
	if err.Error() == "Employee ID is required" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return true
	}
	return false
}

func ownLeaveQuery(employee models.Employee) *gorm.DB {
	return ownLeaveQueryFrom(config.DB, employee)
}

func ownLeaveQueryFrom(database *gorm.DB, employee models.Employee) *gorm.DB {
	return database.Where(
		"employee_id = ? OR ((employee_id IS NULL OR employee_id = '') AND LOWER(TRIM(employee)) = LOWER(TRIM(?)))",
		employee.EmployeeID,
		employee.FullName,
	)
}
