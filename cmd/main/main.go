package main

import (
	"log"
	"net/http"
	"restoportGo/internal/routes"
	"restoportGo/internal/services"
)

func main() {

	db, err := services.EstablishConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new router
	router := routes.NewRouter(db)

	// Start the server
	log.Println("Server running on http://localhost:8000")
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
