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

	// Generate tryout matches
	err := t.GenerateTryouts()
	if err != nil {
		return err
	}

	t.Started = time.Now()
	return nil
}

// GenerateTryouts will generate the tryout matches
//
// The tournament model as it stands right now can handle matches sets of
// 2, 4 and 6 matches. If the amount of players do not match, add sets of
// two matches with runnerups.
func (t *Tournament) GenerateTryouts() error {
	var end int
	p := len(t.Players)

	switch p {
	case 8:
		end = 2
	case 9, 10, 11, 12, 13, 14, 15, 16:
		end = 4
	default:
		end = 6
	}

	for i := 0; i < end; i++ {
		m := Match{}
		t.Tryouts = append(t.Tryouts, m)
	}
	return nil
}

func main() {
	fmt.Println("...and thus there was light.")
}
