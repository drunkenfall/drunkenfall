package main

import (
	"errors"
	"time"
)

// Match represents a game being played
type Match struct {
	Players    []Player
	Judges     []Judge
	Kind       string
	Index      int
	Started    time.Time
	Ended      time.Time
	tournament *Tournament
}

// NewMatch creates a new Match for usage!
func NewMatch(t *Tournament, index int, kind string) Match {
	return Match{
		Index:      index,
		Kind:       kind,
		tournament: t,
	}
}

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(p Player) error {
	if len(m.Players) == 4 {
		return errors.New("cannot add fifth player")
	}

	m.Players = append(m.Players, p)
	return nil
}

// StartMatch starts the match
func (m *Match) StartMatch() error {
	if !m.Started.IsZero() {
		return errors.New("match already started")
	}

	m.Started = time.Now()
	return nil
}

// EndMatch signals that the match has ended
//
// It is also the place that moves players into either the Runnerup bracket
// or into their place in the semis.
func (m *Match) EndMatch() error {
	if !m.Ended.IsZero() {
		return errors.New("match already ended")
	}

	m.tournament.MovePlayers(m)

	m.Ended = time.Now()
	return nil
}

// IsStarted returns boolean whether the match has started or not
func (m *Match) IsStarted() bool {
	return !m.Started.IsZero()
}
