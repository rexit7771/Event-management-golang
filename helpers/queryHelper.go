package helpers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func QueryPagination(c *gin.Context) (page int, limit int, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset = (page - 1) * limit
	return
}

func CountTotalPages(totalRows int64, limit int) int {
	return int(math.Ceil(float64(totalRows) / float64(limit)))
}

func QueryEvent(query *gorm.DB, c *gin.Context) (searchEvent string, searchLocation string) {
	searchEvent = c.DefaultQuery("event", "")
	searchLocation = c.DefaultQuery("location", "")

	if searchEvent != "" {
		searchEventQuery := "%" + searchEvent + "%"
		query = query.Where("title ILIKE ?", searchEventQuery)
	}

	if searchLocation != "" {
		searchLocationQuery := "%" + searchLocation + "%"
		query = query.Where("location ILIKE ?", searchLocationQuery)
	}
	return
}

func QueryUser(query *gorm.DB, c *gin.Context) (searchName string, searchEmail string) {
	searchName = c.DefaultQuery("name", "")
	searchEmail = c.DefaultQuery("email", "")

	if searchName != "" {
		searchNameQuery := "%" + searchName + "%"
		query = query.Where("name ILIKE ?", searchNameQuery)
	}

	if searchEmail != "" {
		searchEmailQuery := "%" + searchEmail + "%"
		query = query.Where("email ILIKE ?", searchEmailQuery)
	}
	return
}

func QueryTicket(query *gorm.DB, c *gin.Context) (searchEventId int, searchTicketType string) {
	searchEventId, _ = strconv.Atoi(c.DefaultQuery("eventId", "0"))
	searchTicketType = c.DefaultQuery("ticket", "")

	if searchEventId != 0 {
		query = query.Where("event_id ILIKE ?", searchEventId)
	}

	if searchTicketType != "" {
		query = query.Where("type ILIKE ?", searchTicketType)
	}
	return
}

func QueryBooking(query *gorm.DB, c *gin.Context) (searchTicketId int, searchUserId int, searchCancelled string) {
	searchTicketId, _ = strconv.Atoi(c.DefaultQuery("ticketId", "0"))
	searchUserId, _ = strconv.Atoi(c.DefaultQuery("userId", "0"))
	searchCancelled = c.DefaultQuery("cancelled", "false")

	if searchTicketId != 0 {
		query = query.Where("ticket_id = ?", searchTicketId)
	}

	if searchUserId != 0 {
		query = query.Where("user_id = ?", searchUserId)
	}

	if searchCancelled != "false" {
		query = query.Where("cancelled = ?", true)
	}

	return
}
