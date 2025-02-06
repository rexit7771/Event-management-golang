package helpers

import (
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

func QuerySearch(query *gorm.DB, c *gin.Context) (searchEvent string, searchLocation string) {
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
