package controllers

import (
	"event-management/database"
	"event-management/structs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllEvents(c *gin.Context) {
	var events []structs.Event
	database.DB.Table("events").Preload("users").Find(&events)
	c.JSON(http.StatusOK, gin.H{
		"result": events,
	})
}

func GetAllApprovedEvents(c *gin.Context) {
	var events []structs.Event
	database.DB.Where("approved = ?", true).Preload("users").Find(&events)
	c.JSON(http.StatusOK, gin.H{
		"result": events,
	})
}

func AddEvent(c *gin.Context) {
	var event structs.Event
	err := c.ShouldBind(&event)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User ID is not found"})
		return
	}
	userIDUint := userID.(uint)
	event.Created_by = userIDUint

	if err = database.DB.Create(&event).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": event.Title + " Event has been created",
	})
}

func UpdateEvent(c *gin.Context) {
	idParam := c.Param("id")
	var event structs.Event
	database.DB.Table("events").Preload("users").First(&event, idParam)
	if event.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event is not found"})
		return
	}

	var eventUpdate structs.Event
	err := c.ShouldBind(&eventUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := database.DB.Model(&event).Updates(eventUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event has been updated",
	})
}

func UpdateApproval(c *gin.Context) {
	eventId := c.Param("id")
	var event structs.Event
	database.DB.Table("events").Preload("users").First(&event, eventId)
	if event.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event is not found"})
		return
	}

	var approvalUpdate struct {
		Approved bool `json:"approved"`
	}
	if err := c.ShouldBindJSON(&approvalUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := database.DB.Model(&event).Update("approved", approvalUpdate.Approved).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event has been updated",
	})
}

func DeleteEvent(c *gin.Context) {
	idParam := c.Param("id")
	var event structs.Event
	if err := database.DB.First(&event, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Event is not found"})
		return
	}

	if err := database.DB.Delete(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Event has been deleted",
	})
}
