package main

import (
	"log"
	"os"

	"github.com/drunkenfall/drunkenfall/towerfall"
	"github.com/drunkenfall/drunkenfall/towerfall/migrations"
)

func main() {
	Kekkonen() // Kekkonen

	// Instantiate the database
	db, err := towerfall.NewDatabase(os.Getenv("DF_DB_PATH"))
	if err != nil {
		log.Fatal(err)
	}

	// Apply any applicable migrations
	err = migrations.Migrate(db.DB)
	if err != nil {
		log.Fatal(err)
	}

	// Create the server instance
	s := towerfall.NewServer(db)

	// Load tournaments from the database
	err = db.LoadTournaments()
	if err != nil {
		log.Fatal(err)
	}

	// Set up the paths and the websocket listeners
	s.RegisterHandlersAndListeners()


	// Actually start serving
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
