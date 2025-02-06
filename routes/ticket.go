package routes

import (
	"event-management/controllers"
	"event-management/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TicketRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/tickets", controllers.GetAllApprovedEventsTickets)
	router.GET("/tickets/:id", controllers.GetTicketById)
	router.GET("/tickets/event/:eventId", controllers.GetTicketsByEventParam)
	ticketsGroup := router.Group("/tickets")
	ticketsGroup.Use(middlewares.Auth())
	{
		ticketsGroup.POST("/", middlewares.IsEventTicketOwnerByBody(), controllers.AddTicket)
		ticketsGroup.Use(middlewares.IsEventTicketOwnerByParam())
		{
			ticketsGroup.PUT("/:id", controllers.UpdateTicket)
			ticketsGroup.DELETE("/:id", controllers.DeleteTicket)
		}
		ticketsGroup.Use(middlewares.IsAdmin())
		{
			ticketsGroup.GET("/all", controllers.GetAllTickets)
		}
	}
}
