package towerfall

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type Database struct {
	Server *Server
	reader dbReader
	writer dbWriter
}

type dbReader interface {
	getPeople() ([]*Person, error)
	getTournament(id string, server *Server) (*Tournament, error)
	getTournaments(server *Server) ([]*Tournament, error)
	getCurrentTournament(server *Server) (*Tournament, error)
	getPerson(id string) (*Person, error)
	getSafePerson(id string) *Person

	close() error
}

type dbWriter interface {
	saveTournament(t *Tournament) error
	overwriteTournament(t *Tournament) error
	clearTestTournaments(s *Server) error

	savePerson(p *Person) error

	close() error
}

// LoadTournament loads a tournament from persisted JSON data
func LoadTournament(data []byte, s *Server) (t *Tournament, err error) {
	t = &Tournament{}
	err = json.Unmarshal(data, t)
	if err != nil {
		log.Print(err)
		return t, err
	}

	t.server = s
	t.db = s.DB

	t.SetMatchPointers()
	return
}

// NewDatabase returns a new database object
func NewDatabase(c *Config) (*Database, error) {
	bolt, err := bolt.Open(c.DbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	r := boltReader{bolt}
	w := boltWriter{bolt}

	if c.DbReader == "postgres" {
		log.Print("not switching to postgres yet")
	}

	db := &Database{reader: r, writer: w}

	return db, nil
}

// GetTournament gets a single tournament
func (d *Database) GetTournament(id string, s *Server) (*Tournament, error) {
	return d.reader.getTournament(id, s)
}

// LoadTournaments loads the tournaments from the database and into memory
func (d *Database) GetTournaments(s *Server) ([]*Tournament, error) {
	return d.reader.getTournaments(s)
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) SaveTournament(t *Tournament) error {
	err := d.writer.saveTournament(t)
	go d.Server.SendWebsocketUpdate("tournament", t)
	return err
}

// OverwriteTournament takes a new foreign Tournament{} object and replaces
// the one with the same ID with that one.
//
// Used from the EditHandler()
func (d *Database) OverwriteTournament(t *Tournament) error {
	return d.writer.overwriteTournament(t)
}

// SavePerson stores a person into the DB
func (d *Database) SavePerson(p *Person) error {
	return d.writer.savePerson(p)
}

// GetPerson gets a Person{} from the DB
func (d *Database) GetPerson(id string) (*Person, error) {
	return d.reader.getPerson(id)
}

// GetPeople gets all Person{} from the DB
func (d *Database) GetPeople() ([]*Person, error) {
	return d.reader.getPeople()
}

// GetSafePerson gets a Person{} from the DB, while being absolutely
// sure there will be no error.
//
// This is only for hardcoded cases where error handling is just pointless.
func (d *Database) GetSafePerson(id string) *Person {
	p, _ := d.GetPerson(id)
	return p
}

// DisablePerson disables or re-enables a person
func (d *Database) DisablePerson(id string) error {
	p, err := d.GetPerson(id)
	if err != nil {
		return err
	}

	p.Disabled = !p.Disabled
	d.SavePerson(p)

	return nil
}

// GetCurrentTournament gets the currently running tournament.
//
// Returns the first matching one, so if there are multiple they will
// be shadowed.
func (d *Database) GetCurrentTournament(s *Server) (*Tournament, error) {
	return d.reader.getCurrentTournament(s)
}

// ClearTestTournaments deletes any tournament that doesn't begin with "DrunkenFall"
func (d *Database) ClearTestTournaments() error {
	err := d.writer.clearTestTournaments(d.Server)

	log.Print("Not sending websocket update; map not implemented")
	// d.Server.SendWebsocketUpdate("all", d.asMap())

	return err
}

// Close closes the database
func (d *Database) Close() error {
	rerr := d.reader.close()
	if rerr != nil {
		return rerr
	}

	return d.writer.close()
}
