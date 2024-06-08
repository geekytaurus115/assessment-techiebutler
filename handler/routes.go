package handler

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(ec *EmployeeController) *gin.Engine {
	router := gin.Default()
	router.POST("/employee", ec.CreateEmployee)
	router.GET("/employees", ec.ListEmployees)
	router.GET("/employees/:id", ec.GetEmployeeByID)
	router.PUT("/employees/:id", ec.UpdateEmployee)
	router.DELETE("/employees/:id", ec.DeleteEmployee)
	return router
}
