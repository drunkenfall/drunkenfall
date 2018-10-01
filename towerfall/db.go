package towerfall

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type Database struct {
	Server *Server
	Reader dbReader
	Writer dbWriter
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

	setUp() error
	migrate(r dbReader, s *Server) error
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

// NewDatabase sets up the database reader and writer
func NewDatabase(c *Config) (*Database, error) {
	var b *bolt.DB
	var pg *gorm.DB
	var r dbReader
	var w dbWriter
	var err error

	log, _ := zap.NewDevelopment()

	if c.DbReader == "bolt" || c.DbWriter == "bolt" {
		b, err = bolt.Open(c.DbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Fatal("bolt open error", zap.Error(err))
		}
	}

	if c.DbReader == "postgres" || c.DbWriter == "postgres" {
		connStr := c.DbPostgresConn
		pg, err = gorm.Open("postgres", connStr)
		if err != nil {
			log.Fatal("postgres open error", zap.Error(err))
		}
	}

	if c.DbReader == "postgres" {
		l := log.With(zap.String("db", "pg_read"))
		r = postgresReader{pg, l}
	} else {
		r = boltReader{b}
	}

	if c.DbWriter == "postgres" {
		l := log.With(zap.String("db", "pg_write"))
		w = postgresWriter{pg, l}
		w.setUp()
	} else {
		w = boltWriter{b}
	}

	db := &Database{Reader: r, Writer: w}

	return db, nil
}

// Migrate migrates the database
func (d *Database) Migrate() error {
	return d.Writer.migrate(d.Reader, d.Server)
}

// GetTournament gets a single tournament
func (d *Database) GetTournament(id string, s *Server) (*Tournament, error) {
	return d.Reader.getTournament(id, s)
}

// LoadTournaments loads the tournaments from the database and into memory
func (d *Database) GetTournaments(s *Server) ([]*Tournament, error) {
	return d.Reader.getTournaments(s)
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) SaveTournament(t *Tournament) error {
	err := d.Writer.saveTournament(t)
	go d.Server.SendWebsocketUpdate("tournament", t)
	return err
}

// OverwriteTournament takes a new foreign Tournament{} object and replaces
// the one with the same ID with that one.
//
// Used from the EditHandler()
func (d *Database) OverwriteTournament(t *Tournament) error {
	return d.Writer.overwriteTournament(t)
}

// SavePerson stores a person into the DB
func (d *Database) SavePerson(p *Person) error {
	return d.Writer.savePerson(p)
}

// GetPerson gets a Person{} from the DB
func (d *Database) GetPerson(id string) (*Person, error) {
	return d.Reader.getPerson(id)
}

// GetPeople gets all Person{} from the DB
func (d *Database) GetPeople() ([]*Person, error) {
	return d.Reader.getPeople()
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
	return d.Reader.getCurrentTournament(s)
}

// ClearTestTournaments deletes any tournament that doesn't begin with "DrunkenFall"
func (d *Database) ClearTestTournaments() error {
	err := d.Writer.clearTestTournaments(d.Server)

	log.Print("Not sending websocket update; map not implemented")
	// d.Server.SendWebsocketUpdate("all", d.asMap())

	return err
}

// Close closes the database
func (d *Database) Close() error {
	rerr := d.Reader.close()
	if rerr != nil {
		return rerr
	}

	return d.Writer.close()
}
