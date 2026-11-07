package controllers

import (
	"hrm-backend/config"
	"hrm-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateDepartment creates a new department record
func CreateDepartment(c *gin.Context) {
	var input models.Department
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department name is required"})
		return
	}

	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": input.ID, "name": input.Name, "head": input.Head,
		"employees": input.Employees, "status": input.Status,
		"created_date": input.CreatedAt.Format("Jan 02, 2006"),
	})
}

func UpdateDepartment(c *gin.Context) {
	var department models.Department
	if err := config.DB.First(&department, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}

	var input models.Department
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Department name is required"})
		return
	}

	department.Name = input.Name
	department.Head = input.Head
	department.Employees = input.Employees
	department.Status = input.Status
	if err := config.DB.Save(&department).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": department.ID, "name": department.Name, "head": department.Head,
		"employees": department.Employees, "status": department.Status,
		"created_date": department.CreatedAt.Format("Jan 02, 2006"),
	})
}

func DeleteDepartment(c *gin.Context) {
	result := config.DB.Delete(&models.Department{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Department not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
