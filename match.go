package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// Match represents a game being played
type Match struct {
	Players    []Player    `json:"players"`
	Judges     []Judge     `json:"judges"`
	Kind       string      `json:"kind"`
	Index      int         `json:"index"`
	Started    time.Time   `json:"started"`
	Ended      time.Time   `json:"ended"`
	Tournament *Tournament `json:"-"`
}

// NewMatch creates a new Match for usage!
func NewMatch(t *Tournament, index int, kind string) *Match {
	m := Match{
		Index:      index,
		Kind:       kind,
		Tournament: t,
	}
	m.Prefill()
	return &m
}

func (m *Match) String() string {
	var tempo string
	var name string

	if !m.IsStarted() {
		tempo = "not started"
	} else if m.IsEnded() {
		tempo = "ended"
	} else {
		tempo = "playing"
	}

	if m.Kind == "final" {
		name = "Final"
	} else {
		name = fmt.Sprintf("%s %d", strings.Title(m.Kind), m.Index+1)
	}

	names := make([]string, 0, len(m.Players))
	for _, p := range m.Players {
		names = append(names, p.Name)
	}

	return fmt.Sprintf(
		"<%s: %s - %s>",
		name,
		strings.Join(names, " / "),
		tempo,
	)
}

// Title returns a title string
func (m *Match) Title() string {
	l := 2
	if m.Kind == "final" {
		return "Final"
	} else if m.Kind == "tryout" {
		l = len(m.Tournament.Tryouts)
	}

	out := fmt.Sprintf(
		"%s %d/%d",
		strings.Title(m.Kind),
		m.Index+1,
		l,
	)
	return out
}

// URL builds the URL to the match
func (m *Match) URL() string {
	out := fmt.Sprintf(
		"/%s/%s/%d",
		m.Tournament.ID,
		m.Kind,
		m.Index,
	)
	return out

}

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(p Player) error {
	if m.ActualPlayers() == 4 {
		return errors.New("cannot add fifth player")
	}

	if len(m.Players) == 4 {
		// Loop through the players and replace the first prefill player that can be found with
		// the actual player.
		for i, o := range m.Players {
			if o.IsPrefill() {
				m.Players[i] = p
				break
			}
		}
	} else {
		m.Players = append(m.Players, p)
	}

	p.Match = m

	return nil
}

// Prefill fills remaining player slots with nil players
func (m *Match) Prefill() error {
	for i := len(m.Players); i < 4; i++ {
		err := m.AddPlayer(Player{})
		if err != nil {
			log.Fatal(err)
		}

	}
	return nil
}

// ActualPlayers returns the number of actual players set in the match
func (m *Match) ActualPlayers() int {
	ret := 0
	for _, p := range m.Players {
		if !p.IsPrefill() {
			ret++
		}
	}

	return ret
}

// Start starts the match
func (m *Match) Start() error {
	if !m.Started.IsZero() {
		return errors.New("match already started")
	}

	// If there are not four players in the match, we need to populate
	// the match with runnerups from the tournament
	if m.ActualPlayers() != 4 {
		m.Tournament.PopulateRunnerups(m)
	}

	for i := range m.Players {
		m.Players[i].Reset()
		m.Players[i].Match = m
	}

	m.Started = time.Now()
	return nil
}

// End signals that the match has ended
//
// It is also the place that moves players into either the Runnerup bracket
// or into their place in the semis.
func (m *Match) End() error {
	if !m.Ended.IsZero() {
		return errors.New("match already ended")
	}

	if m.Kind == "final" {
		m.Tournament.AwardMedals(m)
	} else {
		m.Tournament.MovePlayers(m)
	}
	m.Ended = time.Now()
	return nil
}

// IsStarted returns boolean whether the match has started or not
func (m *Match) IsStarted() bool {
	return !m.Started.IsZero()
}

// IsEnded returns boolean whether the match has ended or not
func (m *Match) IsEnded() bool {
	return !m.Ended.IsZero()
}

// CanEnd returns boolean whether the match can be ended or not
func (m *Match) CanEnd() bool {
	for _, p := range m.Players {
		if p.Kills >= m.Length() {
			return true
		}
	}
	return false
}

// IsOpen returns boolean the match can be controlled or not
func (m *Match) IsOpen() bool {
	return !m.IsEnded()
}
