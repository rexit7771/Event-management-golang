package main

import (
	"event-management/database"
	"event-management/helpers"
	"event-management/middlewares"
	"event-management/routes"
	"event-management/seeders"
	"event-management/structs"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	database.DB.Migrator().DropTable(&structs.User{})
	database.DB.Migrator().DropTable(&structs.Event{})
	database.DB.Migrator().DropTable(&structs.Ticket{})
	database.DB.Migrator().DropTable(&structs.Booking{})
	database.DB.AutoMigrate(&structs.User{})
	database.DB.AutoMigrate(&structs.Event{})
	database.DB.AutoMigrate(&structs.Ticket{})
	database.DB.AutoMigrate(&structs.Booking{})

	helpers.InitRedis()

	// if err := helpers.InitRabbitMQ(); err != nil {
	// 	log.Printf("RabbitMQ initialization failed: %v", err)
	// }

	// if helpers.RabbitMQChannel != nil {
	// 	go consumer.ConsumeEventOperations()
	// }

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
