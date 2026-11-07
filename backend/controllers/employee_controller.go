package controllers

import (
	"errors"
	"hrm-backend/config"
	"hrm-backend/models"
	"net/http" //
	//"time"

	"github.com/gin-gonic/gin"
)

func respondToEmployeeWriteError(c *gin.Context, err error) {
	var sqlStateError interface {
		error
		SQLState() string
	}
	if errors.As(err, &sqlStateError) && sqlStateError.SQLState() == "23505" {
		c.JSON(http.StatusConflict, gin.H{
			"error": "An employee with this ID or email already exists. Please review the values and try again.",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

// Login function - handles employee authentication
func Login(c *gin.Context) {
	var loginInput struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password format"})
		return
	}

	var employee models.Employee
	if err := config.DB.Where("email = ?", loginInput.Email).First(&employee).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Simple password check (in production, use bcrypt or similar)
	if employee.Password != loginInput.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Return success with employee info
	c.JSON(http.StatusOK, gin.H{
		"message":     "Login successful",
		"employee_id": employee.EmployeeID,
		"full_name":   employee.FullName,
		"email":       employee.Email,
		"position":    employee.Position,
		"token":       "dummy_token_" + employee.Email, // In production, generate a proper JWT token
	})
}

// 1. සේවකයෝ ලැයිස්තුවම UI එකට යවන ක්‍රියාවලිය
func GetEmployees(c *gin.Context) {
	var employees []models.Employee
	if err := config.DB.Find(&employees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i := range employees {
		employees[i].Password = ""
	}
	c.JSON(http.StatusOK, employees)
}

// 2. අලුත් සේවකයෙක් ඩේටාබේස් එකට සේව් කරන ක්‍රියාවලිය
func CreateEmployee(c *gin.Context) {
	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		respondToEmployeeWriteError(c, err)
		return
	}
	input.Password = ""
	c.JSON(http.StatusOK, input)
}

func UpdateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := config.DB.Where("employee_id = ?", c.Param("employee_id")).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var input models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee.FullName = input.FullName
	employee.Email = input.Email
	employee.Department = input.Department
	employee.Position = input.Position
	employee.Status = input.Status
	if input.Password != "" {
		employee.Password = input.Password
	}
	if err := config.DB.Save(&employee).Error; err != nil {
		respondToEmployeeWriteError(c, err)
		return
	}
	employee.Password = ""
	c.JSON(http.StatusOK, employee)
}

func DeleteEmployee(c *gin.Context) {
	result := config.DB.Where("employee_id = ?", c.Param("employee_id")).Delete(&models.Employee{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func GetDashboardSummary(c *gin.Context) {
	var totalEmployees int64
	var departmentsCount int64
	var presentToday int64
	var pendingLeaves int64

	if err := config.DB.Model(&models.Employee{}).Count(&totalEmployees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.Department{}).Count(&departmentsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.Attendance{}).Where("status = ?", "Present").Count(&presentToday).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&models.LeaveRequest{}).Where("status = ?", "Pending").Count(&pendingLeaves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_employees":   totalEmployees,
		"total_departments": departmentsCount,
		"present_today":     presentToday,
		"pending_leaves":    pendingLeaves,
	})
}

func GetDepartments(c *gin.Context) {
	var dbDepts []models.Department
	if err := config.DB.Find(&dbDepts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	departments := make([]gin.H, 0, len(dbDepts))
	activeDepartments := 0
	departmentHeads := 0
	for _, department := range dbDepts {
		departments = append(departments, gin.H{
			"id": department.ID, "name": department.Name, "head": department.Head,
			"employees": department.Employees, "status": department.Status,
			"created_date": department.CreatedAt.Format("Jan 02, 2006"),
		})
		if department.Status == "Active" {
			activeDepartments++
		}
		if department.Head != "" {
			departmentHeads++
		}
	}
	var totalEmployees int64
	if err := config.DB.Model(&models.Employee{}).Count(&totalEmployees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"total_departments": len(dbDepts), "active_departments": activeDepartments,
			"total_employees": totalEmployees, "department_heads": departmentHeads,
		},
		"departments": departments,
	})
}

func GetAttendance(c *gin.Context) {
	var attendance []models.Attendance
	if err := config.DB.Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalEmployees int64
	if err := config.DB.Model(&models.Employee{}).Count(&totalEmployees).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var presentToday int64
	var absentToday int64
	var lateComers int64

	for _, record := range attendance {
		switch record.Status {
		case "Present":
			presentToday++
		case "Absent":
			absentToday++
		case "Late":
			lateComers++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"summary":    gin.H{"present_today": presentToday, "absent_today": absentToday, "late_comers": lateComers, "total_employees": totalEmployees},
		"attendance": attendance,
	})
}

func UpdateAttendance(c *gin.Context) {
	var record models.Attendance
	if err := config.DB.First(&record, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}

	var input models.Attendance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	record.Employee = input.Employee
	record.Department = input.Department
	record.CheckIn = input.CheckIn
	record.CheckOut = input.CheckOut
	record.Status = input.Status
	if err := config.DB.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, record)
}

func DeleteAttendance(c *gin.Context) {
	result := config.DB.Delete(&models.Attendance{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance record not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func CreateAttendance(c *gin.Context) {
	var input models.Attendance
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, input)
}


func GetPayroll(c *gin.Context) {
	payrolls := []gin.H{
		{"employee_id": "EMP001", "full_name": "Hashini Nisansala", "email": "hashini@itum.com", "department": "IT Division", "basic_salary": "LKR 150,000", "allowances": "LKR 25,000", "deductions": "LKR 5,000", "net_salary": "LKR 170,000", "status": "PAID"},
		{"employee_id": "EMP002", "full_name": "John Doe", "email": "johndoe@company.com", "department": "HR Management", "basic_salary": "LKR 120,000", "allowances": "LKR 15,000", "deductions": "LKR 3,500", "net_salary": "LKR 131,500", "status": "PAID"},
		{"employee_id": "EMP003", "full_name": "Kasun Perera", "email": "kasun@company.com", "department": "Finance", "basic_salary": "LKR 135,000", "allowances": "LKR 20,000", "deductions": "LKR 4,000", "net_salary": "LKR 151,000", "status": "PENDING"},
	}

	c.JSON(http.StatusOK, gin.H{
		"summary":  gin.H{"total_payroll": "LKR 1,845,000", "processed": "112 Employees", "pending_approvals": "8 Employees"},
		"payrolls": payrolls,
	})
}

func GetReports(c *gin.Context) {
	reports := []gin.H{
		{"name": "May 2026 Payroll Summary", "type": "PDF", "generated_date": "May 17, 2026", "status": "Ready"},
		{"name": "Employee Attendance Log (Q1)", "type": "EXCEL", "generated_date": "April 30, 2026", "status": "Ready"},
	}

	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{"attendance_rate": "94.2%", "performance": []gin.H{{"division": "IT Division", "value": "96%"}, {"division": "Human Resources", "value": "88%"}, {"division": "Finance", "value": "92%"}}},
		"reports": reports,
	})
}
