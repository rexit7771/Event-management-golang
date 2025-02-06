package controllers

import (
	"event-management/database"
	"event-management/helpers"
	"event-management/structs"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllEvents(c *gin.Context) {
	var events []structs.Event
	query := database.DB.Model(&structs.Event{})
	page, limit, offset := helpers.QueryPagination(c)
	searchEvent, searchLocation := helpers.QuerySearch(query, c)

	cacheKey := fmt.Sprintf("events:page:%d:limit:%d:event:%d:location:%d", page, limit, searchEvent, searchLocation)
	err := helpers.CheckCache(cacheKey, c)
	if err == nil {
		return
	}

	var totalRows int64
	query.Count(&totalRows)
	query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email")
	}).
		Offset(offset).
		Limit(limit).
		Find(&events)

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pagination := structs.Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       events,
	}
	helpers.SetCache(pagination, cacheKey)

	c.JSON(http.StatusOK, gin.H{
		"result": pagination,
	})
}

func GetEventById(c *gin.Context) {
	eventID := c.Param("id")
	var event structs.Event
	database.DB.Table("events").Preload("User").First(&event, eventID)
	c.JSON(http.StatusOK, gin.H{
		"result": event,
	})
}

func GetAllApprovedEvents(c *gin.Context) {
	// TODO Tambahkan Pagination dan juga search query
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	cachedKey := fmt.Sprintf("events:page%d:limit:%d", page, limit)

	err := helpers.CheckCache(cachedKey, c)
	if err == nil {
		return
	}

	var events []structs.Event
	var totalRows int64
	database.DB.Model(&structs.Event{}).Count(&totalRows)
	database.DB.Where("approved = ?", true).
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&events)

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pagination := structs.Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
		Data:       events,
	}
	helpers.SetCache(pagination, cachedKey)

	c.JSON(http.StatusOK, gin.H{
		"result": pagination,
	})
}

func GetApprovedEvent(c *gin.Context) {
	eventID := c.Param("id")
	var event structs.Event
	database.DB.Where("approved = ?", true).Preload("User").First(&event, eventID)
	c.JSON(http.StatusOK, gin.H{
		"result": event,
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

	// message := structs.EventMessage{
	// 	EventID:   strconv.Itoa(int(event.ID)),
	// 	Action:    "create",
	// 	EventData: event,
	// 	UserID:    event.Created_by,
	// 	Timestamp: time.Now(),
	// }
	// helpers.PublishToQueue(message, "event_operations")

	helpers.InvalidateCache("events")
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

	// message := structs.EventMessage{
	// 	EventID:   strconv.Itoa(int(event.ID)),
	// 	Action:    "update",
	// 	EventData: event,
	// 	UserID:    event.Created_by,
	// 	Timestamp: time.Now(),
	// }
	// helpers.PublishToQueue(message, "event_operations")

	helpers.InvalidateCache("events")
	c.JSON(http.StatusOK, gin.H{
		"message": "Event has been updated",
	})
}

func UpdateApproval(c *gin.Context) {
	eventId := c.Param("id")
	var event structs.Event
	database.DB.Table("events").Preload("User").First(&event, eventId)
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

	// message := structs.EventMessage{
	// 	EventID:   strconv.Itoa(int(event.ID)),
	// 	Action:    "approval",
	// 	EventData: event,
	// 	UserID:    event.Created_by,
	// 	Timestamp: time.Now(),
	// }
	// helpers.PublishToQueue(message, "event_operations")

	helpers.InvalidateCache("events")

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

	// message := structs.EventMessage{
	// 	EventID:   strconv.Itoa(int(event.ID)),
	// 	Action:    "delete",
	// 	EventData: event,
	// 	UserID:    event.Created_by,
	// 	Timestamp: time.Now(),
	// }
	// helpers.PublishToQueue(message, "event_operations")

	helpers.InvalidateCache("events")

	c.JSON(http.StatusOK, gin.H{"message": "Event has been deleted"})
}
