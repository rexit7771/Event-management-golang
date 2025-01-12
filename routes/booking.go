package routes

import (
	"event-management/controllers"
	"event-management/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BookingRoutes(router *gin.Engine, db *gorm.DB) {
	bookingsGroup := router.Group("/bookings")
	bookingsGroup.Use(middlewares.Auth())
	{
		bookingsGroup.GET("/:id", middlewares.IsBookingTicketOwnerByParam(), controllers.GetDetailBookingByUserId)
		bookingsGroup.POST("/", controllers.AddBooking)

		bookingsGroup.Use(middlewares.IsBookingTicketOwner())
		{
			bookingsGroup.GET("/", controllers.GetAllBookingsByOwner)
			bookingsGroup.PUT("/:id", controllers.UpdateQuantity)
			bookingsGroup.PATCH("/:id", controllers.UpdateCancelled)
		}

		bookingsGroup.Use(middlewares.IsAdmin())
		{
			bookingsGroup.GET("/all", controllers.GetAllBookings)
			bookingsGroup.DELETE("/:id", controllers.DeleteBooking)
		}
	}
}
