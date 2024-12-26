package main

import (
	"event-management/database"
	"event-management/routes"
	"event-management/seeders"
	"event-management/structs"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&structs.User{})
	database.DB.AutoMigrate(&structs.Event{})
	database.DB.AutoMigrate(&structs.Ticket{})

	// seeders.SeedUsers()
	// seeders.SeedEvents()
	seeders.SeedTickets()

	router := gin.Default()
	routes.UserRoutes(router, database.DB)
	routes.EventRoutes(router, database.DB)
	router.Run(":" + os.Getenv("PORT"))
}
