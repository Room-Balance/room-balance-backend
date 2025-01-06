package main

import (
	"log"
	"net/http"

	"github.com/Room-Balance/room-balance-backend.git/db"
	"github.com/Room-Balance/room-balance-backend.git/middlewares"
	"github.com/Room-Balance/room-balance-backend.git/routes"
)

func main() {
	// Connect to the database
	db.ConnectDB()
	db.MigrateDB()

	// Initialize Firebase
	middlewares.InitFirebase()

	// Initialize routes
	router := routes.RegisterRoutes()

	// Start the HTTP server
	port := ":8080"
	log.Printf("Server running on http://localhost%s", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
