package towerfall

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/boltdb/bolt"
)

// Bolt is the persisting class
type Bolt struct {
	Server      *Server
	Tournaments []*Tournament
	People      []*Person
}

type boltWriter struct {
	DB *bolt.DB
}

type boltReader struct {
	DB *bolt.DB
}

var (
	// TournamentKey defines the tournament buckets
	TournamentKey = []byte("tournaments")
	// PeopleKey defines the bucket of people
	PeopleKey = []byte("people")
	// MigrationKey defines the bucket of migration levels
	MigrationKey = []byte("migration")
)

var tournamentMutex = &sync.Mutex{}
var personMutex = &sync.Mutex{}

// Used to signal that a current tournament was found and that the
// scanner should stop iterating.
var ErrTournamentFound = errors.New("found")

// SaveTournament stores the current state of the tournaments into the db
func (d boltWriter) saveTournament(t *Tournament) error {
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

	// go d.Server.SendWebsocketUpdate("tournament", t)
	return ret
}

// OverwriteTournament takes a new foreign Tournament{} object and replaces
// the one with the same ID with that one.
//
// Used from the EditHandler()
func (d boltWriter) overwriteTournament(t *Tournament) error {
	ret := d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TournamentKey)

		json, err := t.JSON()
		if err != nil {
			log.Fatal(err)
		}

		err = b.Put([]byte(t.ID), json)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	return ret
}

// SavePerson stores a person into the DB
func (d boltWriter) savePerson(p *Person) error {
	err := d.DB.Update(func(tx *bolt.Tx) error {
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

	return err
}

// GetPerson gets a Person{} from the DB
func (d boltReader) getPerson(id string) (*Person, error) {
	tx, err := d.DB.Begin(false)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket(PeopleKey)
	if b == nil {
		return nil, errors.New("database not initialized")
	}
	out := b.Get([]byte(id))
	if out == nil {
		return &Person{}, errors.New("user not found")
	}
	p := &Person{}
	_ = json.Unmarshal(out, p)
	return p, nil
}

// GetSafePerson gets a Person{} from the DB, while being absolutely
// sure there will be no error.
//
// This is only for hardcoded cases where error handling is just pointless.
func (d boltReader) getSafePerson(id string) *Person {
	p, _ := d.getPerson(id)
	return p
}

// DisablePerson disables or re-enables a person
func (d boltReader) disablePerson(id string) error {
	// p, err := d.getPerson(id)
	// if err != nil {
	// 	return err
	// }

	// p.Disabled = !p.Disabled
	// d.SavePerson(p)

	return nil
}

// LoadPeople loads the people from the database
func (d boltReader) getPeople() ([]*Person, error) {
	ret := make([]*Person, 0)
	err := d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PeopleKey)
		if b == nil {
			return nil // Test setup - no profiles
		}

		err := b.ForEach(func(k []byte, v []byte) error {
			p, err := LoadPerson(v)

			// If the player is disabled, just skip them
			if err == ErrPlayerDisabled {
				return nil
			}

			if err != nil {
				return err
			}

			ret = append(ret, p)
			return nil
		})
		return err
	})

	return ret, err
}

// getTournament gets a tournament by ID
func (d boltReader) getTournament(id string, s *Server) (*Tournament, error) {
	tx, err := d.DB.Begin(false)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket(TournamentKey)
	if b == nil {
		return nil, errors.New("database not initialized")
	}
	out := b.Get([]byte(id))
	if out == nil {
		return &Tournament{}, errors.New("user not found")
	}
	return LoadTournament(out, s)
}

func (d boltReader) getTournaments(s *Server) ([]*Tournament, error) {
	ret := make([]*Tournament, 0)
	err := d.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TournamentKey)
		if b == nil {
			return nil
		}

		err := b.ForEach(func(k []byte, v []byte) error {
			t, err := LoadTournament(v, s)
			if err != nil {
				return err
			}

			ret = append(ret, t)
			return nil
		})
		return err
	})

	return ret, err
}

// GetCurrentTournament gets the currently running tournament.
//
// Returns the first matching one, so if there are multiple they will
// be shadowed.
func (d boltReader) getCurrentTournament(s *Server) (*Tournament, error) {
	ts, err := d.getTournaments(s)
	if err != nil {
		return &Tournament{}, err
	}
	for _, t := range SortByScheduleDate(ts) {
		if t.IsRunning() {
			return t, nil
		}
	}
	return &Tournament{}, errors.New("no tournament is running")
}

// ClearTestTournaments deletes any tournament that doesn't begin with "DrunkenFall"
func (d boltWriter) clearTestTournaments(s *Server) error {
	err := d.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TournamentKey)

		err := b.ForEach(func(k []byte, v []byte) error {
			t, err := LoadTournament(v, s)
			if err != nil {
				return err
			}

			if !strings.HasPrefix(t.Name, "DrunkenFall") {
				log.Print("Deleting ", t.ID)
				err := b.Delete([]byte(t.ID))
				if err != nil {
					return err
				}
			}
			return nil

		})
		return err
	})

	log.Print("Not sending full update; not implemented")
	// d.Server.SendWebsocketUpdate("all", d.asMap())

	return err
}

// func (d boltReader) asMap() map[string]*Tournament {
// 	tournamentMutex.Lock()
// 	out := make(map[string]*Tournament)
// 	for _, t := range d.Tournaments {
// 		out[t.ID] = t
// 	}
// 	tournamentMutex.Unlock()
// 	return out
// }

// Close closes the database
func (d boltReader) close() error {
	return d.DB.Close()
}

// Close closes the database
func (d boltWriter) close() error {
	return d.DB.Close()
}

// ByScheduleDate is a sort.Interface that sorts tournaments according
// to when they were scheduled.
type ByScheduleDate []*Tournament

func (s ByScheduleDate) Len() int {
	return len(s)
}
func (s ByScheduleDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByScheduleDate) Less(i, j int) bool {
	return s[i].Scheduled.Before(s[j].Scheduled)
}

// SortByScheduleDate returns a list in order of schedule date
func SortByScheduleDate(ps []*Tournament) []*Tournament {
	sort.Sort(ByScheduleDate(ps))
	return ps
}
