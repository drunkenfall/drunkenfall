package towerfall

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	ErrPublishDisconnected = errors.New("not connected; will not publish")
	ErrAlreadyInTournament = errors.New("already in tournament")
)

// Tournament is the main container of data for this app.
type Tournament struct {
	ID            uint             `json:"dbid"`
	Name          string           `json:"name"`
	Slug          string           `json:"id"`
	Players       []PlayerSummary  `json:"players"`
	Runnerups     []*PlayerSummary `json:"-" sql:"-"`
	Casters       []*Person        `json:"-" sql:"-"`
	Matches       []*Match         `json:"-"`
	Opened        time.Time        `json:"opened"`
	Scheduled     time.Time        `json:"scheduled"`
	Started       time.Time        `json:"started"`
	QualifyingEnd time.Time        `json:"qualifying_end"`
	Ended         time.Time        `json:"ended"`
	Color         string           `json:"color"`
	Cover         string           `json:"cover"`
	Length        int              `json:"length"`
	FinalLength   int              `json:"final_length"`
	connected     bool
	db            *Database
	server        *Server
}

const minPlayers = 12
const matchLength = 10
const finalLength = 20

// NewTournament returns a completely new Tournament
func NewTournament(name, id, cover string, scheduledStart time.Time, c *gin.Context, server *Server) (*Tournament, error) {
	t := Tournament{
		Name:        name,
		Slug:        id,
		Opened:      time.Now(),
		Scheduled:   scheduledStart,
		Cover:       cover,
		Length:      matchLength,
		FinalLength: finalLength,
		db:          server.DB,
		server:      server,
	}

	err := t.db.NewTournament(&t)
	return &t, err
}

// Semi returns one of the two semi matches
func (t *Tournament) Semi(index int) *Match {
	return t.Matches[len(t.Matches)-3+index]
}

// Final returns the final match
func (t *Tournament) Final() *Match {
	return t.Matches[len(t.Matches)-1]
}

// Persist tells the database to save this tournament to disk
func (t *Tournament) Persist() error {
	if t.db == nil {
		// This might happen in tests.
		return errors.New("no database instantiated")
	}

	return t.db.SaveTournament(t)
}

// PublishNext sends information about the next match to the game
//
// It only does this if the match already has four players. If it does
// not, it's a semi that needs backfilling, and then the backfilling
// will make the publish. This should always be called before the
// match is started, so t.NextMatch() can always safely be used.
func (t *Tournament) PublishNext() error {
	if !t.connected {
		if t.server.config.Production {
			t.server.log.Info("Not publishing disconnected tournament")
		}
		return ErrPublishDisconnected
	}

	next, err := t.NextMatch()
	if err != nil {
		return err
	}

	if len(next.Players) != 4 {
		return ErrPublishIncompleteMatch
	}

	msg := GameMatchMessage{
		Tournament: t.Slug,
		Level:      next.realLevel(),
		Length:     next.Length,
		Ruleset:    next.Ruleset,
		Kind:       next.Kind,
	}

	for _, p := range next.Players {
		gp := GamePlayer{
			TopName:    p.DisplayNames[0],
			BottomName: p.DisplayNames[1],
			Color:      p.NumericColor(),
			ArcherType: p.Person.ArcherType,
		}
		msg.Players = append(msg.Players, gp)
	}

	t.server.log.Info("Sending publish", zap.Any("match", msg))
	return t.server.publisher.Publish(gMatch, msg)
}

// connect sets the connect variable
func (t *Tournament) connect(connected bool) {
	if t.connected == connected {
		return
	}

	t.server.log.Info(
		"Tournament connection changed",
		zap.Bool("connected", connected),
	)
	t.connected = connected
}

// JSON returns a JSON representation of the Tournament
func (t *Tournament) JSON() (out []byte, err error) {
	return json.Marshal(t)
}

// URL returns the URL for the tournament
func (t *Tournament) URL() string {
	out := fmt.Sprintf("/%s/", t.Slug)
	return out
}

// AddPlayer adds a player into the tournament
func (t *Tournament) AddPlayer(p *PlayerSummary) error {
	p.Person.Correct()

	if t.IsInTournament(p.Person) {
		return ErrAlreadyInTournament
	}

	t.Players = append(t.Players, *p)

	// If the tournament is already started, just add the player into the
	// runnerups so that they will be placed at the end immediately.
	if !t.Started.IsZero() {
		t.Runnerups = append(t.Runnerups, p)
	}

	return t.db.AddPlayer(t, p)
}

