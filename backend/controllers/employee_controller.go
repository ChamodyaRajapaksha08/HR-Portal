package controllers

import (
	"errors"
	"hrm-backend/auth"
	"hrm-backend/config"
	"hrm-backend/models"
	"hrm-backend/utils"
	"net/http" //
	//"time"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request"})
		return
	}
	if loginInput.Email == "" || loginInput.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}
	now := time.Now()
	attemptKey := loginAttemptKey(loginInput.Email, c.ClientIP())
	if !reserveLoginAttempt(attemptKey, now) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many failed login attempts. Try again in 15 minutes."})
		return
	}

	var employee models.Employee
	if err := config.DB.Where("email = ?", loginInput.Email).First(&employee).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			finishLoginAttempt(attemptKey, now, false, false)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to authenticate"})
			return
		}
		finishLoginAttempt(attemptKey, now, true, false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPassword(employee.Password, loginInput.Password) {
		finishLoginAttempt(attemptKey, now, true, false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	finishLoginAttempt(attemptKey, now, false, true)
	token, err := auth.GenerateToken(employee.ID, employee.EmployeeID, employee.Email, employee.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create login session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":          employee.ID,
			"employee_id": employee.EmployeeID,
			"name":        employee.FullName,
			"email":       employee.Email,
			"position":    employee.Position,
			"department":  employee.Department,
			"role":        employee.Role,
		},
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
	var input struct {
		EmployeeID string `json:"employee_id" binding:"required"`
		FullName   string `json:"full_name" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required"`
		Department string `json:"department"`
		Position   string `json:"position"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to secure password"})
		return
	}

	employee := models.Employee{
		EmployeeID: input.EmployeeID,
		FullName:   input.FullName,
		Email:      input.Email,
		Password:   hashed,
		Department: input.Department,
		Position:   input.Position,
		Status:     input.Status,
		Role:       "Employee",
	}
	if err := config.DB.Create(&employee).Error; err != nil {
		respondToEmployeeWriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, employee)
}

func UpdateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := config.DB.Where("employee_id = ?", c.Param("employee_id")).First(&employee).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	var input struct {
		FullName   string `json:"full_name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Department string `json:"department"`
		Position   string `json:"position"`
		Status     string `json:"status"`
	}
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
		hashed, err := utils.HashPassword(input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to secure password"})
			return
		}
		employee.Password = hashed
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

	var input struct {
		EmployeeID string `json:"employee_id"`
		Employee   string `json:"employee"`
		Department string `json:"department"`
		CheckIn    string `json:"check_in"`
		CheckOut   string `json:"check_out"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.EmployeeID != "" || input.Employee != "" {
		employee, err := resolveRecordEmployee(input.EmployeeID, input.Employee)
		if err != nil {
			if respondEmployeeResolutionError(c, err) {
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		record.EmployeeID = employee.EmployeeID
		record.Employee = employee.FullName
		record.Department = employee.Department
	}
	if input.Department != "" {
		record.Department = input.Department
	}
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
	var input struct {
		EmployeeID string `json:"employee_id"`
		Employee   string `json:"employee"`
		Department string `json:"department"`
		CheckIn    string `json:"check_in"`
		CheckOut   string `json:"check_out"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	employee, err := resolveRecordEmployee(input.EmployeeID, input.Employee)
	if err != nil {
		if respondEmployeeResolutionError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	record := models.Attendance{
		EmployeeID: employee.EmployeeID,
		Employee:   employee.FullName,
		Department: employee.Department,
		CheckIn:    input.CheckIn,
		CheckOut:   input.CheckOut,
		Status:     input.Status,
	}
	if input.Department != "" {
		record.Department = input.Department
	}

	if err := config.DB.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func GetPayroll(c *gin.Context) {
	type row struct {
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
	payrolls := []row{}
	err := config.DB.Raw(`SELECT p.id, e.employee_id, e.name AS full_name, e.email, e.department,
		p.basic_salary, p.allowances, p.deductions, p.net_salary,
		TO_CHAR(p.paid_date, 'YYYY-MM-DD') AS paid_date,
		CASE WHEN p.paid_date IS NULL THEN 'PENDING' ELSE 'PAID' END AS status
		FROM payroll p JOIN employees e ON e.id = p.employee_id
		WHERE e.deleted_at IS NULL ORDER BY e.employee_id, p.id`).Scan(&payrolls).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var total float64
	paid, pending := 0, 0
	for _, p := range payrolls {
		total += p.NetSalary
		if p.Status == "PAID" {
			paid++
		} else {
			pending++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"summary":  gin.H{"total_payroll": total, "processed": paid, "pending_approvals": pending},
		"payrolls": payrolls,
	})
}

func GetReports(c *gin.Context) {
	var total, present int64
	config.DB.Model(&models.Attendance{}).Count(&total)
	config.DB.Model(&models.Attendance{}).Where("status IN ?", []string{"Present", "Late"}).Count(&present)
	rate := 0.0
	if total > 0 {
		rate = float64(present) * 100 / float64(total)
	}
	type perf struct {
		Division string  `json:"division"`
		Value    float64 `json:"value"`
	}
	performance := []perf{}
	config.DB.Raw(`SELECT department AS division,
		ROUND(100.0 * SUM(CASE WHEN status IN ('Present','Late') THEN 1 ELSE 0 END) / COUNT(*), 1) AS value
		FROM attendances WHERE deleted_at IS NULL AND department <> '' GROUP BY department`).Scan(&performance)
	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{"attendance_rate": rate, "performance": performance},
		"reports": []gin.H{},
	})
}
