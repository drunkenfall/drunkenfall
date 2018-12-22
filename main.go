package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/one-eye/drunkenfall/towerfall"
)

func main() {
	Kekkonen() // Kekkonen

	// Load the configuration
	config := towerfall.ParseConfig()
	config.Print()

	// Instantiate the database
	db, err := towerfall.NewDatabase(config)
	if err != nil {
		log.Fatal(err)
	}

	// Apply any applicable migrations
	// err = migrations.Migrate(db.DB)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Create the server instance
	s := towerfall.NewServer(config, db)

	// Start the live updater
	listener, err := towerfall.NewListener(config, db)
	if err != nil {
		log.Fatal(err)
	}
	go listener.Serve()

	// Catch termination signals so we can close the databas properly
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Print("Catching SIGTERM, closing database...")
		db.Close()
		log.Print("Done. Exiting.")
		os.Exit(1)
	}()

	// Actually start serving
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}
