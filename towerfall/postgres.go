package towerfall

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

type postgresWriter struct {
	DB  *gorm.DB
	log *zap.Logger
}

type postgresReader struct {
	DB  *gorm.DB
	log *zap.Logger
}

func (d postgresWriter) setUp() error {
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

func (p postgresWriter) migrate(r dbReader, s *Server) error {
	if true {
		return nil
	}

	people, err := r.getPeople()
	for _, ps := range people {
		ps.Correct()
		ps.PreferredColor = ps.ColorPreference[0]
		p.DB.NewRecord(ps)
		p.DB.Create(ps)
	}

	ts, err := r.getTournaments(s)
	if err != nil {
		return err
	}

	for x, t := range ts {
		for y, m := range t.Matches {
			for z, ps := range m.Players {
				ts[x].Matches[y].Players[z].PersonID = ps.Person.PersonID
			}

			for _, r := range m.Rounds {
				m.Commits = append(m.Commits, r.asCommit())
			}

			for z, _ := range m.Messages {
				ts[x].Matches[y].Messages[z].serialize()
			}
		}
	}

	for _, t := range ts {
		p.DB.NewRecord(t)
		p.DB.Create(t)
	}
	return nil
}

// SaveTournament stores the current state of the tournaments into the db
func (d postgresWriter) saveTournament(t *Tournament) error {
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
func (d postgresWriter) overwriteTournament(t *Tournament) error {
	return nil
}

// SavePerson stores a person into the DB
func (d postgresWriter) savePerson(p *Person) error {
	d.DB.Save(p)
	return nil
}

// GetPerson gets a Person{} from the DB
func (d postgresReader) getPerson(id string) (*Person, error) {
	p := Person{}
	d.DB.Where("id = ?", id).First(&p)

	return &p, d.DB.Error
}

// GetSafePerson gets a Person{} from the DB, while being absolutely
// sure there will be no error.
//
// This is only for hardcoded cases where error handling is just pointless.
func (d postgresReader) getSafePerson(id string) *Person {
	p, _ := d.getPerson(id)
	return p
}

// DisablePerson disables or re-enables a person
func (d postgresReader) disablePerson(id string) error {
	// p, err := d.getPerson(id)
	// if err != nil {
	// 	return err
	// }

	// p.Disabled = !p.Disabled
	// d.SavePerson(p)

	return nil
}

// LoadPeople loads the people from the database
func (d postgresReader) getPeople() ([]*Person, error) {
	ret := make([]*Person, 0)
	return ret, nil
}

// getTournament gets a tournament by ID
func (d postgresReader) getTournament(id string, s *Server) (*Tournament, error) {
	return &Tournament{}, errors.New("tournament not found")
}

func (d postgresReader) getTournaments(s *Server) ([]*Tournament, error) {
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
func (d postgresReader) getCurrentTournament(s *Server) (*Tournament, error) {
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
func (d postgresWriter) clearTestTournaments(s *Server) error {
	log.Print("Not sending full update; not implemented")
	// d.Server.SendWebsocketUpdate("all", d.asMap())

	return nil
}

// func (d postgresReader) asMap() map[string]*Tournament {
// 	tournamentMutex.Lock()
// 	out := make(map[string]*Tournament)
// 	for _, t := range d.Tournaments {
// 		out[t.ID] = t
// 	}
// 	tournamentMutex.Unlock()
// 	return out
// }

// Close closes the database
func (d postgresReader) close() error {
	return d.DB.Close()
}

// Close closes the database
func (d postgresWriter) close() error {
	return d.DB.Close()
}
