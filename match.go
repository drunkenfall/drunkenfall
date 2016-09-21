package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Match represents a game being played
type Match struct {
	Players       []Player    `json:"players"`
	Judges        []Judge     `json:"judges"`
	Kind          string      `json:"kind"`
	Index         int         `json:"index"`
	Length        int         `json:"length"`
	Started       time.Time   `json:"started"`
	Ended         time.Time   `json:"ended"`
	Tournament    *Tournament `json:"-"`
	presentColors map[string]bool
}

// NewMatch creates a new Match for usage!
func NewMatch(t *Tournament, index int, kind string) *Match {
	m := Match{
		Index:      index,
		Kind:       kind,
		Tournament: t,
		Length:     10,
	}
	m.presentColors = make(map[string]bool)

	// Finals are longer <3
	if kind == "final" {
		m.Length = 20
	}

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
		names = append(names, p.Name())
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
		"/%s/%s/%d/",
		m.Tournament.ID,
		m.Kind,
		m.Index,
	)
	return out
}

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(p Player) error {
	if len(m.Players) == 4 {
		return errors.New("cannot add fifth player")
	}

	// Reset all possible scores
	p.Reset()

	c := p.Person.ColorPreference[0]
	p.OriginalColor = c
	if _, ok := m.presentColors[c]; ok {
		// Color is already present - give the player a new random one.
		c = Colors.Available(m).Random()
	}

	// Set the player color
	p.Color = c
	m.presentColors[c] = true

	// Also set the match pointer
	p.Match = m

	m.Players = append(m.Players, p)

	return nil
}

// UpdatePlayer updates a player for the given match
func (m *Match) UpdatePlayer(p Player) error {
	for i, o := range m.Players {
		if o.Name() == p.Name() {
			m.Players[i] = p
		}
	}
	return nil
}

// Commit adds a state of the players
func (m *Match) Commit(scores [][]int, shots []bool) {
	for i, score := range scores {
		shotGiven := false
		ups := score[0]
		downs := score[1]

		// If the score is 3, then this player killed everyone else.
		// Count that as a sweep.
		if ups == 3 {
			m.Players[i].AddSweep()
			// Since AddSweep gives a shot, we shouldn't give another further down.
			shotGiven = true

			// If we have a sweep and a down, we need to redact the
			// extra shot, because no player should get more than
			// one shot per round. Removing one from here makes
			// sure that the `downs` calculation below still adds
			// the one.
			if downs == 1 {
				m.Players[i].RemoveShot()
			}
		} else if ups > 0 {
			m.Players[i].AddKill(ups)
		}

		if downs != 0 {
			m.Players[i].AddSelf()
		}

		if shots[i] && !shotGiven {
			m.Players[i].AddShot()
		}
	}

	_ = m.Tournament.Persist()
}

// Start starts the match
func (m *Match) Start() error {
	if !m.Started.IsZero() {
		return errors.New("match already started")
	}

	// If there are not four players in the match, we need to populate
	// the match with runnerups from the tournament
	if len(m.Players) != 4 {
		m.Tournament.PopulateRunnerups(m)
	}

	for i := range m.Players {
		m.Players[i].Reset()
		m.Players[i].Match = m
	}

	m.Started = time.Now()
	if m.Tournament != nil {
		m.Tournament.Persist()
	}
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

	// Give the winner one last shot
	ps := ByScore(m.Players)
	winner := ps[0].Name()
	for i, p := range m.Players {
		if p.Name() == winner {
			m.Players[i].AddShot()
			break
		}
	}

	m.Ended = time.Now()
	// TODO: This is for the tests not to break. Fix by setting up better tests.
	if m.Tournament != nil {
		if m.Kind == "final" {
			m.Tournament.AwardMedals(m)
		} else {
			m.Tournament.MovePlayers(m)
		}

		m.Tournament.Persist()
	}
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

// CanStart returns boolean the match can be controlled or not
func (m *Match) CanStart() bool {
	return !m.IsStarted() && !m.IsEnded()
}

// CanEnd returns boolean whether the match can be ended or not
func (m *Match) CanEnd() bool {
	if !m.IsOpen() {
		return false
	}
	for _, p := range m.Players {
		if p.Kills >= m.Length {
			return true
		}
	}
	return false
}

// IsOpen returns boolean the match can be controlled or not
func (m *Match) IsOpen() bool {
	return m.IsStarted() && !m.IsEnded()
}
