package main

import (
	"fmt"
	"time"
)

// Tournament is the main container of data for this app.
type Tournament struct {
	Players     [24]Player
	Judges      []Judge
	Tryouts     []Match
	Semis       []Match
	Final       Match
	Started     time.Time
	Ended       time.Time
	length      int
	finalLength int
}

func main() {
	fmt.Println("...and thus there was light.")
}
