package main

import (
	"log"

	"github.com/Room-Balance/room-balance-backend.git/db"
)

func main() {
	log.Println("Starting Room Balance Backend...")

	// Connect to database
	db.ConnectDB()

	log.Println("Application started successfully!")

	db.MigrateDB()

}
