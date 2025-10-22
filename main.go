package main

import (
	"backend-hanssen-hilman/database"
	"backend-hanssen-hilman/migrations"
	"backend-hanssen-hilman/routes"
	"backend-hanssen-hilman/util"
	"fmt"

	"log"
)

func main() {
	// Load environment variables
	util.LoadEnv()

	// Initialize database connection
	err := database.Database()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Database connection successful.")

	// Run migrations
	migrations.Migrate(database.DB)

	// Setup and run the router
	routes.SetupRoutes()
}
