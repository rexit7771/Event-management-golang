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

var Ticket structs.Ticket

func IsEventTicketOwnerByBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := c.ShouldBind(&Ticket)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		var event structs.Event
		database.DB.Table("events").First(&event, Ticket.Event_id)

		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, User ID is not found"})
			return
		}
		userIDUint := userID.(uint)

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, User's Role is not found"})
			return
		}
		roleStr := role.(string)

		if roleStr == "admin" {
			c.Next()
			return
		} else if userIDUint != event.Created_by {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, You are not the owner of this event"})
			return
		}

		c.Next()
	}
}

func IsEventTicketOwnerByParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketID := c.Param("id")
		var ticket structs.Ticket
		database.DB.Table("tickets").Preload("Event").First(&ticket, ticketID)

		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, User ID is not found"})
			return
		}
		userIDUint := userID.(uint)

		role, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, User's Role is not found"})
			return
		}
		roleStr := role.(string)

		if roleStr == "admin" {
			c.Next()
			return
		} else if userIDUint != ticket.Event.Created_by {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, You are not the owner of this event"})
			return
		}
		c.Next()
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
