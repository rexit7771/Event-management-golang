package main

import (
	"event-management/database"
	"event-management/seeders"
	"event-management/structs"
)

func main() {
	database.Connect()
	database.DB.AutoMigrate(&structs.User{})
	seeders.SeedUsers()
}
