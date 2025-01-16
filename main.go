package main

import (
	"event-management/database"
	"event-management/routes"
	"event-management/seeders"
	"event-management/structs"
	"os"

	"github.com/gin-contrib/cors"
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

	corsConfig := cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:5173"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router.Use(cors.New(corsConfig))

	routes.UserRoutes(router, database.DB)
	routes.EventRoutes(router, database.DB)
	routes.TicketRoutes(router, database.DB)
	routes.BookingRoutes(router, database.DB)
	router.Run(":" + os.Getenv("PORT"))
}
