package seeders

import (
	"event-management/database"
	"event-management/structs"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	users := []structs.User{
		{Name: "Admin", Email: "admin@mail.com", Password: "admin", Role: "admin"},
		{Name: "Bagus", Email: "pramaskoro@gmail.com", Password: "bagus", Role: "admin"},
		{Name: "Fakhry", Email: "fakhry@gmail.com", Password: "fakhry"},
	}
	for _, user := range users {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			panic(err)
		}
		user.Password = string(hashedPassword)

		result := database.DB.Create(&user)
		if result.Error != nil {
			log.Printf("Failed to seed user %s: %v", user.Name, result.Error)
		} else {
			log.Printf("User %s seeded successfully!", user.Name)
		}
	}
}
