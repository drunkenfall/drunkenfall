package main

import (
	"github.com/boltdb/bolt"
	"log"
)

// Database is the persisting class
type Database struct {
	DB            *bolt.DB
	Server        *Server
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
	db.tournamentRef = make(map[string]*Tournament)

	return db, nil
}

// LoadTournaments loads the tournaments from the database and into memory
func (d *Database) LoadTournaments() error {
	err := d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TournamentKey)
		if b == nil {
			// If there is no bucket, bail silently.
			// This only really happens in tests.
			// TODO: Fix pls
			return nil
		}

		err := b.ForEach(func(k []byte, v []byte) error {
			t, err := LoadTournament(v, d)
			if err != nil {
				return err
			}

			d.Tournaments = append(d.Tournaments, t)
			d.tournamentRef[t.ID] = t
			return nil
		})
		return err
	})

	return err
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) SaveTournament(t *Tournament) error {
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
