package main

import (
	"time"
)

// Match represents a game being played
type Match struct {
	Players [4]Player
	Judges  []Judge
	Started time.Time
	Ended   time.Time
}
