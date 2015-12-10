package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Tournament is the main container of data for this app.
type Tournament struct {
	Name        string    `json:"name"`
	ID          string    `json:"id"`
	Players     []Player  `json:"players"`
	Winners     []Player  `json:"winners"` // TODO: Refactor to pointer
	Runnerups   []*Player `json:"runnerups"`
	Judges      []Judge   `json:"judges"`
	Tryouts     []*Match  `json:"tryouts"`
	Semis       []*Match  `json:"semis"`
	Final       *Match    `json:"final"`
	Opened      time.Time `json:"opened"`
	Started     time.Time `json:"started"`
	Ended       time.Time `json:"ended"`
	playerRef   map[string]*Player
	db          *Database
	length      int
	finalLength int
}

// NewTournament returns a completely new Tournament
func NewTournament(name, id string, db *Database) (*Tournament, error) {
	t := Tournament{
		Name:   name,
		ID:     id,
		Opened: time.Now(),
		db:     db,
	}
	t.playerRef = make(map[string]*Player)
	t.Persist()
	return &t, nil
}

// LoadTournament loads a tournament from persisted JSON data
func LoadTournament(data []byte, db *Database) (t *Tournament, e error) {
	t = &Tournament{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return t, err
	}

	t.db = db
	t.playerRef = make(map[string]*Player)
	return
}

// Persist tells the database to save this tournament to disk
func (t *Tournament) Persist() error {
	if t.db == nil {
		// This might happen in tests.
		return errors.New("no database instantiated")
	}
	return t.db.Persist(t)
}

// JSON returns a JSON representation of the Tournament
func (t *Tournament) JSON() (out []byte, err error) {
	out, err = json.Marshal(t)
	return
}

// AddPlayer adds a player into the tournament
//
// When adding new players, this means:
//   Generating tryouts
//   Shuffling players into positions
func (t *Tournament) AddPlayer(name string) error {
	p := Player{Name: name}
	if !t.CanJoin(name) {
		return errors.New("player already in match")
	}

	t.Players = append(t.Players, p)
	t.playerRef[name] = &p

	ts := len(t.Tryouts)
	ps := len(t.Players)
	if ts == 0 {
		// No matches yet - add four
		for i := 0; i < 4; i++ {
			match := NewMatch(t, i, "tryout")
			t.Tryouts = append(t.Tryouts, match)
		}
	} else if ts == 4 && ps == 17 {
		// More than 16 players - add four more matches
		for i := 0; i < 4; i++ {
			match := NewMatch(t, i+4, "tryout")
			t.Tryouts = append(t.Tryouts, match)
		}
	}

	// Right now there are only cases where we have two matches in the semis.
	if len(t.Semis) == 0 {
		t.Semis = []*Match{NewMatch(t, 0, "semi"), NewMatch(t, 1, "semi")}
		t.Semis[0].Prefill()
		t.Semis[1].Prefill()
	}
	if t.Final == nil {
		t.Final = NewMatch(t, 0, "final")
		t.Final.Prefill()
	}
	t.ShufflePlayers()
	t.Persist() // TODO: Error handling

	return nil
}

// ShufflePlayers will reposition players into matches
func (t *Tournament) ShufflePlayers() {
	// Reset the set matches
	for _, t := range t.Tryouts {
		t.Players = []Player{}
	}

	// Shuffle all the players
	slice := t.Players
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	// Loop the players and set them into the matches
	for i, p := range slice {
		m := t.Tryouts[i/4]
		m.AddPlayer(p)
	}

	for _, m := range t.Tryouts {
		m.Prefill()
	}
}

// StartTournament will generate the tournament.
//
// This includes:
//  Generating Tryout matches
//  Setting Started date
//
// It will fail if there are not between 8 and 24 players.
func (t *Tournament) StartTournament() error {
	ps := len(t.Players)
	if ps < 8 {
		return fmt.Errorf("Tournament needs at least 8 players, got %d", ps)
	}
	if ps > 24 {
		return fmt.Errorf("Tournament can only host 24 players, got %d", ps)
	}

	t.Started = time.Now()
	return nil
}

// PopulateRunnerups fills a match with the runnerups with best scores
func (t *Tournament) PopulateRunnerups(m *Match) error {
	c := 4 - m.ActualPlayers()
	r, err := t.GetRunnerups(c)
	if err != nil {
		return err
	}

	for _, p := range r {
		m.AddPlayer(p)
	}

	if m.ActualPlayers() != 4 {
		return errors.New("not enough runnerups to populate match")
	}

	return nil
}

