package towerfall

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strings"

	"github.com/deckarep/golang-set"
	"github.com/fatih/camelcase"
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

// ScoreData is a structured Key/Value pair list for scores
type ScoreData struct {
	Key    string
	Value  int
	Player *Player
}

// Player is a Participant that is actively participating in battles.
type Player struct {
	Person         *Person     `json:"person"`
	Color          string      `json:"color"`
	PreferredColor string      `json:"preferred_color"`
	ArcherType     int         `json:"archer_type"`
	Shots          int         `json:"shots"`
	Sweeps         int         `json:"sweeps"`
	Kills          int         `json:"kills"`
	Self           int         `json:"self"`
	Matches        int         `json:"matches"`
	TotalScore     int         `json:"score"`
	State          PlayerState `json:"state"`
	Match          *Match      `json:"-"`
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
	Killer    int    `json:"killer"`
}

// NewPlayer returns a new instance of a player
func NewPlayer(ps *Person) *Player {
	p := &Player{
		Person:     ps,
		ArcherType: ps.ArcherType,
		State:      NewPlayerState(),
	}
	if len(ps.ColorPreference) > 0 {
		p.PreferredColor = ps.ColorPreference[0]
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

func (p *Player) String() string {
	return fmt.Sprintf(
		"<%s: %dsh %dsw %dk %ds>",
		p.Name(),
		p.Shots,
		p.Sweeps,
		p.Kills,
		p.Self,
	)
}

// Name returns the nickname
func (p *Player) Name() string {
	return p.Person.Nick
}

// NumericColor is the numeric representation of the color the player has
func (p *Player) NumericColor() int {
	for x, c := range AllColors {
		if p.Color == c {
			return x
		}
	}

	// No color was found - this is a bug. Return default.
	log.Printf("Player '%s' did not match a color for '%s'", p.Name(), p.Color)
	return 0
}

// Score calculates the score to determine runnerup positions.
func (p *Player) Score() (out int) {
	// This algorithm is probably flawed, but at least it should be able to
	// determine who is the most entertaining.
	// When executed, a sweep is basically 14 points since scoring a sweep
	// also comes with a shot and three kills.

	out += p.Sweeps * 5
	out += p.Shots * 3
	out += p.Kills * 2
	out += p.Self

	return
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

// Classes returns the CSS color classes for the player
func (p *Player) Classes() string {
	if p.Match != nil && p.Match.IsEnded() {
		ps := ByScore(p.Match.Players)
		if ps[0].Name() == p.Name() {
			// Always gold for the winner
			return "gold"
		} else if ps[1].Name() == p.Name() {
			// Silver for the second, unless there is a short amount of playoffs
			if p.Match.Kind != playoff || len(p.Match.tournament.Matches)-3 <= 4 {
				return "silver"
			}
		} else if ps[2].Name() == p.Name() && p.Match.Kind == final {
			return "bronze"
		}

		return "out"
	}
	return p.PreferredColor
}

// Index returns the index in the current match
func (p *Player) Index() int {
	if p.Match != nil {
		for i, o := range p.Match.Players {
			if p.Name() == o.Name() {
				return i
			}
		}
	}
	return -1
}

// AddShot increases the shot count
func (p *Player) AddShot() {
	log.Printf("Adding shot to %s", p.Name())
	p.Shots++
}

// RemoveShot decreases the shot count
// Fails silently if shots are zero.
func (p *Player) RemoveShot() {

	if p.Shots == 0 {
		log.Printf("Not removing shot from %s; already at zero", p.Name())
		return
	}
	log.Printf("Removing shot from %s", p.Name())
	p.Shots--
}

// AddSweep increases the sweep count
func (p *Player) AddSweep() {
	log.Printf("Adding sweep to %s", p.Name())
	p.Sweeps++
}

// AddKills increases the kill count and adds a sweep if necessary
func (p *Player) AddKills(kills int) {
	log.Printf("Adding %d kills to %s", kills, p.Name())
	p.Kills += kills
	if kills == 3 {
		p.AddSweep()
	}
}

// RemoveKill decreases the kill count
// Fails silently if kills are zero.
func (p *Player) RemoveKill() {
	if p.Kills == 0 {
		log.Printf("Not removing kill from %s; already at zero", p.Name())
		return
	}
	log.Printf("Removing kill from %s", p.Name())
	p.Kills--
}

// AddSelf increases the self count and decreases the kill
func (p *Player) AddSelf() {
	log.Printf("Adding self to %s", p.Name())
	p.Self++
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
	p.Matches = 0
	p.State = NewPlayerState()
}

// Update updates a player with the scores of another
//
// This is primarily used by the tournament score calculator
func (p *Player) Update(other Player) {
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

// HTML renders the HTML of a player
func (p *Player) HTML() (out string) {
	return
}

// ByColorConflict is a sort.Interface that sorts players by their score
type ByColorConflict []Player

func (s ByColorConflict) Len() int { return len(s) }

func (s ByColorConflict) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByColorConflict) Less(i, j int) bool {
	if s[i].Person.Userlevel != s[j].Person.Userlevel {
		return s[i].Person.Userlevel > s[j].Person.Userlevel
	}
	return s[i].Score() > s[j].Score()
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
func SortByColorConflicts(ps []Player) (tmp []Player, err error) {
	var tp *Player
	tmp = make([]Player, len(ps))
	for i, p := range ps {
		// TODO(thiderman): This is not very elegant and should be replaced.
		tp, err = p.Match.tournament.getTournamentPlayerObject(p.Person)
		if err != nil {
			return
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
		return s[i].Score() > s[j].Score()
	}
	// Lower is better - the ones that have not played should be at the top
	return s[i].Matches < s[j].Matches
}

// SortByRunnerup returns a list in order of the kills the players have
func SortByRunnerup(ps []Player) []Player {
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

// transformName
func transformName(name string) (string, string) {
	name = strings.TrimSpace(name)

	rxp := regexp.MustCompile("[^0-9A-Za-z()]+")
	// Run it twice so that double spaces are removed
	name = rxp.ReplaceAllString(name, " ")
	name = rxp.ReplaceAllString(name, " ")

	if strings.Contains(name, " ") {
		spl := strings.Split(name, " ")

		if len(spl) == 2 {
			// Perfect - two spaces; just do it!
			return spl[0], spl[1]
		}

		if (spl[0]) == "Mr" {
			return spl[1], spl[2]
		}

		if len(spl) == 3 {
			// Sigh.
			if strings.HasPrefix(spl[2], "(") {
				return spl[0], spl[1]
			}

			front := fmt.Sprintf("%s %s", spl[0], spl[1])
			back := fmt.Sprintf("%s %s", spl[1], spl[2])
			if len(front) >= len(spl[2]) {
				return front, spl[2]
			} else {
				return spl[0], back
			}
		}
	}

	spl := camelcase.Split(name)
	if len(spl) > 1 {
		// Oh noes, CamelCase!
		if len(spl) == 2 {
			return spl[0], spl[1]
		}

		if spl[0] == "The" {
			return fmt.Sprintf("%s %s", spl[0], spl[1]), spl[2]
		}
	}

	prefixes := []string{
		"The",
		"El Grande",
		"Sweet",
		"Big",
		"Strong",
		"Boss",
		"Master",
		"Vicious",
		"Vengeful",
		"Sassy",
		"Gay",
		"Yaaas",
		"Epic",
		"Grandmaster",
		"Ser",
		"Jarl of",
		"Bishop",
		"Lady",
		"Lord",
		"Duke",
		"Their Majesty",
		"Royal",
		"Brother",
		"Sister",
		"Governor",
		"Papa",
		"Monsieur",
		"Madamoiselle",
		"1337",
		"Motherfuckin",
		"Wannabe",
		"Call Me",
		"Mr.",
		"Ms.",
	}

	suffixes := []string{
		"The Great",
		"The Not So Great",
		"K?",
		"AS FUCK",
		"The Goon",
		"The Fucker",
		"Fuck Yeah",
		"for the win",
		"for the lulz",
		"The Supreme",
		"Yaaaas",
		"The Fabulous",
		"The Gay",
		"The Unicorn",
		"The Wannabe",
		"The Usurper",
		"Esq.",
		"2000",
		"The Idiot",
		"The Drunk",
		"The Sober",
		"Is Dead",
		"Is Lame",
		"Sucks",
		"The Homeless",
	}

	// One-namer! Add a prefix or a suffix; 50% distrib
	if rand.Intn(100) >= 50 {
		return prefixes[rand.Intn(len(prefixes))], name
	}
	return name, suffixes[rand.Intn(len(suffixes))]
}
