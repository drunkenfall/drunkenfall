package main

import (
	"sort"
)

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

	shots      int
	sweeps     int
	kills      int
	self       int
	explosions int
}

// NewPlayer returns a new instance of a player
func NewPlayer() Player {
	p := Player{}
	return p
}

// Score calculates the score to determine runnerup positions.
func (p *Player) Score() (out int) {
	// This algorithm is probably flawed, but at least it should be able to
	// determine who is the most entertaining.
	// When executed, a sweep is basically 11 points since scoring a sweep
	// also comes with a shot and three kills.

	out += p.sweeps * 5
	out += p.shots * 3
	out += p.kills * 2
	out += p.self
	out += p.explosions

	return
}

// Color returns the color that the player prefers.
func (p *Player) Color() string {
	return p.preferredColor
}

// AddShot increases the shot count
func (p *Player) AddShot() {
	p.shots++
}

// RemoveShot decreases the shot count
// Fails silently if shots are zero.
func (p *Player) RemoveShot() {
	if p.shots == 0 {
		return
	}
	p.shots--
}

// AddSweep increases the sweep count, gives three kills and a shot.
func (p *Player) AddSweep() {
	p.sweeps++
	p.AddShot()
	p.AddKill()
	p.AddKill()
	p.AddKill()
}

// RemoveSweep decreases the sweep count, three kills and a shot
// Fails silently if sweeps are zero.
func (p *Player) RemoveSweep() {
	if p.sweeps == 0 {
		return
	}
	p.sweeps--
	p.RemoveShot()
	p.RemoveKill()
	p.RemoveKill()
	p.RemoveKill()
}

// AddKill increases the kill count
func (p *Player) AddKill(kills ...int) {
	// This is basically only to help out with testing.
	// Adding an optional argument with the amount of kills lets us just use
	// one call to AddKill() rather than 10.
	if len(kills) > 0 {
		p.kills += kills[0]
	} else {
		p.kills++
	}
}

// RemoveKill decreases the kill count
// Fails silently if kills are zero.
func (p *Player) RemoveKill() {
	if p.kills == 0 {
		return
	}
	p.kills--
}

// AddSelf increases the self count, decreases the kill, and gives a shot
func (p *Player) AddSelf() {
	p.self++
	p.RemoveKill()
	p.AddShot()
}

// RemoveSelf decreases the self count and a shot
// Fails silently if selfs are zero.
func (p *Player) RemoveSelf() {
	if p.self == 0 {
		return
	}
	p.self--
	p.AddKill()
	p.RemoveShot()
}

// AddExplosion increases the explosion count, the kill count and gives a shot
func (p *Player) AddExplosion() {
	p.explosions++
	p.AddShot()
	p.AddKill()
}

// RemoveExplosion decreases the explosion count, a shot and a kill
// Fails silently if explosions are zero.
func (p *Player) RemoveExplosion() {
	if p.explosions == 0 {
		return
	}
	p.explosions--
	p.RemoveShot()
	p.RemoveKill()
}

// Reset resets the stats on a Player to 0
//
// It is to be run in Match.Start()
func (p *Player) Reset() {
	p.shots = 0
	p.sweeps = 0
	p.kills = 0
	p.self = 0
	p.explosions = 0
}

// ByScore is a sort.Interface that sorts players by their score
type ByScore []Player

func (s ByScore) Len() int {
	return len(s)

}
func (s ByScore) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByScore) Less(i, j int) bool {
	// Technically not Less, but we want biggest first...
	return s[i].Score() > s[j].Score()
}

// SortByScore returns a list in order of the score the players have
func SortByScore(ps []Player) []Player {
	sort.Sort(ByScore(ps))
	return ps
}

// ByKills is a sort.Interface that sorts players by their kills
type ByKills []Player

func (s ByKills) Len() int {
	return len(s)

}
func (s ByKills) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByKills) Less(i, j int) bool {
	// Technically not Less, but we want biggest first...
	return s[i].kills > s[j].kills
}

// SortByKills returns a list in order of the kills the players have
func SortByKills(ps []Player) []Player {
	sort.Sort(ByKills(ps))
	return ps
}

// Judge is a Participant that has access to the judge functions
type Judge struct {
}