// TogglePlayer toggles a player in a tournament
func (t *Tournament) TogglePlayer(id string) error {
	p, err := t.db.GetPerson(id)
	if err != nil {
		return err
	}

	// If the player is already in the tournament, it is time to remove them
	if t.IsInTournament(p) {
		ps, err := t.db.GetPlayerSummary(t, id)
		if err != nil {
			return err
		}
		return t.db.RemovePlayer(ps)
	}

	ps := NewPlayerSummary(p)
	return t.AddPlayer(ps)
}

// SetCasters resets the casters to the input people
func (t *Tournament) SetCasters(ids []string) error {
	t.Casters = make([]*Person, 0)
	for _, id := range ids {
		ps, _ := t.db.GetPerson(id)
		t.Casters = append(t.Casters, ps)
	}

	return t.Persist()
}

// StartTournament will generate the tournament.
func (t *Tournament) StartTournament(c *gin.Context) error {
	ps := len(t.Players)
	if ps < minPlayers {
		return fmt.Errorf("tournament needs %d or more players, got %d", minPlayers, ps)
	}

	if t.IsRunning() {
		return errors.New("tournament is already running")
	}

	// Set the two first matches
	m1 := NewMatch(t, qualifying)
	t.Matches = append(t.Matches, m1)

	m2 := NewMatch(t, qualifying)
	t.Matches = append(t.Matches, m2)

	// Add the first 8 players, in modulated joining order (first to
	// first match, second to second, third to first etc)
	for i, p := range t.Players[0:8] {
		ps := p.Player()
		err := t.Matches[i%2].AddPlayer(ps)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	t.Started = time.Now()

	// Get the first match and set the scheduled date to be now.
	err := t.Matches[0].SetTime(c, 0)
	if err != nil {
		return err
	}

	err = t.PublishNext()
	if err != nil && err != ErrPublishDisconnected {
		return err
	}
	return t.Persist()
}

// End marks the end of the tournament
func (t *Tournament) End() error {
	t.Ended = time.Now()

	err := t.server.publisher.SendTournamentComplete(t)
	if err != nil {
		return err
	}

	return t.Persist()
}

// EndQualifyingRounds marks when the qualifying rounds end
func (t *Tournament) EndQualifyingRounds(ts time.Time) error {
	t.QualifyingEnd = ts
	return t.Persist()
}

// UsurpTournament adds a batch of eight random players
func (t *Tournament) UsurpTournament() error {
	err := t.db.UsurpTournament(t, 8)
	if err != nil {
		return err
	}

	err = t.server.SendWebsocketUpdate("tournament", t)
	if err != nil {
		t.server.log.Error("Sending websocket update failed")
	}

	return err
}

// AutoplaySection runs through all the matches in the current section
// of matches
//
// E.g. if we're in the playoffs, they will all be finished and we're
// left at semi 1.
func (t *Tournament) AutoplaySection() error {
	if !t.IsRunning() {
		err := t.StartTournament(nil)
		if err != nil {
			return err
		}
	}

	m, err := t.CurrentMatch()
	if err != nil {
		return err
	}

	kind := m.Kind

	for kind == m.Kind {
		err := m.Autoplay()
		if err != nil {
			return err
		}

		if kind == final {
			// If we just finished the finals, then we should just exit
			break
		}

		m, err = t.CurrentMatch()
		if err != nil {
			return err
		}
	}

	return t.server.SendWebsocketUpdate("tournament", t)
}

// GetRunnerups gets the runnerups for this tournament
//
// The returned list is sorted descending by matches and ascending by
// score.
func (t *Tournament) GetRunnerups() ([]*PlayerSummary, error) {
	return t.db.GetRunnerups(t)
}

// MovePlayers moves the winner(s) of a Match into the next bracket of matches
// or into the Runnerup bracket.
func (t *Tournament) MovePlayers(m *Match) error {
	if m.Kind == qualifying {
		// If we have not yet set the qualifying end, we will keep making
		// matches and we have not passed it
		if t.QualifyingEnd.IsZero() || time.Now().Before(t.QualifyingEnd) {
			log.Print("Scheduling the next match")

			nm := NewMatch(t, qualifying)
			t.Matches = append(t.Matches, nm)
			rups, err := t.GetRunnerups()
			if err != nil {
				return err
			}

			// The runnerups are in order for the next match - add them
			for _, p := range rups {
				err = nm.AddPlayer(p.Player())
				if err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		}

		done, err := t.db.QualifyingMatchesDone(t)
		if err != nil {
			return errors.WithStack(err)
		}
		if done {
			err = t.ScheduleEndgame()
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	// For the playoffs, just place the winner into the final
	if m.Kind == playoff {
		p := SortByKills(m.Players)[0]
		err := t.Matches[len(t.Matches)-1].AddPlayer(p)
		if err != nil {
			return err
		}
	}

	return nil
}

// ScheduleEndgame makes the playoff matches and the final
func (t *Tournament) ScheduleEndgame() error {
	log.Print("Scheduling endgame")

	// Get sorted list of players that made it to the playoffs
	players, err := t.db.GetPlayoffPlayers(t)
	if err != nil {
		return errors.WithStack(err)
	}

	lp := len(players)
	if lp != 16 {
		return errors.New(fmt.Sprintf("Needed 16 players, got %d", lp))
	}

	// Bucket the players for inclusion in the playoffs
	buckets, err := DividePlayoffPlayers(players)
	if err != nil {
		return errors.WithStack(err)
	}

	// Add four playoff matches
	for x := 0; x < 4; x++ {
		m := NewMatch(t, playoff)

		// Add players to the match
		for _, p := range buckets[x] {
			err = m.AddPlayer(p.Player())
			if err != nil {
				return err
			}
		}

		t.Matches = append(t.Matches, m)
	}

	// Add the final, without players for now
	m := NewMatch(t, final)
	t.Matches = append(t.Matches, m)

	return nil
}

// NextMatch returns the next match
func (t *Tournament) NextMatch() (*Match, error) {
	if !t.IsRunning() {
		return nil, errors.New("tournament not running")
	}

	return t.db.NextMatch(t)
}

// CurrentMatch returns the current match
func (t *Tournament) CurrentMatch() (*Match, error) {
	return t.db.CurrentMatch(t)
}

// IsRunning returns boolean true if the tournament is running or not
func (t *Tournament) IsRunning() bool {
	return !t.Started.IsZero() && t.Ended.IsZero()
}

// IsInTournament checks if a player is allowed to join or is already in the tournament
func (t *Tournament) IsInTournament(ps *Person) bool {
	in, err := t.db.IsInTournament(t, ps)
	if err != nil {
		log.Printf("Getting tournament state failed: %+v", err)
		return false
	}
	return in
}

// GetCredits returns the credits object needed to display the credits roll.
func (t *Tournament) GetCredits() (*Credits, error) {
	// if t.Ended.IsZero() {
	// 	return nil, errors.New("cannot roll credits for unfinished tournament")
	// }

	// // TODO(thiderman): Many of these values are currently hardcoded or
	// // broadly grabs everything.
	// // We should move towards specifying these live when setting
	// // up the tournament itself.

	// executive := t.db.GetSafePerson("1279099058796903") // thiderman
	// producers := []*Person{
	// 	t.db.GetSafePerson("10153943465786915"), // GoosE
	// 	t.db.GetSafePerson("10154542569541289"), // Queen Obscene
	// 	t.db.GetSafePerson("10153964695568099"), // Karl-Astrid
	// 	t.db.GetSafePerson("10153910124391516"), // Hest
	// 	t.db.GetSafePerson("10154040229117471"), // Skolpadda
	// 	t.db.GetSafePerson("10154011729888111"), // Moijra
	// 	t.db.GetSafePerson("10154296655435218"), // Dalan
	// }

	// players := []*Person{
	// 	t.Winners[0].Person,
	// 	t.Winners[1].Person,
	// 	t.Winners[2].Person,
	// }
	// players = append(players, t.Runnerups...)

	// c := &Credits{
	// 	Executive:     executive,
	// 	Producers:     producers,
	// 	Players:       players,
	// 	ArchersHarmed: t.ArchersHarmed(),
	// }

	return &Credits{}, nil
}

// ArchersHarmed returns the amount of killed archers during the tournament
func (t *Tournament) ArchersHarmed() int {
	ret := 0
	for _, m := range t.Matches {
		ret += m.ArchersHarmed()
	}

	return ret
}

// GetPlayerSummary returns the tournament-wide player object.
func (t *Tournament) GetPlayerSummary(ps *Person) (*PlayerSummary, error) {
	return t.db.GetPlayerSummary(t, ps.PersonID)
}

// ByScheduleDate is a sort.Interface that sorts tournaments according
// to when they were scheduled.
type ByScheduleDate []*Tournament

func (s ByScheduleDate) Len() int {
	return len(s)
}
func (s ByScheduleDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByScheduleDate) Less(i, j int) bool {
	return s[i].Scheduled.Before(s[j].Scheduled)
}

// SortByScheduleDate returns a list in order of schedule date
func SortByScheduleDate(ps []*Tournament) []*Tournament {
	sort.Sort(ByScheduleDate(ps))
	return ps
}
