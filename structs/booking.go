package structs

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	User_id     uint   `json:"user_id" gorm:"not null"`
	Ticket_id   uint   `json:"ticket_id" gorm:"not null"`
	Quantity    int    `json:"quantity" gorm:"not null"`
	Total_price int    `json:"total_price" gorm:"not null"`
	Cancelled   bool   `json:"cancelled" gorm:"default:false"`
	User        User   `gorm:"foreignKey:User_id;references:ID"`
	Ticket      Ticket `gorm:"foreignKey:Ticket_id;references:ID"`
}
