package main

// Participant someone having a role in the tournament
type Participant interface {
	ID() string
	Name() string
	// Matches returns three integers - matches participated, matches won, matches judged.
	Matches() [3]int
}

// Player is a Participant that is actively participating in battles.
type Player struct {
	name           string
	preferredColor string
}

// Color returns the color that the player prefers.
func (p *Player) Color() string {
	return p.preferredColor
}

// Judge is a Participant that has access to the judge functions
type Judge struct {
}
