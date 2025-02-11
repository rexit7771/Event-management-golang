package controllers

import (
	"event-management/database"
	"event-management/helpers"
	"event-management/structs"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllBookings(c *gin.Context) {
	var bookings []structs.Booking
	query := database.DB.Model(&structs.Booking{})
	page, limit, offset := helpers.QueryPagination(c)
	searchTicketId, searchUserId, searchCancelled := helpers.QueryBooking(query, c)

	cacheKey := fmt.Sprintf("bookings:page:%d:limit:%d:ticketId:%d:userId:%d:cancelled:%s", page, limit, searchTicketId, searchUserId, searchCancelled)
	err := helpers.CheckCache(cacheKey, c)
	if err == nil {
		return
	}

	var totalRows int64
	query.Count(&totalRows)
	query.Preload("Ticket").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Preload("Ticket.Event").
		Preload("Ticket.Event.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Offset(offset).
		Limit(limit).
		Find(&bookings)

	totalPages := helpers.CountTotalPages(totalRows, limit)
	pagination := helpers.PaginationFormat(page, limit, totalRows, totalPages, bookings)
	c.JSON(http.StatusOK, gin.H{"result": pagination})
}

func GetAllBookingsByOwner(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized, You Have to login first"})
		return
	}
	userIDUint := userID.(uint)
	fmt.Print(userID)

	var bookings []structs.Booking
	query := database.DB.Model(&structs.Booking{})
	page, limit, offset := helpers.QueryPagination(c)
	searchTicketId, _, searchCancelled := helpers.QueryBooking(query, c)

	cacheKey := fmt.Sprintf("bookings:page:%d:limit:%d:ticketId:%d:cancelled:%s", page, limit, searchTicketId, searchCancelled)
	err := helpers.CheckCache(cacheKey, c)
	if err == nil {
		return
	}

	var totalRows int64
	query.Count(&totalRows)
	query.Preload("Ticket").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Preload("Ticket.Event").
		Preload("Ticket.Event.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "email")
		}).
		Where("user_id = ?", userIDUint).
		Offset(offset).
		Limit(limit).
		Find(&bookings)

	totalPages := helpers.CountTotalPages(totalRows, limit)
	pagination := helpers.PaginationFormat(page, limit, totalRows, totalPages, bookings)
	c.JSON(http.StatusOK, gin.H{"result": pagination})
}

func GetDetailBookingByUserId(c *gin.Context) {
	bookingId := c.Param("id")
	var booking structs.Booking
	if err := database.DB.Table("bookings").
		Preload("Ticket").
		Preload("User").
		Preload("Ticket.Event").
		Preload("Ticket.Event.User").
		First(&booking, bookingId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Detail booking is not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": booking})
}

func AddBooking(c *gin.Context) {
	// ? Binding json dari body {ticket_id, quantity}
	var booking structs.Booking
	err := c.ShouldBind(&booking)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User ID is not found"})
		return
	}
	userIDUint := userID.(uint)
	booking.User_id = userIDUint

	// ? Get Ticket berdasarkan ticket_id
	var ticket structs.Ticket
	database.DB.Table("tickets").Preload("Event").First(&ticket, &booking.Ticket_id)
	// ? Cek jika tiket yang dicari ada
	if ticket.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ticket not found"})
		return

		// ? Cek jika tiket sudah kosong
	} else if ticket.Quantity == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ticket has been sold out"})
		return
	} else if ticket.Quantity < booking.Quantity {
		qtyStr := strconv.Itoa(ticket.Quantity)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "The tickets are only available for " + qtyStr})
		return
	}
	ticket.Quantity -= booking.Quantity
	booking.Total_price = ticket.Price * booking.Quantity

	if err := database.DB.Model(&ticket).Update("quantity", ticket.Quantity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	bookingQtyStr := strconv.Itoa(booking.Quantity)
	c.JSON(http.StatusCreated, gin.H{"message": bookingQtyStr + " Tickets has been booked successfully"})
}