// GetRunnerups gets the runnerups for this tournament
//
// The returned list is sorted descending by score.
func (t *Tournament) GetRunnerups(amount int) (ps []Player, err error) {
	err = t.UpdatePlayers()
	if err != nil {
		return
	}

	p := make([]Player, 0, len(t.Runnerups))
	for _, r := range t.Runnerups {
		p = append(p, *r)
	}
	bs := ByRunnerup(p)
	for i := 0; i < amount; i++ {
		// Add the runnerup to the return list
		runnerup := bs[i]
		ps = append(ps, runnerup)
	}

	if len(ps) != amount {
		return ps, errors.New("not enough players added")
	}
	return
}

// UpdatePlayers updates all the player objects with their scores from
// all the matches they have participated in.
func (t *Tournament) UpdatePlayers() error {
	var tp *Player

	// Make sure all players have their score reset to nothing
	for _, p := range t.Players {
		tp = t.playerRef[p.Name]
		tp.Reset()
	}

	for _, m := range t.Tryouts {
		for _, p := range m.Players {
			tp = t.playerRef[p.Name]
			tp.Update(&p)
		}
	}

	for _, m := range t.Semis {
		for _, p := range m.Players {
			tp = t.playerRef[p.Name]
			tp.Update(&p)
		}
	}

	for _, p := range t.Final.Players {
		tp = t.playerRef[p.Name]
		tp.Update(&p)
	}

	return nil
}

// MovePlayers moves the winner(s) of a Match into the next bracket of matches
// or into the Runnerup bracket.
func (t *Tournament) MovePlayers(m *Match) error {
	if m.Kind == "tryout" {
		ps := SortByKills(m.Players)
		for i := 0; i < len(ps); i++ {
			p := ps[i]
			// If we are in a four-match tryout, both the winner and the second-place
			// are to be sent to the semis
			if len(t.Tryouts) == 4 && i < 2 || i == 0 {
				// This spreads the winners into the semis so that the winners do not
				// face off immediately in the semis
				index := (i + m.Index) % 2
				t.Semis[index].AddPlayer(p)

				// If the player is also inside of the runnerups, move them from the
				// runnerup roster since they now have advanced to the finals. This
				// only happens for players that win the runnerup rounds.
				for j := 0; j < len(t.Runnerups); j++ {
					r := t.Runnerups[j]
					if r.Name == p.Name {
						t.Runnerups = append(t.Runnerups[:j], t.Runnerups[j+1:]...)
						break
					}
				}

			} else {
				// For everyone else, add them into the Runnerup bracket unless they are
				// already in there.
				found := false
				for j := 0; j < len(t.Runnerups); j++ {
					r := t.Runnerups[j]
					if r.Name == p.Name {
						found = true
						break
					}
				}
				if !found {
					t.Runnerups = append(t.Runnerups, &p)
				}
			}
		}
	}

	if m.Kind == "semi" {
		// For the semis, just place the winner and silver into the final
		for i, p := range SortByKills(m.Players) {
			if i < 2 {
				t.Final.AddPlayer(p)
			}
		}
	}

	return nil
}

// NextMatch returns the next match
func (t *Tournament) NextMatch() (m *Match, err error) {
	// Firstly, check the tryouts
	for x := range t.Tryouts {
		m = t.Tryouts[x]
		if !m.IsEnded() {
			return
		}
	}

	// If we don't have any tryouts, or there are no tryouts left,
	// check the semis
	for x := range t.Semis {
		m = t.Semis[x]
		if !m.IsEnded() {
			return
		}
	}

	if !t.Final.IsEnded() {
		return t.Final, nil
	}

	return m, errors.New("all matches have been played")
}

// AwardMedals places the winning players in the Winners position
func (t *Tournament) AwardMedals(m *Match) error {
	if m.Kind != "final" {
		return errors.New("awarding medals outside of the final")
	}

	ps := SortByKills(m.Players)
	t.Winners = ps[0:3]

	return nil
}

// IsOpen returns boolean true if the tournament is open for registration
func (t *Tournament) IsOpen() bool {
	return t.Started.IsZero()
}

// CanJoin checks if a player is allowed to join or is already in the tournament
func (t *Tournament) CanJoin(name string) bool {
	for _, p := range t.Players {
		if p.Name == name {
			return false
		}
	}
	return true
}
