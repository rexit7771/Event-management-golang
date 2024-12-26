package seeders

import (
	"event-management/database"
	"event-management/structs"
	"log"
)

func SeedEvents() {
	events := []structs.Event{
		{Title: "Linkin Park Concert", Image_url: "https://awsimages.detik.net.id/community/media/visual/2024/11/26/poster-konser-linkin-park-di-jakarta_34.jpeg?w=700&q=90", Description: "Linkin Park Concert at Jakarta, They'll gonna introduce us to their new singer", Date: "2025-02-08", Location: "Gelora Bung Karno, Jakarta, Indonesia", Created_by: 1, Approved: true},
		{Title: "Tech Conference 2025", Image_url: "https://img.evbuc.com/https%3A%2F%2Fcdn.evbuc.com%2Fimages%2F873765669%2F167645577512%2F1%2Foriginal.20241014-084309?crop=focalpoint&fit=crop&w=1000&auto=format%2Ccompress&q=75&sharp=10&fp-x=0.5&fp-y=0.5&s=39efc273e967b744b2aad0311456f8ff", Description: "Annual Tech Conference focusing on AI and Machine Learning", Date: "2025-05-15", Location: "Moscone Center, San Francisco, USA", Created_by: 2, Approved: false},
		{Title: "Art Expo 2025", Image_url: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSf9wvhnW6khbamV6ORq5B6J8hKjPRphWUMOw&s", Description: "International Art Expo showcasing modern art", Date: "2025-09-20", Location: "Louvre Museum, Paris, France", Created_by: 3, Approved: true},
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
