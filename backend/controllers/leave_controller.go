package controllers

import (
	"errors"
	"hrm-backend/auth"
	"hrm-backend/config"
	"hrm-backend/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type leaveResponse struct {
	ID         uint   `json:"id"`
	EmployeeID string `json:"employee_id"`
	Employee   string `json:"employee"`
	LeaveType  string `json:"leave_type"`
	FromDate   string `json:"from_date"`
	ToDate     string `json:"to_date"`
	Days       int    `json:"days"`
	Reason     string `json:"reason"`
	Status     string `json:"status"` // Pending, Approved, Rejected
	AppliedOn  string `json:"applied_on"`
}

func GetLeaves(c *gin.Context) {
	var leaves []models.LeaveRequest
	if err := config.DB.Order("id DESC").Find(&leaves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	requests := make([]leaveResponse, 0, len(leaves))
	summary := gin.H{
		"total_leave_requests": len(leaves),
		"pending_requests":     0,
		"approved":             0,
		"rejected":             0,
	}
	for _, leave := range leaves {
		requests = append(requests, leaveResponse{
			ID: leave.ID, EmployeeID: leave.EmployeeID, Employee: leave.Employee, LeaveType: leave.LeaveType,
			FromDate: leave.FromDate, ToDate: leave.ToDate, Days: leave.Days,
			Reason: leave.Reason, Status: leave.Status, AppliedOn: leave.AppliedOn,
		})
		switch leave.Status {
		case "Pending":
			summary["pending_requests"] = summary["pending_requests"].(int) + 1
		case "Approved":
			summary["approved"] = summary["approved"].(int) + 1
		case "Rejected":
			summary["rejected"] = summary["rejected"].(int) + 1
		}
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary, "requests": requests})
}

func CreateLeave(c *gin.Context) {
	var input struct {
		EmployeeID string `json:"employee_id"`
		Employee   string `json:"employee"`
		LeaveType  string `json:"leave_type"`
		FromDate   string `json:"from_date"`
		ToDate     string `json:"to_date"`
		Days       int    `json:"days"`
		Reason     string `json:"reason"`
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
	req := models.LeaveRequest{
		EmployeeID: employee.EmployeeID,
		Employee:   employee.FullName,
		LeaveType:  input.LeaveType,
		FromDate:   input.FromDate,
		ToDate:     input.ToDate,
		Days:       input.Days,
		Reason:     input.Reason,
		Status:     "Pending",
		AppliedOn:  time.Now().Format("Jan 02, 2006"),
	}
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leaveResponse{
		ID: req.ID, EmployeeID: req.EmployeeID, Employee: req.Employee, LeaveType: req.LeaveType,
		FromDate: req.FromDate, ToDate: req.ToDate, Days: req.Days,
		Reason: req.Reason, Status: req.Status, AppliedOn: req.AppliedOn,
	})
}

func UpdateLeave(c *gin.Context) {
	var leave models.LeaveRequest
	if err := config.DB.First(&leave, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		return
	}
	if leave.Status != "Pending" {
		c.JSON(http.StatusConflict, gin.H{"error": "Only Pending leave requests can be changed"})
		return
	}

	var input struct {
		EmployeeID string `json:"employee_id"`
		LeaveType  string `json:"leave_type"`
		FromDate   string `json:"from_date"`
		ToDate     string `json:"to_date"`
		Days       int    `json:"days"`
		Reason     string `json:"reason"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.EmployeeID != "" && input.EmployeeID != leave.EmployeeID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "A leave request cannot be reassigned to another employee"})
		return
	}
	if input.Status != "" && input.Status != "Pending" && input.Status != "Approved" && input.Status != "Rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be Pending, Approved, or Rejected"})
		return
	}
	if input.Status == "Approved" || input.Status == "Rejected" {
		claims, ok := auth.CurrentClaims(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		var currentEmployee models.Employee
		if err := config.DB.Select("id", "name").First(&currentEmployee, claims.EmployeeDBID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		isOwnRequest := leave.EmployeeID != "" && leave.EmployeeID == claims.EmployeeID
		if leave.EmployeeID == "" {
			isOwnRequest = strings.EqualFold(strings.TrimSpace(leave.Employee), strings.TrimSpace(currentEmployee.FullName))
		}
		if isOwnRequest {
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot approve or reject your own leave request"})
			return
		}
	}
	leave.LeaveType = input.LeaveType
	leave.FromDate = input.FromDate
	leave.ToDate = input.ToDate
	leave.Days = input.Days
	leave.Reason = input.Reason
	if input.Status != "" {
		leave.Status = input.Status
	}
	if err := config.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaveResponse{
		ID: leave.ID, EmployeeID: leave.EmployeeID, Employee: leave.Employee, LeaveType: leave.LeaveType,
		FromDate: leave.FromDate, ToDate: leave.ToDate, Days: leave.Days,
		Reason: leave.Reason, Status: leave.Status, AppliedOn: leave.AppliedOn,
	})
}

func DeleteLeave(c *gin.Context) {
	var leave models.LeaveRequest
	if err := config.DB.First(&leave, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if leave.Status != "Pending" {
		c.JSON(http.StatusConflict, gin.H{"error": "Only Pending leave requests can be changed"})
		return
	}
	result := config.DB.Delete(&models.LeaveRequest{}, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Leave request not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
