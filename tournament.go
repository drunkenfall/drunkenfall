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

	t.Started = time.Now()
	return nil
}

func main() {
	fmt.Println("...and thus there was light.")
}
