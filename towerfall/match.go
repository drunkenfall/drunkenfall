package towerfall

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/deckarep/golang-set"
	"github.com/mitchellh/mapstructure"
)

const (
	playoff = "playoff"
	semi    = "semi"
	final   = "final"

	EnvironmentKill = -1
)

type Message struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Kill reasons
const (
	rArrow = iota
	rExplosion
	rBrambles
	rJumpedOn
	rLava
	rShock
	rSpikeBall
	rFallingObject
	rSquish
	rCurse
	rMiasma
	rEnemy
	rChalice
)

// Arrows
const (
	aNormal = iota
	aBomb
	aSuperBomb
	aLaser
	aBramble
	aDrill
	aBolt
	aToy
	aFeather
	aTrigger
	aPrism
)

// Message types
const (
	inKill       = "kill"
	inRoundStart = "round_start"
	inRoundEnd   = "round_end"
	inMatchStart = "match_start"
	inMatchEnd   = "match_end"
	inPickup     = "arrows_collected"
	inShot       = "arrow_shot"
	inShield     = "shield_state"
	inWings      = "wings_state"
	inOrbLava    = "lava_orb_state"
	// TODO(thiderman): Non-player orbs are not implemented
	inOrbSlow   = "slow_orb_state"
	inOrbDark   = "dark_orb_state"
	inOrbScroll = "scroll_orb_state"
)

type KillMessage struct {
	Player int `json:"player"`
	Killer int `json:"killer"`
	Cause  int `json:"cause"`
}

type ArrowMessage struct {
	Player int    `json:"player"`
	Arrows Arrows `json:"arrows"`
}

type ShieldMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

type WingsMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

type SlowOrbMessage struct {
	State bool `json:"state"`
}

type DarkOrbMessage struct {
	State bool `json:"state"`
}

type ScrollOrbMessage struct {
	State bool `json:"state"`
}

type LavaOrbMessage struct {
	Player int  `json:"player"`
	State  bool `json:"state"`
}

// List of integers where one item is an arrow type as described in
// the arrow types above.
type Arrows []int

type StartRoundMessage struct {
	Arrows []Arrows `json:"arrows"`
}

// Match represents a game being played
//
// Match.ScoreOrder stores the index to the player in the relative position.
// E.g. if player 3 is in the lead, ScoreOrder[0] will be 2 (the index of
// player 3).
//
// Match.Commits is a list of one commit per round and represents the
// changeset of what happened in the match.
type Match struct {
	Players       []Player      `json:"players"`
	Casters       []*Person     `json:"casters"`
	Kind          string        `json:"kind"`
	Index         int           `json:"index"`
	Length        int           `json:"length"`
	Pause         time.Duration `json:"pause"`
	Scheduled     time.Time     `json:"scheduled"`
	Started       time.Time     `json:"started"`
	Ended         time.Time     `json:"ended"`
	Events        []*Event      `json:"events"`
	Tournament    *Tournament   `json:"-"`
	KillOrder     []int         `json:"kill_order"`
	Rounds        []Round       `json:"commits"`
	Messages      []Message     `json:"messages"`
	Level         string        `json:"level"`
	currentRound  Round
	presentColors mapset.Set
	tournament    *Tournament
}

// Round is a state commit for a round of a match
type Round struct {
	Kills     [][]int `json:"kills"`
	Shots     []bool  `json:"shots"`
	Committed string  `json:"committed"` // ISO-8601
}

