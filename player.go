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

// RunnerupScore calculates the score to determine runnerup positions.
func (p *Player) RunnerupScore() (out int) {
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
