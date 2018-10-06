package towerfall

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

type Database struct {
	Server *Server
	DB     *gorm.DB
	log    *zap.Logger
}

func (d *Database) setUp() error {
	// d.DB.Exec("DROP TABLE tournaments; DROP TABLE matches; DROP TABLE players; DROP TABLE people; DROP TABLE commits; DROP TABLE messages;")
	// d.DB.AutoMigrate(&Tournament{})
	// d.DB.AutoMigrate(&PlayerSummary{})
	// d.DB.AutoMigrate(&Match{})
	// d.DB.AutoMigrate(&Player{})
	// d.DB.AutoMigrate(&Person{})
	// d.DB.AutoMigrate(&Commit{})
	// d.DB.AutoMigrate(&Message{})
	return nil
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) SaveTournament(t *Tournament) error {
	d.DB.Save(t)
	err := d.DB.Error
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// OverwriteTournament takes a new foreign Tournament{} object and replaces
// the one with the same ID with that one.
//
// Used from the EditHandler()
func (d *Database) OverwriteTournament(t *Tournament) error {
	return nil
}

// SavePerson stores a person into the DB
func (d *Database) SavePerson(p *Person) error {
	d.DB.Save(p)
	return nil
}

// GetPerson gets a Person{} from the DB
func (d *Database) GetPerson(id string) (*Person, error) {
	p := Person{}
	d.DB.Where("id = ?", id).First(&p)

	return &p, d.DB.Error
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
	// p, err := d.getPerson(id)
	// if err != nil {
	// 	return err
	// }

	// p.Disabled = !p.Disabled
	// d.SavePerson(p)

	return nil
}

// LoadPeople loads the people from the database
func (d *Database) GetPeople() ([]*Person, error) {
	ret := make([]*Person, 0)
	return ret, nil
}

// getTournament gets a tournament by ID
func (d *Database) GetTournament(id string, s *Server) (*Tournament, error) {
	return &Tournament{}, errors.New("tournament not found")
}

func (d *Database) GetTournaments(s *Server) ([]*Tournament, error) {
	ret := make([]*Tournament, 0)
	d.DB.Find(&ret)
	errs := d.DB.GetErrors()
	if len(errs) != 0 {
		for _, err := range errs {
			d.log.Error("getTournament error", zap.Error(err))
		}
		return ret, errors.New("errors found")
	}
	return ret, nil
}

// GetCurrentTournament gets the currently running tournament.
//
// Returns the first matching one, so if there are multiple they will
// be shadowed.
func (d *Database) GetCurrentTournament(s *Server) (*Tournament, error) {
	ts, err := d.GetTournaments(s)
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
func (d *Database) ClearTestTournaments() error {
	log.Print("Not sending full update; not implemented")
	// d.Server.SendWebsocketUpdate("all", d.asMap())

	return nil
}

// Close closes the database
func (d *Database) Close() error {
	return d.DB.Close()
}
