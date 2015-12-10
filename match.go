package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Match represents a game being played
type Match struct {
	Players    []Player  `json:"players"`
	Judges     []Judge   `json:"judges"`
	Kind       string    `json:"kind"`
	Index      int       `json:"index"`
	Started    time.Time `json:"started"`
	Ended      time.Time `json:"ended"`
	tournament *Tournament
}

// NewMatch creates a new Match for usage!
func NewMatch(t *Tournament, index int, kind string) *Match {
	m := Match{
		Index:      index,
		Kind:       kind,
		tournament: t,
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

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(p Player) error {
	if len(m.Players) == 4 {
		return errors.New("cannot add fifth player")
	}

	m.Players = append(m.Players, p)
	return nil
}

// Prefill fills remaining player slots with nil players
func (m *Match) Prefill() error {
	for i := len(m.Players); i < 4; i++ {
		m.AddPlayer(Player{})
	}
	return nil
}

// Start starts the match
func (m *Match) Start() error {
	if !m.Started.IsZero() {
		return errors.New("match already started")
	}

	// If there are not four players in the match, we need to populate
	// the match with runnerups from the tournament
	if len(m.Players) != 4 {
		m.tournament.PopulateRunnerups(m)
	}

	for i := range m.Players {
		m.Players[i].Reset()
		m.Players[i].match = m
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
		m.tournament.AwardMedals(m)
	} else {
		m.tournament.MovePlayers(m)
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
