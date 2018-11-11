package towerfall

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"

	"github.com/deckarep/golang-set"
	"github.com/pkg/errors"
)

// AllColors is a list of the available player colors
var AllColors = []interface{}{
	"green",
	"blue",
	"pink",
	"orange",
	"white",
	"yellow",
	"cyan",
	"purple",
	"red",
}

// Colors is the definitive set of all the colors
var Colors = mapset.NewSetFromSlice(AllColors)

const scoreMultiplier = 7
const scoreSweep = (97 * scoreMultiplier)
const scoreKill = (21 * scoreMultiplier)
const scoreSelf = (-35 * scoreMultiplier) // Negative 1.66 times of a kill
const scoreWinner = (350 * scoreMultiplier)
const scoreSecond = (150 * scoreMultiplier)
const scoreThird = (70 * scoreMultiplier)
const scoreFourth = (30 * scoreMultiplier)

const finalMultiplier = 2.5
const finalExponential = 1.05

// ScoreData is a structured Key/Value pair list for scores
type ScoreData struct {
	Key    string
	Value  int
	Player *Player
}

// Player is a representation of one player in a match
type Player struct {
	ID             uint `json:"id"`
	MatchID        uint
	PersonID       string      `sql:",pk"`
	Person         *Person     `json:"person" sql:"-"`
	Nick           string      `json:"nick"`
	Color          string      `json:"color"`
	PreferredColor string      `json:"preferred_color"`
	ArcherType     int         `json:"archer_type" sql:",notnull"`
	Shots          int         `json:"shots" sql:",notnull"`
	Sweeps         int         `json:"sweeps" sql:",notnull"`
	Kills          int         `json:"kills" sql:",notnull"`
	Self           int         `json:"self" sql:",notnull"`
	MatchScore     int         `json:"match_score" sql:",notnull"`
	TotalScore     int         `json:"total_score" sql:",notnull"`
	State          PlayerState `json:"state" sql:"-"`
	Match          *Match      `json:"-" sql:"-"`
	DisplayNames   []string    `sql:",array"`
}

// A PlayerSummary is a tournament-wide summary of the scores a player has
type PlayerSummary struct {
	ID uint `json:"id"`

	TournamentID uint
	PersonID     string  `json:"person_id"`
	Person       *Person `json:"person" sql:"-"`
	Shots        int     `json:"shots" sql:",notnull"`
	Sweeps       int     `json:"sweeps" sql:",notnull"`
	Kills        int     `json:"kills" sql:",notnull"`
	Self         int     `json:"self" sql:",notnull"`
	Matches      int     `json:"matches" sql:",notnull"`
	TotalScore   int     `json:"score" sql:",notnull"`
	SkillScore   int     `json:"skill_score" sql:",notnull"`
}

type PlayerState struct {
	Arrows    Arrows `json:"arrows"`
	Shield    bool   `json:"shield"`
	Wings     bool   `json:"wings"`
	Hat       bool   `json:"hat"`
	Invisible bool   `json:"invisible"`
	Speed     bool   `json:"speed"`
	Alive     bool   `json:"alive"`
	Lava      bool   `json:"lava"`
	Killer    int    `json:"killer" sql:",notnull"`
}

// NewPlayer returns a new instance of a player
func NewPlayer(ps *Person) *Player {
	p := &Player{
		PersonID:       ps.PersonID,
		Person:         ps,
		Nick:           ps.Nick,
		ArcherType:     ps.ArcherType,
		State:          NewPlayerState(),
		PreferredColor: ps.PreferredColor,
		DisplayNames:   ps.DisplayNames,
	}

	if p.PreferredColor != "" {
		p.PreferredColor = ps.PreferredColor
	} else {
		p.PreferredColor = RandomColor(Colors)
	}

	return p
}

func NewPlayerState() PlayerState {
	ps := PlayerState{
		Arrows: make(Arrows, 0),
		Alive:  true,
		Hat:    true,
		Killer: -2,
	}
	return ps
}

// NewPlayerSummary returns a new instance of a tournament player
func NewPlayerSummary(ps *Person) *PlayerSummary {
	p := &PlayerSummary{
		PersonID: ps.PersonID,
		Person:   ps,
	}
	return p
}

func (p *Player) String() string {
	return fmt.Sprintf(
		"<(%d) %s: %dsh %dsw %dk %ds>",
		p.ID,
		p.Nick,
		p.Shots,
		p.Sweeps,
		p.Kills,
		p.Self,
	)
}

