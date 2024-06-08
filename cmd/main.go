package main

import (
	"log"
	"os"

	"github.com/geekytaurus115/assessment_techiebutler/handler"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	handler.InitDB()

	store := handler.NewEmployeeStore()
	employeeController := handler.NewEmployeeController(store)
	router := handler.SetupRoutes(employeeController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server is running on port %s", port)
	log.Fatal(router.Run(":" + port))
}
