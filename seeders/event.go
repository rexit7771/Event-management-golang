package seeders

import (
	"encoding/json"
	"event-management/database"
	"event-management/structs"
	"log"
	"os"
)

func SeedEvents() {
	jsonData, err := os.ReadFile("database/dummy/event.json")
	if err != nil {
		log.Printf("Error reading event.json file : %v", err)
		return
	}

	var events []structs.Event
	if err := json.Unmarshal(jsonData, &events); err != nil {
		log.Printf("Error parsing event.json: %v", err)
	}

	for _, event := range events {
		result := database.DB.Where(&structs.Event{
			Title:       event.Title,
			Image_url:   event.Image_url,
			Description: event.Description,
			Date:        event.Date,
			Location:    event.Location,
			Created_by:  event.Created_by,
			Approved:    event.Approved,
		}).FirstOrCreate(&event)
		if result.Error != nil {
			log.Printf("Failed to seed event %s: %v", event.Title, result.Error)
		} else {
			log.Printf("Event %s seeded successfully!", event.Title)
		}
	}
}
