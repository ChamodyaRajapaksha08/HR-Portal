package routes

import (
	"hrm-backend/auth"
	"hrm-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/login", controllers.Login)

	protected := api.Group("")
	protected.Use(auth.AuthRequired())

	management := protected.Group("")
	management.Use(auth.RequireRole("Admin", "HR Manager"))
	{
		management.GET("/employees", controllers.GetEmployees)
		management.POST("/employees", controllers.CreateEmployee)
		management.PUT("/employees/:employee_id", controllers.UpdateEmployee)
		management.DELETE("/employees/:employee_id", controllers.DeleteEmployee)
		management.GET("/departments", controllers.GetDepartments)
		management.POST("/departments", controllers.CreateDepartment)
		management.PUT("/departments/:id", controllers.UpdateDepartment)
		management.DELETE("/departments/:id", controllers.DeleteDepartment)
		management.GET("/attendance", controllers.GetAttendance)
		management.POST("/attendance", controllers.CreateAttendance)
		management.PUT("/attendance/:id", controllers.UpdateAttendance)
		management.DELETE("/attendance/:id", controllers.DeleteAttendance)
		management.GET("/leaves", controllers.GetLeaves)
		management.POST("/leaves", controllers.CreateLeave)
		management.PUT("/leaves/:id", controllers.UpdateLeave)
		management.DELETE("/leaves/:id", controllers.DeleteLeave)
		management.GET("/payroll", controllers.GetPayroll)
		management.POST("/payroll", controllers.CreatePayroll)
		management.PUT("/payroll/:id/pay", controllers.MarkPayrollPaid)
		management.PUT("/payroll/:id", controllers.UpdatePayroll)
		management.DELETE("/payroll/:id", controllers.DeletePayroll)
		management.GET("/reports", controllers.GetReports)
		management.GET("/dashboard-summary", controllers.GetDashboardSummary)
		management.GET("/settings", controllers.GetSettings)
		management.GET("/roles", controllers.GetRoles)
		management.GET("/users-roles", controllers.GetUsersRoles)
	}

	admin := protected.Group("")
	admin.Use(auth.RequireRole("Admin"))
	{
		admin.PUT("/settings/company", controllers.UpdateCompanySettings)
		admin.PUT("/settings/attendance", controllers.UpdateAttendanceSettings)
		admin.PUT("/settings/preferences", controllers.UpdatePreferenceSettings)
		admin.PUT("/users-roles/:id", controllers.UpdateUserRole)
	}

	protected.GET("/profile", controllers.GetProfile)
	protected.PUT("/profile", controllers.UpdateProfile)
	protected.POST("/profile/change-password", controllers.ChangeProfilePassword)

	protected.GET("/me/profile", controllers.GetMyProfile)
	protected.PUT("/me/profile", controllers.UpdateMyProfile)
	protected.GET("/me/attendance", controllers.GetMyAttendance)
	protected.GET("/me/leaves", controllers.GetMyLeaves)
	protected.POST("/me/leaves", controllers.CreateMyLeave)
	protected.GET("/me/leave-balance", controllers.GetMyLeaveBalance)
	protected.GET("/me/payslips", controllers.GetMyPayslips)
}