func UpdateQuantity(c *gin.Context) {
	bookingId := c.Param("id")
	var booking structs.Booking
	database.DB.Table("bookings").Preload("Ticket").Preload("User").First(&booking, bookingId)
	if booking.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Booking Ticket is not found"})
		return
	}

	var bookingUpdate struct {
		Quantity    int `json:"quantity"`
		Total_price int `json:"total_price"`
	}
	if err := c.ShouldBind(&bookingUpdate); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	var ticket structs.Ticket
	database.DB.Table("tickets").Preload("Event").First(&ticket, booking.Ticket_id)
	if ticket.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ticket is not found"})
		return
	} else if ticket.Quantity == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Ticket has been sold out"})
		return
	}

	var totalBookingUpdate int
	if bookingUpdate.Quantity > booking.Quantity {
		// ? Check jika jumlah ticket yang tersedia dibawah permintaan update
		if ticket.Quantity < bookingUpdate.Quantity {
			qtyStr := strconv.Itoa(ticket.Quantity)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "The tickets are only available for " + qtyStr})
			return
		}

		// ? Menghitung selisih perbedaan
		totalBookingUpdate = bookingUpdate.Quantity - booking.Quantity
		ticket.Quantity -= totalBookingUpdate

		// ? Update jumlah ticket ke database
		if err := database.DB.Model(&ticket).Update("quantity", ticket.Quantity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		// ? Update jumlah booking ticket ke database
		bookingUpdate.Total_price = ticket.Price * bookingUpdate.Quantity
		if err := database.DB.Model(&booking).Updates(bookingUpdate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking ticket quantity has been Increased"})

	} else if bookingUpdate.Quantity < booking.Quantity {
		// ? Menghitung selisih perbedaannya
		totalBookingUpdate = booking.Quantity - bookingUpdate.Quantity
		// ? Menambahkan jumlah tiket karena permintaan booking ticket dikurangi
		ticket.Quantity += totalBookingUpdate

		if err := database.DB.Model(&ticket).Update("quantity", ticket.Quantity).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		bookingUpdate.Total_price = ticket.Price * bookingUpdate.Quantity
		if err := database.DB.Model(&booking).Updates(bookingUpdate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking ticket quantity has been Decreased"})
	}
}

func UpdateCancelled(c *gin.Context) {
	bookingID := c.Param("id")
	var booking structs.Booking
	database.DB.Table("bookings").Preload("Ticket").Preload("User").First(&booking, bookingID)
	if booking.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Booking Ticket is not found"})
		return
	}

	var bookingCancelled struct {
		Cancelled bool `json:"cancelled" gorm:"default:true"`
	}
	err := c.ShouldBind(&bookingCancelled)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var ticket structs.Ticket
	database.DB.Table("tickets").Preload("Event").First(&ticket, booking.Ticket_id)
	if ticket.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ticket is not found"})
		return
	}

	ticket.Quantity += booking.Quantity
	if err := database.DB.Model(&ticket).Update("quantity", &ticket.Quantity).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	bookingCancelled.Cancelled = true
	if err := database.DB.Model(&booking).Update("cancelled", bookingCancelled.Cancelled).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking Ticket has been cancelled"})
}

func DeleteBooking(c *gin.Context) {
	bookingID := c.Param("id")
	var booking structs.Booking
	database.DB.Table("bookings").Preload("Ticket").Preload("User").Preload("Ticket.Event").First(&booking, bookingID)
	if booking.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Booking Ticket is not found"})
		return
	}

	var ticket structs.Ticket
	database.DB.Table("tickets").Preload("Event").First(&ticket, booking.Ticket_id)
	if ticket.ID == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Ticket is not found"})
		return
	}

	ticket.Quantity += booking.Quantity
	if err := database.DB.Model(&ticket).Update("quantity", &ticket.Quantity).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if err := database.DB.Delete(&booking).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Booking Ticket has been deleted"})
}
