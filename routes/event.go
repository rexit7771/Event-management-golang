package routes

import (
	"event-management/controllers"
	"event-management/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EventRoutes(router *gin.Engine, db *gorm.DB) {
	router.GET("/events", controllers.GetAllApprovedEvents)
	eventsGroup := router.Group("/events")
	eventsGroup.Use(middlewares.Auth())
	{
		eventsGroup.POST("/", controllers.AddEvent)
		eventsGroup.Use(middlewares.IsEventOwner())
		{
			eventsGroup.PUT("/:id", controllers.UpdateEvent)
			eventsGroup.DELETE("/:id", middlewares.IsEventOwner(), controllers.DeleteEvent)
		}

		eventsGroup.Use(middlewares.IsAdmin())
		{
			eventsGroup.GET("/all", controllers.GetAllEvents)
			eventsGroup.PATCH("/:id", controllers.UpdateApproval)
		}
	}
}
