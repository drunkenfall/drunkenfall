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
	Tryouts     []Match   `json:"tryouts"`
	Semis       []Match   `json:"semis"`
	Final       Match     `json:"final"`
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
func LoadTournament(data []byte) (t *Tournament, e error) {
	t = &Tournament{}
	err := json.Unmarshal(data, t)
	if err != nil {
		return t, err
	}

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
func (t *Tournament) AddPlayer(name string) error {
	p := Player{Name: name}
	t.Players = append(t.Players, p)
	t.playerRef[name] = &p
	return nil
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

	// Generate tryouts and semis
	err := t.GenerateMatches()
	if err != nil {
		return err
	}

	// Populate the matches with players
	perr := t.PopulateMatches()
	if perr != nil {
		return perr
	}

	t.Started = time.Now()
	return nil
}

// GenerateMatches will generate all the matches for the semis
//
// The tournament model as it stands right now can handle matches sets of
// 2, 4 and 7 matches. If the amount of players do not match, add sets of
// two matches with runnerups.
func (t *Tournament) GenerateMatches() error {
	var tryouts int
	var kind string

	ps := len(t.Players)

	switch ps {
	case 8:
		tryouts = 0
	case 9, 10, 11, 12, 13, 14, 15, 16:
		tryouts = 4
	default:
		tryouts = 8
	}

	for i := 0; i < tryouts; i++ {
		// Matches that contain only new players are tryouts
		// Any match that contains replayers are runnerups
		if (i+1)*4 <= ps {
			kind = "tryout"
		} else {
			kind = "runnerup"
		}

		m := NewMatch(t, i, kind)
		t.Tryouts = append(t.Tryouts, m)
	}

	// Right now there are only cases where we have two matches in the semis.
	t.Semis = []Match{NewMatch(t, 0, "semi"), NewMatch(t, 1, "semi")}
	t.Final = NewMatch(t, 0, "final")

	return nil
}

// PopulateMatches shuffles players into the matches
func (t *Tournament) PopulateMatches() error {
	// Make a copy of the players map...
	slice := make([]Player, 0, len(t.Players))
	for _, p := range t.Players {
		slice = append(slice, p)
	}

	// ...and shuffle it!
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}

	// The case of eight players is different since they are to go directly
	// to the finals.
	if len(slice) == 8 {
		for x, p := range slice {
			t.Semis[x/4].AddPlayer(p)
		}
		return nil
	}

	// Fill the tryouts as much as possible!
	for x, p := range slice {
		t.Tryouts[x/4].AddPlayer(p)
	}

	return nil
}

// PopulateRunnerups fills a match with the runnerups with best scores
func (t *Tournament) PopulateRunnerups(m *Match) error {
	c := 4 - len(m.Players)
	r, err := t.GetRunnerups(c)
	if err != nil {
		return err
	}

	for _, p := range r {
		m.AddPlayer(p)
	}

	if len(m.Players) != 4 {
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
	if m.Kind == "tryout" || m.Kind == "runnerup" {
		for i, p := range SortByKills(m.Players) {
			// If we are in a four-match tryout, both the winner and the second-place
			// are to be sent to the semis
			if len(t.Tryouts) == 4 && i < 2 || i == 0 {
				// This spreads the winners into the semis so that the winners do not
				// face off immediately in the semis
				index := (i + m.Index) % 2
				t.Semis[index].AddPlayer(p)

				// If the player is also inside of the runnerups, move them from the
				//runnerup roster since they now have advanced to the finals. This
				// only happens for players that win the runnerup rounds.
				for j := 0; j < len(t.Runnerups); j++ {
					r := t.Runnerups[j]
					if r.Name == p.Name {
						t.Runnerups = append(t.Runnerups[:j], t.Runnerups[j+1:]...)
						break
					}
				}

			} else if m.Kind == "tryout" {
				// If we are not already inside a runnerup match, put the non-winners
				// in the runnerup bracket.
				t.Runnerups = append(t.Runnerups, &p)
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
		m = &t.Tryouts[x]
		if !m.IsEnded() {
			return
		}
	}

	// If we don't have any tryouts, or there are no tryouts left,
	// check the semis
	for x := range t.Semis {
		m = &t.Semis[x]
		if !m.IsEnded() {
			return
		}
	}

	if !t.Final.IsEnded() {
		return &t.Final, nil
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
