package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Tournament is the main container of data for this app.
type Tournament struct {
	Players     []Player
	Judges      []Judge
	Tryouts     []Match
	Semis       []Match
	Final       Match
	Opened      time.Time
	Started     time.Time
	Ended       time.Time
	length      int
	finalLength int
}

// NewTournament returns a completely new Tournament
func NewTournament() (*Tournament, error) {
	t := Tournament{
		Opened: time.Now(),
	}
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
		m := NewMatch()
		t.Tryouts = append(t.Tryouts, m)
	}

	// Right now there are only cases where we have two matches in the semis.
	t.Semis = []Match{NewMatch(), NewMatch()}
	t.Final = NewMatch()

	return nil
}

// PopulateMatches shuffles players into the matches
func (t *Tournament) PopulateMatches() error {
	// Make a copy of the players list and shuffle it
	slice := t.Players
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

func main() {
	fmt.Println("...and thus there was light.")
}
