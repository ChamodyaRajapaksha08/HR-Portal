package controllers

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"hrm-backend/config"
	"hrm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type payrollInput struct {
	EmployeeID  string  `json:"employee_id"`
	BasicSalary float64 `json:"basic_salary"`
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
}

type payrollResponse struct {
	ID          uint    `json:"id"`
	EmployeeID  string  `json:"employee_id"`
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	Department  string  `json:"department"`
	BasicSalary float64 `json:"basic_salary"`
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
	NetSalary   float64 `json:"net_salary"`
	Status      string  `json:"status"`
	PaidDate    *string `json:"paid_date"`
}

func validatePayrollAmounts(c *gin.Context, input payrollInput) bool {
	if math.IsNaN(input.BasicSalary) || math.IsInf(input.BasicSalary, 0) || input.BasicSalary <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Basic salary must be greater than 0"})
		return false
	}
	if math.IsNaN(input.Allowances) || math.IsInf(input.Allowances, 0) || input.Allowances < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Allowances must be 0 or greater"})
		return false
	}
	if math.IsNaN(input.Deductions) || math.IsInf(input.Deductions, 0) || input.Deductions < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deductions must be 0 or greater"})
		return false
	}
	if input.Deductions > input.BasicSalary+input.Allowances {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deductions cannot exceed basic salary plus allowances"})
		return false
	}
	if math.IsInf(input.BasicSalary+input.Allowances-input.Deductions, 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payroll amounts are too large"})
		return false
	}
	return true
}

func bindPayrollInput(c *gin.Context) (payrollInput, bool) {
	var input payrollInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll request body"})
		return input, false
	}
	return input, true
}

func payrollID(c *gin.Context) (uint, bool) {
	parsed, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || parsed == 0 || parsed > uint64(^uint(0)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payroll ID"})
		return 0, false
	}
	return uint(parsed), true
}

func loadPayrollResponse(id uint) (payrollResponse, error) {
	var result payrollResponse
	err := config.DB.Raw(`SELECT p.id, e.employee_id, e.name AS full_name, e.email, e.department,
		p.basic_salary, p.allowances, p.deductions, p.net_salary,
		TO_CHAR(p.paid_date, 'YYYY-MM-DD') AS paid_date,
		CASE WHEN p.paid_date IS NULL THEN 'PENDING' ELSE 'PAID' END AS status
		FROM payroll p JOIN employees e ON e.id = p.employee_id
		WHERE p.id = ? AND e.deleted_at IS NULL`, id).Scan(&result).Error
	if err != nil {
		return result, err
	}
	if result.ID == 0 {
		return result, gorm.ErrRecordNotFound
	}
	return result, nil
}

func CreatePayroll(c *gin.Context) {
	input, ok := bindPayrollInput(c)
	if !ok {
		return
	}
	input.EmployeeID = strings.TrimSpace(input.EmployeeID)
	if input.EmployeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
		return
	}
	if !validatePayrollAmounts(c, input) {
		return
	}

	var employee models.Employee
	if err := config.DB.Where("employee_id = ?", input.EmployeeID).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payroll := models.Payroll{
		EmployeeID:  employee.ID,
		BasicSalary: input.BasicSalary,
		Allowances:  input.Allowances,
		Deductions:  input.Deductions,
		NetSalary:   input.BasicSalary + input.Allowances - input.Deductions,
		PaidDate:    nil,
	}
	if err := config.DB.Create(&payroll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := loadPayrollResponse(payroll.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

func UpdatePayroll(c *gin.Context) {
	id, ok := payrollID(c)
	if !ok {
		return
	}
	input, ok := bindPayrollInput(c)
	if !ok {
		return
	}
	if !validatePayrollAmounts(c, input) {
		return
	}

	var payroll models.Payroll
	if err := config.DB.First(&payroll, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Payroll record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	payroll.BasicSalary = input.BasicSalary
	payroll.Allowances = input.Allowances
	payroll.Deductions = input.Deductions
	payroll.NetSalary = input.BasicSalary + input.Allowances - input.Deductions
	if err := config.DB.Save(&payroll).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := loadPayrollResponse(payroll.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func MarkPayrollPaid(c *gin.Context) {
	id, ok := payrollID(c)
	if !ok {
		return
	}
	paidDate, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to load Asia/Colombo timezone: %v", err)})
		return
	}
	result := config.DB.Model(&models.Payroll{}).
		Where("id = ? AND paid_date IS NULL", id).
		Update("paid_date", time.Now().In(paidDate))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		var payroll models.Payroll
		if err := config.DB.First(&payroll, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Payroll record not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusConflict, gin.H{"error": "Already paid"})
		return
	}
	c.Status(http.StatusOK)
}

func DeletePayroll(c *gin.Context) {
	id, ok := payrollID(c)
	if !ok {
		return
	}
	result := config.DB.Delete(&models.Payroll{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payroll record not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
