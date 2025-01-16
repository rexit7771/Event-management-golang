package main

import (
	"event-management/database"
	"event-management/middlewares"
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
	database.DB.AutoMigrate(&structs.Booking{})

	seeders.SeedUsers()
	seeders.SeedEvents()
	seeders.SeedTickets()

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	routes.UserRoutes(router, database.DB)
	routes.EventRoutes(router, database.DB)
	routes.TicketRoutes(router, database.DB)
	routes.BookingRoutes(router, database.DB)
	router.Run(":" + os.Getenv("PORT"))
}
