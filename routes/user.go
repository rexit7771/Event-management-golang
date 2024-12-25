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
	router.Use(middlewares.Auth())
	router.GET("/users", middlewares.IsAdmin(), controllers.GetAllUser)
	router.Use(middlewares.IsAccountOwner())
	router.GET("/users/:id", controllers.GetUserById)
	router.PUT("/users/:id", controllers.UpdateUserById)
	router.DELETE("/users/:id", controllers.DeleteUserById)
}