// Name returns the nickname
func (p *Player) Name() string {
	return p.Nick
}

// NumericColor is the numeric representation of the color the player has
func (p *Player) NumericColor() int {
	for x, c := range AllColors {
		if p.Color == c {
			return x
		}
	}

	// No color was found - this is a bug. Return default.
	log.Printf("Player '%s' did not match a color for '%s'", p.Nick, p.Color)
	return 0
}

// Score calculates the score to determine runnerup positions.
func (p *PlayerSummary) Score() (out int) {
	out += p.Sweeps * scoreSweep
	out += p.Kills * scoreKill
	out += p.Self * scoreSelf

	// Negative score is not allowed
	if out <= 0 {
		out = 0
	}

	return
}

// Score calculates the score to determine runnerup positions.
func (p *Player) Score() (out int) {
	out += p.Sweeps * scoreSweep
	out += p.Kills * scoreKill
	out += p.Self * scoreSelf

	// Negative score is not allowed
	if out <= 0 {
		out = 0
	}

	// Match score is added afterwards so that no one is stuck on 0.
	out += p.MatchScore

	return
}

// Summary resturns a Summary{} object for the player
func (p *Player) Summary() PlayerSummary {
	return PlayerSummary{
		PersonID: p.PersonID,
		Person:   p.Person,
		Shots:    p.Shots,
		Sweeps:   p.Sweeps,
		Kills:    p.Kills,
		Self:     p.Self,
	}
}

// Player returns a new Player{} object from the summary
func (p *PlayerSummary) Player() Player {
	// TODO(thiderman): It would be better to always make sure that the
	// person object is set
	if p.Person == nil {
		var err error
		p.Person, err = globalDB.GetPerson(p.PersonID)
		if err != nil {
			log.Fatal(err)
		}
	}

	return *NewPlayer(p.Person)
}

// ScoreData returns this players set of ScoreData
func (p *Player) ScoreData() []ScoreData {
	sd := []ScoreData{
		{Key: "kills", Value: p.Kills, Player: p},
		{Key: "shots", Value: p.Shots, Player: p},
		{Key: "sweeps", Value: p.Sweeps, Player: p},
		{Key: "self", Value: p.Self, Player: p},
	}
	return sd
}

// URL returns the URL to this player
//
// Used for action URLs
func (p *Player) URL() string {
	out := fmt.Sprintf(
		"%s/%d",
		p.Match.URL(),
		p.Index(),
	)
	return out
}

// Index returns the index in the current match
func (p *Player) Index() int {
	if p.Match != nil {
		for i, o := range p.Match.Players {
			if p.Nick == o.Name() {
				return i
			}
		}
	}
	return -1
}

// AddShot increases the shot count
func (p *Player) AddShot() {
	p.Shots++
}

// RemoveShot decreases the shot count
// Fails silently if shots are zero.
func (p *Player) RemoveShot() {
	if p.Shots == 0 {
		log.Printf("Not removing shot from %s; already at zero", p.Nick)
		return
	}
	p.Shots--
}

// AddSweep increases the sweep count
func (p *Player) AddSweep() {
	p.Sweeps++
}

// AddKills increases the kill count and adds a sweep if necessary
func (p *Player) AddKills(kills int) {
	p.Kills += kills
	if kills == 3 {
		p.AddSweep()
	}
}

// RemoveKill decreases the kill count
// Doesn't to anything if kills are at zero.
func (p *Player) RemoveKill() {
	if p.Kills == 0 {
		log.Printf("Not removing kill from %s; already at zero", p.Nick)
		return
	}
	p.Kills--
}

// AddSelf increases the self count and decreases the kill
func (p *Player) AddSelf() {
	p.Self++
	p.RemoveKill()
}

// Reset resets the stats on a PlayerSummary to 0
//
// It is to be run in Match.Start()
func (p *Player) Reset() {
	p.Shots = 0
	p.Sweeps = 0
	p.Kills = 0
	p.Self = 0
	p.State = NewPlayerState()
}

// HTML renders the HTML of a player
func (p *Player) HTML() (out string) {
	return
}

// ByColorConflict is a sort.Interface that sorts players by their score
type ByColorConflict []PlayerSummary

func (s ByColorConflict) Len() int { return len(s) }

