package towerfall

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type Database struct {
	Server *Server
	DB     *pg.DB
	log    *zap.Logger
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) SaveTournament(t *Tournament) error {
	err := d.DB.Update(t)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) NewTournament(t *Tournament) error {
	return d.DB.Insert(t)
}

// AddPlayer adds a player object to the tournament
func (d *Database) AddPlayer(t *Tournament, ps *PlayerSummary) error {
	ps.TournamentID = t.ID
	return d.DB.Insert(ps)
}

// AddPlayerToMatch adds a player object to a match
func (d *Database) AddPlayerToMatch(m *Match, p *Player) error {
	p.MatchID = m.ID
	return d.DB.Insert(p)
}

// AddMatch adds a match
func (d *Database) AddMatch(t *Tournament, m *Match) error {
	m.TournamentID = t.ID
	return d.DB.Insert(m)
}

// AddCommit adds a commit on a match
func (d *Database) AddCommit(m *Match, c *Commit) error {
	c.MatchID = m.ID
	return d.DB.Insert(c)
}

// StoreMessage stores a message for a match
func (d *Database) StoreMessage(m *Match, msg *Message) error {
	msg.MatchID = m.ID
	return d.DB.Insert(msg)
}

// UpdatePlayer updates one player instance
func (d *Database) UpdatePlayer(p *Player) error {
	if p.ID == 0 {
		panic(fmt.Sprintf("player id was zero: %+v", p))
	}

	// Set the computed score on every update
	p.TotalScore = p.Score()
	return d.DB.Update(p)
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
	d.DB.Update(p)
	return nil
}

// GetPerson gets a Person{} from the DB
func (d *Database) GetPerson(id string) (*Person, error) {
	p := Person{
		PersonID: id,
	}
	err := d.DB.Select(&p)

	return &p, err
}

// GetRandomPerson gets a random Person{} from the DB
func (d *Database) GetRandomPerson(used []string) (*Person, error) {
	p := Person{}
	q := d.DB.Model(&p).OrderExpr("random()")

	if len(used) != 0 {
		args := make([]interface{}, 0)
		for _, u := range used {
			args = append(args, u)
		}
		q = q.WhereIn("person_id NOT IN (?)", args...)
	}

	err := q.First()
	return &p, err
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

// getTournament gets a tournament by slug
func (d *Database) GetTournament(slug string) (*Tournament, error) {
	t := Tournament{}
	err := d.DB.Model(&t).Where("slug = ?", slug).First()
	return &t, err
}

func (d *Database) GetTournaments() ([]*Tournament, error) {
	ret := make([]*Tournament, 0)
	err := d.DB.Model(&ret).Select()
	return ret, err
}

// GetCurrentTournament gets the currently running tournament.
//
// Returns the first matching one, so if there are multiple they will
// be shadowed.
func (d *Database) GetCurrentTournament() (*Tournament, error) {
	ts, err := d.GetTournaments()
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

// GetMatch gets a match
func (d *Database) GetMatch(id uint) (*Match, error) {
	m := &Match{
		ID: id,
	}

	err := d.DB.Select(&m)
	return m, err
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
