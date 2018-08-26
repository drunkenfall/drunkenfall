package towerfall

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

// Postgres is the persisting class
type Postgres struct {
	DB            *sql.DB
	Server        *Server
	Tournaments   []*Tournament
	People        []*Person
	tournamentRef map[string]*Tournament
}

// NewPostgres returns a new database object
func NewPostgres(fn string) (*Postgres, error) {
	connStr := "user=postgres dbname=drunkenfall sslmode=verify-full"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	db := &Postgres{DB: conn}
	db.tournamentRef = make(map[string]*Tournament)
	db.LoadPeople()

	return db, nil
}

// LoadTournaments loads the tournaments from the database and into memory
func (d *Postgres) LoadTournaments() error {
	return nil
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Postgres) SaveTournament(t *Tournament) error {
	// If the tournament isn't already in the cache, we should add it
	found := false
	for _, ct := range d.Tournaments {
		if ct.ID == t.ID {
			found = true
			break
		}
	}
	if !found {
		log.Printf("Adding new tournament %s into the memory cache", t.ID)
		d.Tournaments = append(d.Tournaments, t)
		d.tournamentRef[t.ID] = t
	}

	go d.Server.SendWebsocketUpdate("tournament", t)
	return nil
}

// OverwriteTournament takes a new foreign Tournament{} object and replaces
// the one with the same ID with that one.
//
// Used from the EditHandler()
func (d *Postgres) OverwriteTournament(t *Tournament) error {
	return nil
}

// SavePerson stores a person into the DB
func (d *Postgres) SavePerson(p *Person) error {
	return nil
}

// GetPerson gets a Person{} from the DB
func (d *Postgres) GetPerson(id string) (*Person, error) {
	p := &Person{}
	// _ = json.Unmarshal(out, p)
	return p, nil
}

// GetSafePerson gets a Person{} from the DB, while being absolutely
// sure there will be no error.
//
// This is only for hardcoded cases where error handling is just pointless.
func (d *Postgres) GetSafePerson(id string) *Person {
	p, _ := d.GetPerson(id)
	return p
}

// DisablePerson disables or re-enables a person
func (d *Postgres) DisablePerson(id string) error {
	p, err := d.GetPerson(id)
	if err != nil {
		return err
	}

	p.Disabled = !p.Disabled
	d.SavePerson(p)

	return nil
}

// LoadPeople loads the people from the database and into memory
func (d *Postgres) LoadPeople() error {
	d.People = make([]*Person, 0)
	return nil
}

// GetCurrentTournament gets the currently running tournament.
//
// Returns the first matching one, so if there are multiple they will
// be shadowed.
func (d *Postgres) GetCurrentTournament() (*Tournament, error) {
	for _, t := range SortByScheduleDate(d.Tournaments) {
		if t.IsRunning() {
			return t, nil
		}
	}
	return &Tournament{}, errors.New("no tournament is running")
}

// ClearTestTournaments deletes any tournament that doesn't begin with "DrunkenFall"
func (d *Postgres) ClearTestTournaments() error {
	d.Tournaments = make([]*Tournament, 0)
	err := d.LoadTournaments()
	if err != nil {
		return err
	}

	d.Server.SendWebsocketUpdate("all", d.asMap())

	return err
}

func (d *Postgres) asMap() map[string]*Tournament {
	tournamentMutex.Lock()
	out := make(map[string]*Tournament)
	for _, t := range d.Tournaments {
		out[t.ID] = t
	}
	tournamentMutex.Unlock()
	return out
}

// Close closes the database
func (d *Postgres) Close() error {
	return d.DB.Close()
}