// NewMatch creates a new Match
func NewMatch(t *Tournament, kind string) *Match {
	index := len(t.Matches)
	m := Match{
		Index:      index,
		Kind:       kind,
		Tournament: t,
		Length:     t.length,
		Pause:      time.Minute * 5,
		Rounds:     make([]Round, 0),
		currentRound: Round{
			Kills: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
			Shots: []bool{false, false, false, false},
		},
	}
	m.presentColors = mapset.NewSet()

	// Finals are longer <3
	if kind == final {
		m.Length = t.finalLength
	}

	m.Level = m.getRandomLevel()

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

	if m.Kind == final {
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
	if m.Kind == final {
		return "Final"
	} else if m.Kind == playoff {
		l = len(m.Tournament.Matches) - 3
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
		"/%s/%d/",
		m.Tournament.ID,
		m.Index,
	)
	return out
}

// LogEvent makes an event and stores it on the tournament object
func (m *Match) LogEvent(kind, message string, items ...interface{}) {
	ev, err := NewEvent(kind, message, items...)
	if err != nil {
		log.Fatal(err)
	}

	m.Events = append(m.Events, ev)
}

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(p Player) error {
	if len(m.Players) == 4 {
		return errors.New("cannot add fifth player")
	}

	// Reset all possible scores
	p.Reset()

	// Add all the previous players' colors.
	// This is to fix a bug with the presentColors map if the app has been
	// restarted. They cannot be added multiple tines anyway.
	for _, p := range m.Players {
		m.presentColors.Add(p.Color)
	}

	p.Color = p.PreferredColor
	m.presentColors.Add(p.Color)

	// Also set the match pointer
	p.Match = m

	m.Players = append(m.Players, p)

	// If we're adding the fourth player, it's time to correct the conflicts
	if len(m.Players) == 4 && len(m.presentColors.ToSlice()) != 4 {
		if err := m.CorrectFuckingColorConflicts(); err != nil {
			return err
		}
	}

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

// CorrectFuckingColorConflicts corrects color conflicts :@
func (m *Match) CorrectFuckingColorConflicts() error {
	// Make a map of conflicting players keyed on the color
	pairs := make(map[string][]Player)
	for _, color := range m.presentColors.ToSlice() {
		c := color.(string)
		for _, p := range m.Players {
			if p.PreferredColor == c {
				pairs[c] = append(pairs[c], p)
			}
		}
	}

	// Loop over the colors and
	for _, pair := range pairs {
		// If there are two or more players in the group, there is a conflict and
		// they need to be corrected.
		if len(pair) >= 2 {
			// We want to sort them by score, so that we can let the player with the
			// highest score keep their color.
			ps, err := SortByColorConflicts(pair)
			if err != nil {
				return err
			}

			for _, p := range ps[1:] {
				// For the players with lower scores, set their new colors
				new := RandomColor(AvailableColors(m))
				m.presentColors.Add(new)
				p.Color = new

				// Since we are using the tournament level Player object, the compound
				// scores from all other matches are currently on it. Reset that.
				p.Reset()

				if err := m.UpdatePlayer(p); err != nil {
					return err
				}
				m.LogEvent(
					"color_conflict",
					"{nick} corrected from {preferred} to {new}", // Unfortunately we cannot reuse person from below..
					"nick", p.Person.Nick,
					"preferred", p.PreferredColor,
					"new", new,
					"person", p.Person)
			}
		}
	}
	return nil
}

// Commit applies the round actions to the state of the players
// TODO(thiderman): It should not be possible to commit to a non-started match
func (m *Match) Commit(round Round) {
	if round.IsShotUpdate() {
		// The only thing submitted was shots, just update the players
		for i, s := range round.Shots {
			if s {
				m.Players[i].AddShot()
			}
		}
	} else {
		// Apply normally
		for i, score := range round.Kills {
			kills := score[0]
			self := score[1]

			m.Players[i].AddKills(kills)
			if self == -1 {
				m.Players[i].AddSelf()
			}
			if self == -1 || kills == 3 || round.Shots[i] {
				m.Players[i].AddShot()
			}
		}
		m.Rounds = append(m.Rounds, round)
		m.KillOrder = m.MakeKillOrder()
	}

	_ = m.Tournament.Persist()
}

// storeMessage stores a message on the match
func (m *Match) storeMessage(msg Message) error {
	m.Messages = append(m.Messages, msg)
	// TODO(thiderman): Persist here
	return nil
}

// handleMessage decides what to do with an incoming message
func (m *Match) handleMessage(msg Message) error {
	// Store the message. Do this before figuring out the type and even
	// if it would not be parsed.
	err := m.storeMessage(msg)
	if err != nil {
		return nil
	}

	switch msg.Type {
	case inKill:
		km := KillMessage{}
		err := mapstructure.Decode(msg.Data, &km)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}

		return m.Kill(km)

	case inRoundStart:
		sr := StartRoundMessage{}
		err := mapstructure.Decode(msg.Data, &sr)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return m.StartRound(sr)

	case inRoundEnd:
		return m.EndRound()

	case inMatchStart:
		nm, err := m.Tournament.NextMatch()
		if err != nil {
			return err
		}

		return nm.Start(nil)

	case inMatchEnd:
		return m.End(nil)

	case inShot, inPickup:
		am := ArrowMessage{}
		err := mapstructure.Decode(msg.Data, &am)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return m.ArrowUpdate(am)

	case inShield:
		sm := ShieldMessage{}
		err := mapstructure.Decode(msg.Data, &sm)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return m.ShieldUpdate(sm)

	case inWings:
		wm := WingsMessage{}
		err := mapstructure.Decode(msg.Data, &wm)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return m.WingsUpdate(wm)

	case inOrbLava:
		lm := LavaOrbMessage{}
		err := mapstructure.Decode(msg.Data, &lm)
		if err != nil {
			fmt.Println("Error: Could not decode mapstructure", err.Error())
		}
		return m.LavaOrb(lm)

	default:
		log.Printf("Warning: Unknown message type '%s'", msg.Type)
	}

	return nil
}

// sendPlayerUpdate sends a status update for a single player
func (m *Match) sendPlayerUpdate(idx int) error {
	return m.Tournament.server.SendWebsocketUpdate(
		"player",
		PlayerStateUpdateMessage{
			m.Tournament.ID,
			m.Index,
			idx,
			m.Players[idx].State,
		},
	)
}

// EndRound is similar to Commit, but does not alter the score other
// than to manage shots
func (m *Match) EndRound() error {
	for i, score := range m.currentRound.Kills {
		kills := score[0]
		self := score[1]

		if kills == 3 {
			m.Players[i].AddSweep()
		}

		if self == -1 || kills == 3 || m.currentRound.Shots[i] {
			m.Players[i].AddShot()
		}
	}

	m.currentRound.Committed = time.Now().UTC().Format(time.RFC3339)
	m.Rounds = append(m.Rounds, m.currentRound)
	m.KillOrder = m.MakeKillOrder()

	// Reset the Round object
	m.currentRound = Round{
		Kills: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
		Shots: []bool{false, false, false, false},
	}

	return m.Tournament.Persist()
}

// StartRound sets the initial state of player arrows.
func (m *Match) StartRound(sr StartRoundMessage) error {
	for i, as := range sr.Arrows {
		m.Players[i].State.Arrows = as
		m.Players[i].State.Alive = true
		m.Players[i].State.Hat = true
		m.Players[i].State.Lava = false
		m.Players[i].State.Killer = -2
	}
	return m.Tournament.Persist()
}

// ArrowUpdate updates the arrow state for a player
func (m *Match) ArrowUpdate(am ArrowMessage) error {
	m.Players[am.Player].State.Arrows = am.Arrows
	return m.sendPlayerUpdate(am.Player)
}

// ShieldUpdate updates the shield state for a player
func (m *Match) ShieldUpdate(sm ShieldMessage) error {
	m.Players[sm.Player].State.Shield = sm.State
	if sm.State {
		m.LogEvent(
			"shield", "{player} gets a shield",
			"player", m.Players[sm.Player].Name(),
			"person", m.Players[sm.Player].Person,
		)
	} else {
		m.LogEvent(
			"shield_off", "{player}'s shield breaks",
			"player", m.Players[sm.Player].Name(),
			"person", m.Players[sm.Player].Person,
		)
	}
	return m.sendPlayerUpdate(sm.Player)
}

// WingsUpdate updates the wings state for a player
func (m *Match) WingsUpdate(wm WingsMessage) error {
	m.Players[wm.Player].State.Wings = wm.State
	if wm.State {
		m.LogEvent(
			"wings", "{player} grows wings",
			"player", m.Players[wm.Player].Name(),
			"person", m.Players[wm.Player].Person,
		)
	} else {
		m.LogEvent(
			"wings_off", "{player} flies no more",
			"player", m.Players[wm.Player].Name(),
			"person", m.Players[wm.Player].Person,
		)
	}
	return m.sendPlayerUpdate(wm.Player)
}

// LavaOrb sets or unsets the lava for a player
func (m *Match) LavaOrb(lm LavaOrbMessage) error {
	m.Players[lm.Player].State.Lava = lm.State

	if lm.State {
		m.LogEvent(
			"lava", "{player} set the map on fire",
			"player", m.Players[lm.Player].Name(),
			"person", m.Players[lm.Player].Person,
		)
	} else {
		m.LogEvent(
			"lava_off", "{player}'s fire sizzles away",
			"player", m.Players[lm.Player].Name(),
			"person", m.Players[lm.Player].Person,
		)
	}

	return m.sendPlayerUpdate(lm.Player)
}

// Kill records a Kill
func (m *Match) Kill(km KillMessage) error {
	m.Players[km.Player].State.Alive = false
	m.Players[km.Player].State.Killer = km.Killer

	if km.Killer == EnvironmentKill {
		m.Players[km.Player].AddSelf()
		m.currentRound.AddSelf(km.Player)

		m.LogEvent(
			"kill_environ", "{player} was killed by the environment via {cause}",
			"player", m.Players[km.Player].Name(),
			"person", m.Players[km.Player].Person,
			"cause", km.Cause,
		)
	} else if km.Killer == km.Player {
		m.Players[km.Player].AddSelf()
		m.currentRound.AddSelf(km.Player)

		m.LogEvent(
			"suicide", "{player} committed suicide via {cause}",
			"player", m.Players[km.Player].Name(),
			"person", m.Players[km.Player].Person,
			"cause", km.Cause,
		)
	} else {
		m.Players[km.Killer].AddKills(1)
		m.currentRound.AddKill(km.Killer)
		m.LogEvent(
			"kill", "{killer} killed {player} with {cause}",
			"killer", m.Players[km.Killer].Name(),
			"player", m.Players[km.Player].Name(),
			"person", m.Players[km.Killer].Person,
			"cause", km.Cause,
		)
	}
	return m.Tournament.Persist()
}

// Start starts the match
func (m *Match) Start(r *http.Request) error {
	if !m.Started.IsZero() {
		return errors.New("match already started")
	}

	log.Printf("Starting match %d", m.Index)

	for i := range m.Players {
		m.Players[i].Reset()
		m.Players[i].Match = m
	}

	// Set the casters
	m.Casters = m.Tournament.Casters

	// Increment the current match, but only if we're not at the first.
	if m.Index != 0 {
		log.Printf("Increasing current from %d", m.Tournament.Current)
		m.Tournament.Current++
	} else {
		log.Print("Not increasing current when starting first match")
	}

	m.Started = time.Now()
	m.LogEvent(
		"started", "{match} started",
		"match", m.Title(),
		"person", PersonFromSession(m.Tournament.server, r))

	return m.Tournament.Persist()
}

// End signals that the match has ended
//
// It is also the place that moves players into either the Runnerup bracket
// or into their place in the semis.
func (m *Match) End(r *http.Request) error {
	log.Printf("Ending match %d", m.Index)
	if !m.Ended.IsZero() {
		return errors.New("match already ended")
	}

	// XXX(thiderman): In certain test cases a Commit() might not have been run
	// and therefore this might not have been set. Since the calculation is
	// quick and has no side effects, it's easier to just add it here now. In
	// the future, make the tests better.
	m.KillOrder = m.MakeKillOrder()

	// Give the winner one last shot
	winner := m.KillOrder[0]
	m.Players[winner].AddShot()

	m.Ended = time.Now()
	m.LogEvent(
		"ended", "{match} ended",
		"match", m.Title(),
		"person", PersonFromSession(m.Tournament.server, r))

	if m.Kind == final {
		if err := m.Tournament.AwardMedals(r, m); err != nil {
			return err
		}
	} else {
		if err := m.Tournament.MovePlayers(m); err != nil {
			return err
		}
	}

	m.Tournament.Persist()
	return nil
}

// Reset resets all the player scores to zero and removes all the commits
func (m *Match) Reset() error {
	// Reset dem players
	for i := range m.Players {
		m.Players[i].Reset()
	}

	// And remove all the rounds
	m.Rounds = make([]Round, 0)

	m.Tournament.Persist()
	return nil
}

// Autoplay runs through the entire match simulating real play
func (m *Match) Autoplay() {
	if !m.IsStarted() {
		m.Start(nil)
	}
	for !m.CanEnd() {
		m.Commit(NewAutoplayRound())
	}
	m.End(nil)
}

// SetTime sets the scheduled time based on the Pause attribute
func (m *Match) SetTime(r *http.Request, minutes int) {
	m.Scheduled = time.Now().Add(time.Minute * time.Duration(minutes))

	m.LogEvent(
		"time_set", "{match} scheduled in {minutes}m",
		"minutes", minutes,
		"match", m.Title(),
		"person", PersonFromSession(m.Tournament.server, r))
	m.Tournament.Persist()
}

// IsStarted returns boolean whether the match has started or not
func (m *Match) IsStarted() bool {
	return !m.Started.IsZero()
}

// IsEnded returns boolean whether the match has ended or not
func (m *Match) IsEnded() bool {
	return !m.Ended.IsZero()
}

// IsScheduled returns boolean whether the match has been scheduled or not
func (m *Match) IsScheduled() bool {
	return !m.Scheduled.IsZero()
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

// MakeKillOrder returns the score in order of the number of kills in the match.
func (m *Match) MakeKillOrder() (ret []int) {
	ps := SortByKills(m.Players)
	for _, p := range ps {
		for i, o := range m.Players {
			if p.Name() == o.Name() {
				ret = append(ret, i)
				break
			}
		}
	}

	return
}

// ArchersHarmed returns the number of killed archers during the match
func (m *Match) ArchersHarmed() int {
	ret := 0

	for _, r := range m.Rounds {
		for _, k := range r.Kills {
			ret += k[0]

			// If someone suicided, it shows up as a minus one. This means
			// an archer was harmed and should count towards the total.
			if k[1] == -1 {
				ret++
			}
		}
	}

	return ret
}

// Duration returns how long the match took
func (m *Match) Duration() time.Duration {
	return m.Ended.Sub(m.Started)
}

func (m *Match) getRandomLevel() string {
	l := m.Tournament.Levels[m.Kind]
	return l[m.Index%len(l)]
}

// NewMatchCommit makes a new MatchCommit object from a CommitRequest
func NewMatchCommit(c CommitRequest) Round {
	states := c.State
	m := Round{
		[][]int{
			{states[0].Ups, states[0].Downs},
			{states[1].Ups, states[1].Downs},
			{states[2].Ups, states[2].Downs},
			{states[3].Ups, states[3].Downs},
		},
		[]bool{
			states[0].Shot,
			states[1].Shot,
			states[2].Shot,
			states[3].Shot,
		},
		// ISO-8601 timestamp
		time.Now().UTC().Format(time.RFC3339),
	}

	return m
}

// IsShotUpdate returns true if the only thing that happened was shots
func (r *Round) IsShotUpdate() bool {
	for _, y := range r.Kills {
		for _, z := range y {
			if z != 0 {
				return false
			}
		}
	}

	for _, s := range r.Shots {
		if s {
			return true
		}
	}

	return false
}

// NewAutoplayRound fakes player activity in a round
//
// It randomizes how many kills the players get, and it randomizes
// shots every now and again. This does not fully represent actual
// gameplay, since technically all four players could get a sweep in
// the same match. However, since this is for testing purposes it is
// acceptable that such is the case.
func NewAutoplayRound() Round {
	r := Round{
		[][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
		[]bool{false, false, false, false},
		time.Now().UTC().Format(time.RFC3339),
	}

	rand.Seed(time.Now().UnixNano())
	for x := 0; x < 4; x++ {
		y := rand.Intn(100)
		// 5% of the times - sweep
		// 20% of the times - 2 kills
		// 70% of the times - 1 kill
		if y <= 5 {
			r.Kills[x][0] = 3
		} else if y <= 20 {
			r.Kills[x][0] = 2
		} else if y <= 70 {
			r.Kills[x][0] = 1
		}

		// 10% of the time - accidental self
		if rand.Intn(10)%10 == 0 {
			r.Kills[x][1] = -1
		}

		// 10% of the time - shot from the judges
		if rand.Intn(10)%10 == 0 {
			r.Shots[x] = true
		}
	}

	return r
}

// AddKill adds one kill to the specified player
func (r *Round) AddKill(p int) {
	if len(r.Kills) == 0 {
		r.Reset()
	}

	r.Kills[p][0] += 1
}

// AddSelf adds one self to the specified player
func (r *Round) AddSelf(p int) {
	if len(r.Kills) == 0 {
		r.Reset()
	}

	r.Kills[p][1] -= 1
}

func (r *Round) Reset() {
	r.Kills = [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}}
	r.Shots = []bool{false, false, false, false}
}
