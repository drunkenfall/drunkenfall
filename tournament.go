package main

import (
	"fmt"
	"time"
)

// Tournamennt is the main container of data for this app.
type Tournament struct {
	Players     []Player
	Judges      []Judge
	Tryouts     []Match
	Semis       []Match
	Final       Match
	Started     time.Time
	Ended       time.Time
	length      int
	finalLength int
}

// NewTournament returns a completely new Tournament
func NewTournament() (*Tournament, error) {
	t := Tournament{}
	return &t, nil
}

// StartTournament will generate the tournament.
//
// This includes:
//  Generating Tryout matches
//  Setting Started date
//
// It will fail if there are not between 8 and 24 players.
func (t *Tournament) StartTournament() error {
	if len(t.Players) < 8 {
		return fmt.Errorf("Tournament needs at least 8 players, got %s", t.Players)
	}
	if len(t.Players) > 24 {
		return fmt.Errorf("Tournament can only host 24 players, got %s", t.Players)
	}

	// Generate tryouts and semis
	err := t.GenerateMatches()
	if err != nil {
		return err
	}

	t.Started = time.Now()
	return nil
}

// GenerateMatches will generate all the matches for the semis
//
// The tournament model as it stands right now can handle matches sets of
// 2, 4 and 6 matches. If the amount of players do not match, add sets of
// two matches with runnerups.
func (t *Tournament) GenerateMatches() error {
	var tryouts int

	switch len(t.Players) {
	case 8:
		tryouts = 0
	case 9, 10, 11, 12, 13, 14, 15, 16:
		tryouts = 4
	default:
		tryouts = 6
	}

	for i := 0; i < tryouts; i++ {
		m := Match{}
		t.Tryouts = append(t.Tryouts, m)
	}

	// Right now there are only cases where we have two matches in the semis.
	t.Semis = []Match{Match{}, Match{}}

	return nil
}

func main() {
	fmt.Println("...and thus there was light.")
}
