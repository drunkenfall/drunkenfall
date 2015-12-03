package main

import (
	"github.com/boltdb/bolt"
	"log"
)

// Database is the persisting class
type Database struct {
	DB            *bolt.DB
	Tournaments   []*Tournament
	tournamentRef map[string]*Tournament
}

var (
	// TournamentKey is the byte string identifying the tournament buckets
	TournamentKey = []byte("tournaments")
)

// NewDatabase returns a new database object
func NewDatabase(fn string) (*Database, error) {
	bolt, err := bolt.Open(fn, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := &Database{DB: bolt}

	return db, nil
}

// Persist stores the current state of the tournaments into the db
func (d *Database) Persist(t *Tournament) error {
	ret := d.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(TournamentKey)
		if err != nil {
			return err
		}

		json, _ := t.JSON()
		err = b.Put([]byte(t.ID), json)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	return ret
}

// Close closes the database
func (d *Database) Close() error {
	return d.DB.Close()
}
