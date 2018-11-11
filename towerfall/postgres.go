package towerfall

import (
	"fmt"
	"log"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Database struct {
	Server       *Server
	DB           *pg.DB
	log          *zap.Logger
	persistcalls int
}

// SaveTournament stores the current state of the tournaments into the db
func (d *Database) SaveTournament(t *Tournament) error {
	d.persistcalls++
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

	// Reset the scores.
	// TODO(thiderman): Replace this
	p.Shots = 0
	p.Sweeps = 0
	p.Kills = 0
	p.Self = 0
	p.MatchScore = 0
	p.TotalScore = 0
	return d.DB.Insert(p)
}

// AddMatch adds a match
func (d *Database) AddMatch(t *Tournament, m *Match) error {
	m.TournamentID = t.ID
	return d.DB.Insert(m)
}

// SaveMatch saves a match
func (d *Database) SaveMatch(m *Match) error {
	err := d.DB.Update(m)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// AddCommit adds a commit on a match
func (d *Database) AddCommit(m *Match, c *Commit) error {
	c.MatchID = m.ID
	return d.DB.Insert(c)
}

// StoreMessage stores a message for a match
func (d *Database) StoreMessage(m *Match, msg *Message) error {
	msg.MatchID = m.ID

	// Spin off as a goroutine, and print error if it fails; don't care
	// what the caller thinks. Without this, this operation becomes
	// crazy slow since we do it so often.
	go func() {
		err := d.DB.Insert(msg)
		if err != nil {
			log.Printf("Error when saving message: %+v", err)
		}
	}()

	return nil
}

// UpdatePlayer updates one player instance
func (d *Database) UpdatePlayer(m *Match, p *Player) error {
	if p.ID == 0 {
		panic(fmt.Sprintf("player id was zero: %+v", p))
	}

	// Set the computed score on every update
	p.TotalScore = p.Score()
	err := d.DB.Update(p)
	if err != nil {
		return err
	}

	return d.UpdatePlayerSummary(m.Tournament, p)
}

// UpdatePlayerSummary updates the total player data for the tourmament
func (d *Database) UpdatePlayerSummary(t *Tournament, p *Player) error {
	query := `UPDATE player_summaries ps
   SET (shots, sweeps, kills, self, matches, total_score, skill_score)
   =
   (SELECT SUM(shots),
           SUM(sweeps),
           SUM(kills),
           SUM(self),
           COUNT(*),
           SUM(total_score),
           (SUM(total_score) / COUNT(*))
      FROM players P
      INNER JOIN matches M ON p.match_id = m.id
      WHERE m.tournament_id = ?
        AND m.started IS NOT NULL
        AND person_id = ?)
    WHERE person_id = ?
      AND tournament_id = ?;`

	_, err := d.DB.Exec(query, t.ID, p.PersonID, p.PersonID, t.ID)
	if err != nil {
		log.Printf("Summary update failed: %+v", err)
	}
	return err
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
	q := d.DB.Model(&p).Where("NOT disabled").OrderExpr("random()")

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

// NextMatch the next playable match of a tournament
func (d *Database) NextMatch(t *Tournament) (*Match, error) {
	m := Match{
		Tournament: t,
	}

	q := t.db.DB.Model(&m).Where("tournament_id = ? AND started IS NULL", t.ID)
	q = q.Order("id").Limit(1)

	err := q.Select()
	if err != nil {
		return nil, err
	}

	ps := []Player{}
	q = t.db.DB.Model(&ps).Where("match_id = ?", m.ID)
	err = q.Select()
	m.Players = ps

	return &m, err
}

// GetRunnerups gets the next four runnerups, excluding those already
// booked to matches
func (d *Database) GetRunnerups(t *Tournament) ([]*PlayerSummary, error) {
	ps := []*Player{}
	subq := d.DB.Model(&ps).Column("person_id").Join("INNER JOIN matches m on m.id = match_id")
	subq = subq.Where("m.ended IS NULL").Where("m.tournament_id = ?", t.ID)

	ret := []*PlayerSummary{}
	q := d.DB.Model(&ret).Where("tournament_id = ?", t.ID)
	q = q.Where("person_id NOT IN (?)", subq)
	q = q.Order("matches ASC", "skill_score DESC").Limit(4)

	err := q.Select()
	return ret, err
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
