package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
)

// Database is the persisting class
type Database struct {
	DB            *bolt.DB
	Server        *Server
	Tournaments   []*Tournament
	People        []*Person
	tournamentRef map[string]*Tournament
}

var (
	// TournamentKey defines the tournament buckets
	TournamentKey = []byte("tournaments")
	// PeopleKey defines the bucket of people
	PeopleKey = []byte("people")
	// MigrationKey defines the bucket of migration levels
	MigrationKey = []byte("migration")
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

// SavePerson stores a person into the DB
func (d *Database) SavePerson(p *Person) error {
	ret := d.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(PeopleKey)
		if err != nil {
			return err
		}

		json, _ := p.JSON()
		err = b.Put([]byte(p.ID), json)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	return ret
}

// GetPerson gets a Person{} from the DB
func (d *Database) GetPerson(id string) *Person {
	tx, err := d.DB.Begin(false)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer tx.Rollback()

	b := tx.Bucket(PeopleKey)
	out := b.Get([]byte(id))
	p := &Person{}
	_ = json.Unmarshal(out, p)
	return p
}

// LoadPeople loads the people from the database and into memory
func (d *Database) LoadPeople() error {
	d.People = make([]*Person, 0)
	err := d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PeopleKey)

		err := b.ForEach(func(k []byte, v []byte) error {
			p, err := LoadPerson(v)
			if err != nil {
				return err
			}

			d.People = append(d.People, p)
			return nil
		})
		return err
	})

	return err
}

// Close closes the database
func (d *Database) Close() error {
	return d.DB.Close()
}
