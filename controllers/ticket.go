package controllers

import (
	"event-management/database"
	"event-management/helpers"
	"event-management/middlewares"
	"event-management/structs"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllApprovedEventsTickets(c *gin.Context) {
	var tickets []structs.Ticket
	query := database.DB.Model(&structs.Ticket{})
	page, limit, offset := helpers.QueryPagination(c)
	searchEventId, searchTicketType := helpers.QueryTicket(query, c)

	cacheKey := fmt.Sprintf("tickets:page:%d:limit:%d:eventId:%d:ticket:%s", page, limit, searchEventId, searchTicketType)
	err := helpers.CheckCache(cacheKey, c)
	if err == nil {
		return
	}

	var totalRows int64
	query.Count(&totalRows)

	query.
		Preload("Event.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Offset(offset).
		Limit(limit).
		Find(&tickets)

	var approvedTickets []structs.Ticket
	for _, ticket := range tickets {
		if ticket.Event.Approved {
			approvedTickets = append(approvedTickets, ticket)
		}
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pagination := helpers.PaginationFormat(page, limit, totalRows, totalPages, approvedTickets)
	c.JSON(http.StatusOK, gin.H{"result": pagination})
}

func GetAllTickets(c *gin.Context) {
	var tickets []structs.Ticket
	query := database.DB.Model(&structs.Ticket{})
	page, limit, offset := helpers.QueryPagination(c)
	searchEventId, searchTicketType := helpers.QueryTicket(query, c)

	cacheKey := fmt.Sprintf("tickets:page:%d:limit:%d:eventId:%d:ticket:%s",
		page, limit, searchEventId, searchTicketType)

	err := helpers.CheckCache(cacheKey, c)
	if err == nil {
		return
	}

	var totalRows int64
	query.Count(&totalRows)

	query.Preload("Event").
		Preload("Event.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Offset(offset).
		Limit(limit).
		Find(&tickets)

	totalPages := int(math.Ceil(float64(totalRows) / float64(limit)))
	pagination := helpers.PaginationFormat(page, limit, totalRows, totalPages, tickets)
	c.JSON(http.StatusOK, gin.H{"result": pagination})
}

func GetTicketById(c *gin.Context) {
	ticketID := c.Param("id")
	var ticket structs.Ticket
	database.DB.Table("tickets").
		Preload("Event", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Preload("Event.User").Find(&ticket, ticketID)
}

func GetTicketsByEventParam(c *gin.Context) {
	eventID := c.Param("eventId")
	var tickets []structs.Ticket
	database.DB.
		Where("event_id = ?", eventID).
		Preload("Event").
		Preload("Event.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Find(&tickets)
	c.JSON(http.StatusOK, gin.H{"result": tickets})
}

func AddTicket(c *gin.Context) {
	if err = database.DB.Create(&middlewares.Ticket).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Ticket has been created"})
}

func UpdateTicket(c *gin.Context) {
	ticketId := c.Param("id")
	var ticket structs.Ticket
	database.DB.Table("tickets").Preload("Event").Preload("User").First(&ticket, ticketId)

	var ticketUpdate structs.Ticket
	err := c.ShouldBind(&ticketUpdate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := database.DB.Model(&ticket).Updates(&ticketUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket has been updated"})
}

func DeleteTicket(c *gin.Context) {
	ticketId := c.Param("id")
	var ticket structs.Ticket
	if err := database.DB.First(&ticket, ticketId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Ticket is not found"})
		return
	}

	if err := database.DB.Delete(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket has been deleted"})
}
