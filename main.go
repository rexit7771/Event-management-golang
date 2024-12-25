package main

import (
	"event-management/database"
	"event-management/routes"
	"event-management/structs"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&structs.User{})
	database.DB.AutoMigrate(&structs.Event{})
	// seeders.SeedUsers()
	// seeders.SeedEvents()

	router := gin.Default()
	routes.UserRoutes(router, database.DB)
	router.Run(":" + os.Getenv("PORT"))
}
