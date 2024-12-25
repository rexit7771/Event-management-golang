package seeders

import (
	"event-management/database"
	"event-management/structs"
	"log"
)

func SeedEvents() {
	events := []structs.Event{
		{Title: "Linkin Park Concert", Description: "Linkin Park Concert at Jakarta, They'll gonna introduce us to their new singer", Date: "2025-02-08", Location: "Gelora Bung Karno, Jakarta, Indonesia", Created_by: 1},
		{Title: "Tech Conference 2025", Description: "Annual Tech Conference focusing on AI and Machine Learning", Date: "2025-05-15", Location: "Moscone Center, San Francisco, USA", Created_by: 2},
		{Title: "Art Expo 2025", Description: "International Art Expo showcasing modern art", Date: "2025-09-20", Location: "Louvre Museum, Paris, France", Created_by: 3},
	}

	for _, event := range events {
		result := database.DB.Create(&event)
		if result.Error != nil {
			log.Printf("Failed to seed event %s: %v", event.Title, result.Error)
		} else {
			log.Printf("Event %s seeded successfully!", event.Title)
		}
	}
}
