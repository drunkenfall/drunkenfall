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

// AddSweep increases the sweep count
func (p *Player) AddSweep() {
	p.sweeps++
}

// RemoveSweep decreases the sweep count
// Fails silently if sweeps are zero.
func (p *Player) RemoveSweep() {
	if p.sweeps == 0 {
		return
	}
	p.sweeps--
}

// AddKill increases the kill count
func (p *Player) AddKill() {
	p.kills++
}

// RemoveKill decreases the kill count
// Fails silently if kills are zero.
func (p *Player) RemoveKill() {
	if p.kills == 0 {
		return
	}
	p.kills--
}

// AddSelf increases the self count
func (p *Player) AddSelf() {
	p.self++
}

// RemoveSelf decreases the self count
// Fails silently if selfs are zero.
func (p *Player) RemoveSelf() {
	if p.self == 0 {
		return
	}
	p.self--
}

// AddExplosion increases the explosion count
func (p *Player) AddExplosion() {
	p.explosions++
}

// RemoveExplosion decreases the explosion count
// Fails silently if explosions are zero.
func (p *Player) RemoveExplosion() {
	if p.explosions == 0 {
		return
	}
	p.explosions--
}

// Judge is a Participant that has access to the judge functions
type Judge struct {
}
