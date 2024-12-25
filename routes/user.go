package routes

import (
	"event-management/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/register", controllers.AddNewUser)
	router.POST("/login", controllers.Login)
	router.GET("/users", controllers.GetAllUser)
	router.GET("/users/:id", controllers.GetUserById)
	router.PUT("/users/:id", controllers.UpdateUserById)
	router.DELETE("/users/:id", controllers.DeleteUserById)
}
