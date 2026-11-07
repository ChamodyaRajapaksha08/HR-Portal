package routes

import (
	"hrm-backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.GET("/employees", controllers.GetEmployees)
		api.PUT("/employees/:employee_id", controllers.UpdateEmployee)
		api.DELETE("/employees/:employee_id", controllers.DeleteEmployee)
		api.GET("/departments", controllers.GetDepartments)
		api.POST("/departments", controllers.CreateDepartment)
		api.PUT("/departments/:id", controllers.UpdateDepartment)
		api.DELETE("/departments/:id", controllers.DeleteDepartment)
		api.GET("/attendance", controllers.GetAttendance)
		api.PUT("/attendance/:id", controllers.UpdateAttendance)
		api.DELETE("/attendance/:id", controllers.DeleteAttendance)
		api.GET("/leaves", controllers.GetLeaves)
		api.PUT("/leaves/:id", controllers.UpdateLeave)
		api.DELETE("/leaves/:id", controllers.DeleteLeave)
		api.GET("/payroll", controllers.GetPayroll)
		api.GET("/reports", controllers.GetReports)
		api.GET("/dashboard-summary", controllers.GetDashboardSummary)
		api.POST("/employees", controllers.CreateEmployee)
		api.POST("/attendance", controllers.CreateAttendance)
		api.POST("/leaves", controllers.CreateLeave)
	}
}
