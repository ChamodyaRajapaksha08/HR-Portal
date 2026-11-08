package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"
	_ "time/tzdata"

	"hrm-backend/auth"
	"hrm-backend/config"
	"hrm-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	annualLeaveEntitlement = 14
	casualLeaveEntitlement = 7
	sickLeaveEntitlement   = 7
)

var leaveEntitlements = []struct {
	Name string
	Days int
}{
	{Name: "Annual", Days: annualLeaveEntitlement},
	{Name: "Casual", Days: casualLeaveEntitlement},
	{Name: "Sick", Days: sickLeaveEntitlement},
}

var (
	errLeaveDateOverlap         = errors.New("leave dates overlap existing requests")
	errInvalidExistingLeaveDate = errors.New("unable to validate existing leave dates")
)

func currentEmployee(c *gin.Context) (models.Employee, error) {
	claims, ok := auth.CurrentClaims(c)
	if !ok {
		return models.Employee{}, gorm.ErrRecordNotFound
	}
	var employee models.Employee
	err := config.DB.First(&employee, claims.EmployeeDBID).Error
	return employee, err
}

func respondCurrentEmployeeError(c *gin.Context, err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return true
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return true
}

func GetMyProfile(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	c.JSON(http.StatusOK, profileFor(employee))
}

func UpdateMyProfile(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	var input struct {
		Name     string `json:"name" binding:"required"`
		Position string `json:"position"`
	}
	if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}
	employee.FullName = strings.TrimSpace(input.Name)
	employee.Position = strings.TrimSpace(input.Position)
	if err := config.DB.Save(&employee).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profileFor(employee))
}

func GetMyAttendance(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	location, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load attendance timezone"})
		return
	}
	now := time.Now().In(location)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, location)
	nextMonth := monthStart.AddDate(0, 1, 0)

	var attendance []models.Attendance
	if err := config.DB.Where("employee_id = ?", employee.EmployeeID).
		Order("created_at DESC").Find(&attendance).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var monthly []models.Attendance
	if err := config.DB.Where("employee_id = ? AND created_at >= ? AND created_at < ?",
		employee.EmployeeID, monthStart, nextMonth).Find(&monthly).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	present, late, absent := 0, 0, 0
	for _, record := range monthly {
		switch strings.ToLower(strings.TrimSpace(record.Status)) {
		case "present":
			present++
		case "late":
			late++
		case "absent":
			absent++
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"summary": gin.H{
			"year": now.Year(), "month": int(now.Month()),
			"present": present, "late": late, "absent": absent,
		},
		"attendance": attendance,
	})
}

