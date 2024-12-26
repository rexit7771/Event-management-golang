package controllers

import (
	"event-management/database"
	"event-management/middlewares"
	"event-management/structs"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllApprovedEventsTickets(c *gin.Context) {
	var tickets []structs.Ticket
	database.DB.Table("tickets").Preload("Event", "approved = ?", true).Preload("Event.User").Find(&tickets)
	var approvedTickets []structs.Ticket
	for _, ticket := range tickets {
		if ticket.Event.Approved {
			approvedTickets = append(approvedTickets, ticket)
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": approvedTickets})
}

func GetAllTickets(c *gin.Context) {
	var tickets []structs.Ticket
	database.DB.Table("tickets").Preload("Event").Preload("Event.User").Find(&tickets)
	c.JSON(http.StatusOK, gin.H{"result": tickets})
}

func AddTicket(c *gin.Context) {
	if err = database.DB.Create(&middlewares.Ticket).Error; err != nil {
		fmt.Println("Ini error waktu create ticket ke database")
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
