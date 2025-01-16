package routes

import (
	"event-management/controllers"
	"event-management/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/register", controllers.AddNewUser)
	router.POST("/login", controllers.Login)
	userGroup := router.Group("/users")
	{
		userGroup.Use(middlewares.Auth())
		userGroup.GET("/", controllers.GetUserByToken)
		userGroup.GET("/all", middlewares.IsAdmin(), controllers.GetAllUser)
		userGroup.Use(middlewares.IsAccountOwner())
		userGroup.GET("/:id", controllers.GetUserById)
		userGroup.PUT("/:id", controllers.UpdateUserById)
		userGroup.DELETE("/:id", controllers.DeleteUserById)
	}

}