func (s ByColorConflict) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByColorConflict) Less(i, j int) bool {
	if s[i].Person.Userlevel != s[j].Person.Userlevel {
		return s[i].Person.Userlevel > s[j].Person.Userlevel
	}
	return s[i].SkillScore > s[j].SkillScore
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

// SortByColorConflicts returns a list in an unspecified order,
// Probably by User level and then score.
func SortByColorConflicts(m *Match, ps []Person) (tmp []PlayerSummary, err error) {
	var tp *PlayerSummary
	tmp = make([]PlayerSummary, len(ps))
	for i, p := range ps {
		// TODO(thiderman): This is not very elegant and should be replaced.
		tp, err = m.Tournament.GetPlayerSummary(&p)
		if err != nil {
			return tmp, errors.WithStack(err)
		}
		tmp[i] = *tp
	}
	sort.Sort(ByColorConflict(tmp))
	return
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
	iKills := s[i].Kills
	jKills := s[j].Kills
	if iKills == jKills {
		return s[i].Score() > s[j].Score()
	}
	return iKills > jKills
}

// SortByKills returns a list in order of the kills the players have
func SortByKills(ps []Player) []Player {
	tmp := make([]Player, len(ps))
	copy(tmp, ps)
	sort.Sort(ByKills(tmp))
	return tmp
}

// ByRunnerup is a sort.Interface that sorts players by their runnerup status
//
// This is almost exactly the same as score, but the number of matches a player
// has played factors in, and players that have played less matches are sorted
// favorably.
type ByRunnerup []PlayerSummary

func (s ByRunnerup) Len() int {
	return len(s)

}
func (s ByRunnerup) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]

}
func (s ByRunnerup) Less(i, j int) bool {
	if s[i].Matches == s[j].Matches {
		// Higher is better - if you have a high score you'll play again
		return s[i].Score() > s[j].Score()
	}
	// Lower is better - the ones that have not played should be at the top
	return s[i].Matches < s[j].Matches
}

// SortByRunnerup returns a list in order of the kills the players have
func SortByRunnerup(ps []PlayerSummary) []PlayerSummary {
	sort.Sort(ByRunnerup(ps))
	return ps
}

// RandomColor returns a random color from the ColorList
func RandomColor(s mapset.Set) string {
	colors := s.ToSlice()
	x := len(colors)
	return colors[rand.Intn(x)].(string)
}

// AvailableColors returns a ColorList with the colors not used in a match
func AvailableColors(m *Match) mapset.Set {
	colors := mapset.NewSetFromSlice(AllColors)
	ret := colors.Difference(m.presentColors)
	return ret
}

// Reset resets the stats on a PlayerSummary to 0
//
// It is to be run in Match.Start()
func (p *PlayerSummary) Reset() {
	p.Shots = 0
	p.Sweeps = 0
	p.Kills = 0
	p.Self = 0
	p.Matches = 0
}

// Update updates a player with the scores of another
//
// This is primarily used by the tournament score calculator
func (p *PlayerSummary) Update(other PlayerSummary) {
	p.Shots += other.Shots
	p.Sweeps += other.Sweeps
	p.Kills += other.Kills
	p.Self += other.Self
	p.TotalScore = p.Score()

	// Every call to this method is per match. Count every call
	// as if a match.
	p.Matches++
	// log.Printf("Updated player: %d, %d", p.TotalScore, p.Matches)
}

// DividePlayoffPlayers divides the playoff players into four buckets based
// on their score
//
// The input is expected to be sorted with score descending
func DividePlayoffPlayers(ps []*PlayerSummary) ([][]*PlayerSummary, error) {
	ret := [][]*PlayerSummary{
		[]*PlayerSummary{},
		[]*PlayerSummary{},
		[]*PlayerSummary{},
		[]*PlayerSummary{},
	}

	for x, p := range ps {
		ret[x%4] = append(ret[x%4], p)
	}

	return ret, nil
}

// FinalMultiplier returns the multiplier used for the winner scores
// in the final
//
// The longer at tournament lasts, the more points you'll get for
// winning the final.
func FinalMultiplier(numMatches int) float64 {
	// We only count extra when there has been more than 16 matches
	x := numMatches - 16

	// If there haven't been, just return the default
	if x <= 0 {
		return finalMultiplier
	}

	return finalMultiplier * math.Pow(finalExponential, float64(x))
}
