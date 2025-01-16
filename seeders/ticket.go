package seeders

import (
	"event-management/database"
	"event-management/structs"
	"log"
)

func SeedTickets() {
	tickets := []structs.Ticket{
		{Event_id: 1, Type: "Reguler", Price: 600000, Quantity: 250},
		{Event_id: 1, Type: "VIP", Price: 4000000, Quantity: 100},
		{Event_id: 2, Type: "Reguler", Price: 400000, Quantity: 400},
		{Event_id: 3, Type: "Reguler", Price: 800000, Quantity: 300},
	}

	for _, ticket := range tickets {
		result := database.DB.Where(&structs.Ticket{
			Event_id: ticket.Event_id,
			Type:     ticket.Type,
			Price:    ticket.Price,
			Quantity: ticket.Quantity,
		}).FirstOrCreate(&ticket)
		if result.Error != nil {
			log.Printf("Failed to seed ticket %d: %v", ticket.Event_id, result.Error)
		} else {
			log.Printf("Ticket %d seeded successfully!", ticket.ID)
		}
	}
}
