package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/geekytaurus115/assessment_techiebutler/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var store *handler.EmployeeStore
var router *gin.Engine

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	handler.InitDB()
	store = handler.NewEmployeeStore()
	employeeController := handler.NewEmployeeController(store)
	router = handler.SetupRoutes(employeeController)
	os.Exit(m.Run())
}

func TestCreateEmployee(t *testing.T) {
	emp := handler.Employee{
		Name:     "John Doe",
		Position: "Software Engineer",
		Salary:   60000.0,
	}
	jsonValue, _ := json.Marshal(emp)
	req, _ := http.NewRequest("POST", "/employee", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusCreated, response.Code)
	log.Println("Response body:", response.Body.String())

	var createdEmp handler.Employee
	err := json.Unmarshal(response.Body.Bytes(), &createdEmp)
	if err != nil {
		log.Println("Error unmarshalling response:", err)
	}
	log.Println("createdEmp--> ", createdEmp)

	assert.NotZero(t, createdEmp.ID)
	assert.Equal(t, emp.Name, createdEmp.Name)
	assert.Equal(t, emp.Position, createdEmp.Position)
	assert.Equal(t, emp.Salary, createdEmp.Salary)
}

func TestGetEmployeeByID(t *testing.T) {
	emp := handler.Employee{
		Name:     "Jane Doe",
		Position: "Product Manager",
		Salary:   80000.0,
	}
	store.CreateEmployee(&emp)

	req, _ := http.NewRequest("GET", "/employees/"+strconv.Itoa(emp.ID), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var fetchedEmp handler.Employee
	json.Unmarshal(response.Body.Bytes(), &fetchedEmp)
	assert.Equal(t, emp.ID, fetchedEmp.ID)
	assert.Equal(t, emp.Name, fetchedEmp.Name)
	assert.Equal(t, emp.Position, fetchedEmp.Position)
	assert.Equal(t, emp.Salary, fetchedEmp.Salary)
}

func TestUpdateEmployee(t *testing.T) {
	emp := handler.Employee{
		Name:     "Alice Smith",
		Position: "Designer",
		Salary:   50000.0,
	}
	store.CreateEmployee(&emp)

	updatedEmp := emp
	updatedEmp.Name = "Alice Johnson"
	updatedEmp.Salary = 55000.0
	jsonValue, _ := json.Marshal(updatedEmp)
	req, _ := http.NewRequest("PUT", "/employees/"+strconv.Itoa(emp.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var responseEmp handler.Employee
	json.Unmarshal(response.Body.Bytes(), &responseEmp)
	assert.Equal(t, updatedEmp, responseEmp)
}

func TestDeleteEmployee(t *testing.T) {
	emp := handler.Employee{
		Name:     "Bob Brown",
		Position: "HR Manager",
		Salary:   70000.0,
	}
	store.CreateEmployee(&emp)

	req, _ := http.NewRequest("DELETE", "/employees/"+strconv.Itoa(emp.ID), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusNoContent, response.Code)

	deletedEmp, err := store.GetEmployeeByID(emp.ID)
	assert.Nil(t, deletedEmp)
	assert.NotNil(t, err)
}

func TestListEmployees(t *testing.T) {
	store.CreateEmployee(&handler.Employee{
		Name:     "Charlie Davis",
		Position: "Marketing Specialist",
		Salary:   65000.0,
	})

	req, _ := http.NewRequest("GET", "/employees?page=1&limit=2", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var employees []handler.Employee
	json.Unmarshal(response.Body.Bytes(), &employees)
	assert.NotEmpty(t, employees)
	assert.LessOrEqual(t, len(employees), 2)
}
