package routes

import (
	"event-management/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/register", controllers.AddNewUser)
	router.POST("/login", controllers.Login)
}
