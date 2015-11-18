package main

import (
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
