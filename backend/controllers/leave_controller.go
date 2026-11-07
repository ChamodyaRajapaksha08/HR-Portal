package controllers

import (
	"hrm-backend/config"
	"hrm-backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type leaveResponse struct {
	ID        uint   `json:"id"`
	Employee  string `json:"employee"`
	LeaveType string `json:"leave_type"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	Days      int    `json:"days"`
	Reason    string `json:"reason"`
	Status    string `json:"status"` // Pending, Approved, Rejected
	AppliedOn string `json:"applied_on"`
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
			ID: leave.ID, Employee: leave.Employee, LeaveType: leave.LeaveType,
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
	var req models.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status == "" {
		req.Status = "Pending"
	}
	req.AppliedOn = time.Now().Format("Jan 02, 2006")
	if err := config.DB.Create(&req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leaveResponse{
		ID: req.ID, Employee: req.Employee, LeaveType: req.LeaveType,
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

	var input models.LeaveRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	leave.Employee = input.Employee
	leave.LeaveType = input.LeaveType
	leave.FromDate = input.FromDate
	leave.ToDate = input.ToDate
	leave.Days = input.Days
	leave.Reason = input.Reason
	leave.Status = input.Status
	if err := config.DB.Save(&leave).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaveResponse{
		ID: leave.ID, Employee: leave.Employee, LeaveType: leave.LeaveType,
		FromDate: leave.FromDate, ToDate: leave.ToDate, Days: leave.Days,
		Reason: leave.Reason, Status: leave.Status, AppliedOn: leave.AppliedOn,
	})
}

func DeleteLeave(c *gin.Context) {
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