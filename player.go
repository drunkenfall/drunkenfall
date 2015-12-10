package main

import (
	"fmt"
	"sort"
)

// Participant someone having a role in the tournament
type Participant interface {
	ID() string
}

// Player is a Participant that is actively participating in battles.
type Player struct {
	Name           string `json:"name"`
	PreferredColor string `json:"preferred_color"`
	Shots          int    `json:"shots"`
	Sweeps         int    `json:"sweeps"`
	Kills          int    `json:"kills"`
	Self           int    `json:"self"`
	Explosions     int    `json:"explosions"`
	Matches        int    `json:"matches"`
	Match          *Match
}

// NewPlayer returns a new instance of a player
func NewPlayer() Player {
	p := Player{}
	return p
}

func (p *Player) String() string {
	return fmt.Sprintf(
		"%s: %dsh %dsw %dk %ds %de",
		p.Name,
		p.Shots,
		p.Sweeps,
		p.Kills,
		p.Self,
		p.Explosions,
	)
}

// Score calculates the score to determine runnerup positions.
func (p *Player) Score() (out int) {
	// This algorithm is probably flawed, but at least it should be able to
	// determine who is the most entertaining.
	// When executed, a sweep is basically 11 points since scoring a sweep
	// also comes with a shot and three kills.

	out += p.Sweeps * 5
	out += p.Shots * 3
	out += p.Kills * 2
	out += p.Self
	out += p.Explosions

	return
}

// Color returns the color that the player prefers.
func (p *Player) Color() string {
	if p.IsPrefill() {
		return "prefill"
	}
	return p.PreferredColor
}

// Index returns the index in the current match
func (p *Player) Index() int {
	if p.Match != nil {
		for i, o := range p.Match.Players {
			if p.Name == o.Name {
				return i
			}
		}
	}
	return -1
}

// IsPrefill returns whether the player is a prefill placeholder or not
func (p *Player) IsPrefill() bool {
	return p.Name == ""
}

// AddShot increases the shot count
func (p *Player) AddShot() {
	p.Shots++
}

// RemoveShot decreases the shot count
// Fails silently if shots are zero.
func (p *Player) RemoveShot() {
	if p.Shots == 0 {
		return
	}
	p.Shots--
}

// AddSweep increases the sweep count, gives three kills and a shot.
func (p *Player) AddSweep() {
	p.Sweeps++
	p.AddShot()
	p.AddKill()
	p.AddKill()
	p.AddKill()
}

// RemoveSweep decreases the sweep count, three kills and a shot
// Fails silently if sweeps are zero.
func (p *Player) RemoveSweep() {
	if p.Sweeps == 0 {
		return
	}
	p.Sweeps--
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
		p.Kills += kills[0]
	} else {
		p.Kills++
	}
}

// RemoveKill decreases the kill count
// Fails silently if kills are zero.
func (p *Player) RemoveKill() {
	if p.Kills == 0 {
		return
	}
	p.Kills--
}

// AddSelf increases the self count, decreases the kill, and gives a shot
func (p *Player) AddSelf() {
	p.Self++
	p.RemoveKill()
	p.AddShot()
}

// RemoveSelf decreases the self count and a shot
// Fails silently if selfs are zero.
func (p *Player) RemoveSelf() {
	if p.Self == 0 {
		return
	}
	p.Self--
	p.AddKill()
	p.RemoveShot()
}

// AddExplosion increases the explosion count, the kill count and gives a shot
func (p *Player) AddExplosion() {
	p.Explosions++
	p.AddShot()
	p.AddKill()
}

// RemoveExplosion decreases the explosion count, a shot and a kill
// Fails silently if Explosions are zero.
func (p *Player) RemoveExplosion() {
	if p.Explosions == 0 {
		return
	}
	p.Explosions--
	p.RemoveShot()
	p.RemoveKill()
}

// Reset resets the stats on a Player to 0
//
// It is to be run in Match.Start()
func (p *Player) Reset() {
	p.Shots = 0
	p.Sweeps = 0
	p.Kills = 0
	p.Self = 0
	p.Explosions = 0
	p.Matches = 0
}

// Update updates a player with the scores of another
//
// This is primarily used by the tournament score calculator
func (p *Player) Update(other *Player) error {
	p.Shots += other.Shots
	p.Sweeps += other.Sweeps
	p.Kills += other.Kills
	p.Self += other.Self
	p.Explosions += other.Explosions

	// Every call to this method is per match. Count every call
	// as if a match.
	p.Matches++

	return nil
}

// HTML renders the HTML of a player
func (p *Player) HTML() (out string) {
	return
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
	return s[i].Kills > s[j].Kills
}

// SortByKills returns a list in order of the kills the players have
func SortByKills(ps []Player) []Player {
	sort.Sort(ByKills(ps))
	return ps
}

// ByRunnerup is a sort.Interface that sorts players by their runnerup status
//
// This is almost exactly the same as score, but the number of matches a player
// has played factors in, and players that have played less matches are sorted
// favorably.
type ByRunnerup []Player

func (s ByRunnerup) Len() int {
	return len(s)

}
func (s ByRunnerup) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByRunnerup) Less(i, j int) bool {
	if s[i].Matches == s[j].Matches {
		// Same as by kills
		return s[i].Kills > s[j].Kills
	}
	// Lower is better - the ones that have not played should be at the top
	return s[i].Matches < s[j].Matches
}

// SortByRunnerup returns a list in order of the kills the players have
func SortByRunnerup(ps []Player) []Player {
	sort.Sort(ByRunnerup(ps))
	return ps
}

// Judge is a Participant that has access to the judge functions
type Judge struct {
}