func GetMyLeaves(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	var leaves []models.LeaveRequest
	if err := ownLeaveQuery(employee).Order("id DESC").Find(&leaves).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	requests := make([]leaveResponse, 0, len(leaves))
	for _, leave := range leaves {
		requests = append(requests, leaveResponse{
			ID: leave.ID, EmployeeID: leave.EmployeeID, Employee: leave.Employee,
			LeaveType: leave.LeaveType, FromDate: leave.FromDate, ToDate: leave.ToDate,
			Days: leave.Days, Reason: leave.Reason, Status: leave.Status, AppliedOn: leave.AppliedOn,
		})
	}
	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

func CreateMyLeave(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	var input struct {
		LeaveType string `json:"leave_type" binding:"required"`
		FromDate  string `json:"from_date" binding:"required"`
		ToDate    string `json:"to_date" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave type, from date, to date, and reason are required"})
		return
	}
	input.LeaveType = strings.TrimSpace(input.LeaveType)
	input.Reason = strings.TrimSpace(input.Reason)
	if input.LeaveType == "" || input.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Leave type, from date, to date, and reason are required"})
		return
	}
	fromDate, err := time.Parse("2006-01-02", input.FromDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "From date must use YYYY-MM-DD format"})
		return
	}
	toDate, err := time.Parse("2006-01-02", input.ToDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "To date must use YYYY-MM-DD format"})
		return
	}
	if toDate.Before(fromDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "To date cannot be before from date"})
		return
	}
	fromDate, toDate = dateOnly(fromDate), dateOnly(toDate)

	location, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load leave timezone"})
		return
	}
	req := models.LeaveRequest{
		EmployeeID: employee.EmployeeID,
		Employee:   employee.FullName,
		LeaveType:  strings.TrimSpace(input.LeaveType),
		FromDate:   fromDate.Format("2006-01-02"),
		ToDate:     toDate.Format("2006-01-02"),
		Days:       int(toDate.Sub(fromDate).Hours()/24) + 1,
		Reason:     strings.TrimSpace(input.Reason),
		Status:     "Pending",
		AppliedOn:  time.Now().In(location).Format("2006-01-02"),
	}
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var lockedEmployee models.Employee
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedEmployee, employee.ID).Error; err != nil {
			return err
		}
		var existing []models.LeaveRequest
		if err := ownLeaveQueryFrom(tx, lockedEmployee).
			Where("status IN ?", []string{"Pending", "Approved"}).
			Find(&existing).Error; err != nil {
			return err
		}
		for _, leave := range existing {
			existingFrom, parseErr := parseLeaveDate(leave.FromDate)
			if parseErr != nil {
				return errInvalidExistingLeaveDate
			}
			existingTo, parseErr := parseLeaveDate(leave.ToDate)
			if parseErr != nil {
				return errInvalidExistingLeaveDate
			}
			if !toDate.Before(dateOnly(existingFrom)) && !fromDate.After(dateOnly(existingTo)) {
				return errLeaveDateOverlap
			}
		}
		return tx.Create(&req).Error
	})
	if errors.Is(err, errLeaveDateOverlap) {
		c.JSON(http.StatusConflict, gin.H{"error": "Leave dates overlap an existing Pending or Approved request"})
		return
	}
	if err != nil {
		if errors.Is(err, errInvalidExistingLeaveDate) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to validate existing leave dates"})
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, leaveResponse{
		ID: req.ID, EmployeeID: req.EmployeeID, Employee: req.Employee,
		LeaveType: req.LeaveType, FromDate: req.FromDate, ToDate: req.ToDate,
		Days: req.Days, Reason: req.Reason, Status: req.Status, AppliedOn: req.AppliedOn,
	})
}

func GetMyLeaveBalance(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	location, err := time.LoadLocation("Asia/Colombo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to load leave timezone"})
		return
	}
	now := time.Now().In(location)
	yearStart := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	nextYear := yearStart.AddDate(1, 0, 0)
	var approved []models.LeaveRequest
	if err := ownLeaveQuery(employee).Where("status = ?", "Approved").Find(&approved).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	usedByType := make(map[string]int)
	for _, leave := range approved {
		fromDate, parseErr := parseLeaveDate(leave.FromDate)
		if parseErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to calculate leave balance"})
			return
		}
		toDate, parseErr := parseLeaveDate(leave.ToDate)
		if parseErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to calculate leave balance"})
			return
		}
		fromDate, toDate = dateOnly(fromDate), dateOnly(toDate)
		if toDate.Before(yearStart) || !fromDate.Before(nextYear) {
			continue
		}
		if fromDate.Before(yearStart) {
			fromDate = yearStart
		}
		lastDayOfYear := nextYear.AddDate(0, 0, -1)
		if toDate.After(lastDayOfYear) {
			toDate = lastDayOfYear
		}
		days := int(toDate.Sub(fromDate).Hours()/24) + 1
		leaveType := leaveCategory(leave.LeaveType)
		if leaveType != "" {
			usedByType[leaveType] += days
		}
	}
	balances := make([]gin.H, 0, len(leaveEntitlements))
	for _, entitlement := range leaveEntitlements {
		used := usedByType[entitlement.Name]
		remaining := entitlement.Days - used
		if remaining < 0 {
			remaining = 0
		}
		balances = append(balances, gin.H{
			"leave_type":  entitlement.Name,
			"entitlement": entitlement.Days,
			"used":        used,
			"remaining":   remaining,
		})
	}
	c.JSON(http.StatusOK, gin.H{"year": now.Year(), "balances": balances})
}

func GetMyPayslips(c *gin.Context) {
	employee, err := currentEmployee(c)
	if err != nil {
		respondCurrentEmployeeError(c, err)
		return
	}
	type payslip struct {
		ID          uint    `json:"id"`
		PaidDate    *string `json:"paid_date"`
		BasicSalary float64 `json:"basic_salary"`
		Allowances  float64 `json:"allowances"`
		Deductions  float64 `json:"deductions"`
		NetSalary   float64 `json:"net_salary"`
		Status      string  `json:"status"`
	}
	rows := []payslip{}
	err = config.DB.Raw(`SELECT p.id,
		TO_CHAR(p.paid_date, 'YYYY-MM-DD') AS paid_date,
		p.basic_salary, p.allowances, p.deductions, p.net_salary,
		CASE WHEN p.paid_date IS NULL THEN 'PENDING' ELSE 'PAID' END AS status
		FROM payroll p
		WHERE p.employee_id = ?
		ORDER BY p.id DESC`, employee.ID).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"payslips": rows})
}

func dateOnly(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

func parseLeaveDate(value string) (time.Time, error) {
	for _, layout := range []string{"2006-01-02", "Jan 02, 2006", "02/01/2006"} {
		if date, err := time.Parse(layout, strings.TrimSpace(value)); err == nil {
			return date, nil
		}
	}
	return time.Time{}, errors.New("invalid leave date")
}

func leaveCategory(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, entitlement := range leaveEntitlements {
		if strings.Contains(value, strings.ToLower(entitlement.Name)) {
			return entitlement.Name
		}
	}
	return ""
}
