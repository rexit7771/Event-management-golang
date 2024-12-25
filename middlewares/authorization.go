package middlewares

import (
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
