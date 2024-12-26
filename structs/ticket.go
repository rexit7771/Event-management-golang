package structs

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	Event_id uint   `json:"event_id" gorm:"not null"`
	Type     string `json:"type" gorm:"not null"`
	Price    int    `json:"price" gorm:"not null"`
	Quantity int    `json:"quantity" gorm:"not null"`
	Event    Event  `gorm:"foreignKey:Event_id;references:ID"`
}
