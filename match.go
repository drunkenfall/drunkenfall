package main

import (
	"errors"
	"time"
)

// Match represents a game being played
type Match struct {
	Players []Player
	Judges  []Judge
	Started time.Time
	Ended   time.Time
}

// NewMatch creates a new Match for usage!
func NewMatch() Match {
	return Match{}
}

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(p Player) error {
	if len(m.Players) == 4 {
		return errors.New("cannot add fifth player")
	}

	m.Players = append(m.Players, p)
	return nil
}
