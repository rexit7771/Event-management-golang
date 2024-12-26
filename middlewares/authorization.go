package middlewares

import (
	"event-management/database"
	"event-management/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsAccountOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User ID is not found"})
			return
		}

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User Role is not found"})
			return
		}

		roleStr := role.(string)
		if roleStr == "admin" {
			c.Next()
		} else if idParam == userID {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You are forbidden to access this",
			})
			return
		}
	}
}

func IsEventOwner() gin.HandlerFunc {
	return func(c *gin.Context) {
		eventId := c.Param("id")
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User ID is not found"})
			return
		}

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User Role is not found"})
			return
		}

		var event structs.Event
		database.DB.First(&event, eventId)

		roleStr := role.(string)
		userIDUint := userID.(uint)
		if roleStr == "admin" {
			c.Next()
		} else if userIDUint == event.Created_by {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You are forbidden to access this event",
			})
			return
		}
	}
}

func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User ID is not found"})
			return
		}

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User Role is not found"})
			return
		}

		roleStr := role.(string)
		if roleStr == "admin" {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You are forbidden to access this",
			})
			return
		}
	}
}
