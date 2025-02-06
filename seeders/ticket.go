package seeders

import (
	"encoding/json"
	"event-management/database"
	"event-management/structs"
	"log"
	"os"
)

func SeedTickets() {
	jsonData, err := os.ReadFile("database/dummy/ticket.json")
	if err != nil {
		log.Printf("Error reading ticket.json file : %v", err)
		return
	}

	var tickets []structs.Ticket
	if err := json.Unmarshal(jsonData, &tickets); err != nil {
		log.Printf("Error parsing jsonData : %v", err)
		return
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
