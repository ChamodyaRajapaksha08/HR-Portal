package controllers

import (
	"errors"
	"hrm-backend/auth"
	"hrm-backend/config"
	"hrm-backend/models"
	"hrm-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type profileResponse struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	EmployeeID string `json:"employee_id"`
	Position   string `json:"position"`
	Department string `json:"department"`
	Role       string `json:"role"`
}

var errLastAdmin = errors.New("cannot remove the last Admin")

func profileFor(employee models.Employee) profileResponse {
	return profileResponse{
		Name:       employee.FullName,
		Email:      employee.Email,
		EmployeeID: employee.EmployeeID,
		Position:   employee.Position,
		Department: employee.Department,
		Role:       employee.Role,
	}
}

func loadCurrentEmployee(c *gin.Context) (models.Employee, error) {
	claims, ok := auth.CurrentClaims(c)
	if !ok {
		return models.Employee{}, gorm.ErrRecordNotFound
	}
	var employee models.Employee
	err := config.DB.First(&employee, claims.EmployeeDBID).Error
	return employee, err
}

func GetProfile(c *gin.Context) {
	employee, err := loadCurrentEmployee(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profileFor(employee))
}

func UpdateProfile(c *gin.Context) {
	employee, err := loadCurrentEmployee(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Position string `json:"position"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and a valid email are required"})
		return
	}
	employee.FullName = strings.TrimSpace(input.Name)
	employee.Email = strings.TrimSpace(input.Email)
	employee.Position = strings.TrimSpace(input.Position)
	if employee.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	if err := config.DB.Save(&employee).Error; err != nil {
		var sqlStateError interface {
			error
			SQLState() string
		}
		if errors.As(err, &sqlStateError) && sqlStateError.SQLState() == "23505" {
			c.JSON(http.StatusConflict, gin.H{"error": "That email is already in use"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profileFor(employee))
}

func ChangeProfilePassword(c *gin.Context) {
	employee, err := loadCurrentEmployee(c)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current, new, and confirmation passwords are required"})
		return
	}
	if !utils.CheckPassword(employee.Password, input.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}
	if utf8.RuneCountInString(input.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password must be at least 8 characters"})
		return
	}
	if input.NewPassword != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirmation do not match"})
		return
	}
	hashed, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to secure password"})
		return
	}
	if err := config.DB.Model(&employee).Update("password", hashed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

type roleDefinition struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func GetRoles(c *gin.Context) {
	c.JSON(http.StatusOK, []roleDefinition{
		{Role: "Admin", Permissions: []string{"Everything"}},
		{Role: "HR Manager", Permissions: []string{"Employees", "Departments", "Attendance", "Leaves", "Payroll", "Reports"}},
		{Role: "Employee", Permissions: []string{"Own profile", "Own attendance", "Apply for leave"}},
	})
}

type userRole struct {
	ID         uint   `json:"id"`
	EmployeeID string `json:"employee_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
}

func GetUsersRoles(c *gin.Context) {
	users := []userRole{}
	err := config.DB.Model(&models.Employee{}).
		Select("id, employee_id, name, email, role").
		Order("employee_id").
		Scan(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func UpdateUserRole(c *gin.Context) {
	var input struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}
	input.Role = strings.TrimSpace(input.Role)
	if input.Role != "Admin" && input.Role != "HR Manager" && input.Role != "Employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be Admin, HR Manager, or Employee"})
		return
	}

	parsedID, parseErr := strconv.ParseUint(c.Param("id"), 10, 64)
	if parseErr != nil || parsedID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}
	id := uint(parsedID)
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var employee models.Employee
		if err := tx.First(&employee, id).Error; err != nil {
			return err
		}
		if employee.Role == "Admin" && input.Role != "Admin" {
			admins := []models.Employee{}
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("role = ?", "Admin").
				Find(&admins).Error; err != nil {
				return err
			}
			if len(admins) <= 1 {
				return errLastAdmin
			}
		}
		return tx.Model(&employee).Update("role", input.Role).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	if errors.Is(err, errLastAdmin) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot remove the last Admin"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User role updated"})
}
