package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Person someone having a role in the tournament
type Person struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Email           string   `json:"email"`
	Nick            string   `json:"nick"`
	ColorPreference []string `json:"color_preference"`
	FacebookID      string   `json:"facebook_id"`
	FacebookToken   string   `json:"facebook_token"`
	AvatarURL       string   `json:"avatar_url"`
	Userlevel       int      `json:"userlevel"`
}

// Userlevels. Designed so that we can insert new ones in between them.
const (
	PermissionProducer    = 100
	PermissionCommentator = 50
	PermissionJudge       = 30
	PermissionPlayer      = 10
)

type score map[int]string

// ScoreSummary is a collection of scores for a Person
type ScoreSummary struct {
	Totals      score
	Tournaments map[string]score
}

func (p *Person) String() string {
	return fmt.Sprintf(
		"<Player %s (%s)>",
		p.Name,
		p.Nick,
	)
}

// JSON returns the person as a JSON representation
func (p *Person) JSON() (out []byte, err error) {
	out, err = json.Marshal(p)
	return
}

// Score gets the score of the Person
//
// Returned as a map of the total score and an array of maps - one per
// tournament participated in.
func (p *Person) Score() *ScoreSummary {
	return nil
}

// CreateFromFacebook adds a new player via Facebook login
func CreateFromFacebook(s *Server, req *FacebookAuthResponse) *Person {
	p := &Person{
		ID:            req.ID,
		FacebookID:    req.ID,
		FacebookToken: req.Token,
		Name:          req.Name,
		Email:         req.Email,
		Userlevel:     PermissionPlayer,
	}

	p.PrefillNickname()

	s.DB.SavePerson(p)

	return p
}

// PrefillNickname makes a suggestion to the nick based on the person
func (p *Person) PrefillNickname() {
	// TODO(thiderman): Move this into data files
	switch p.Name {
	case "Karl Johan Krantz":
		p.Nick = "Qrl-Astrid"
		p.Userlevel = PermissionProducer
	case "Ida Andreasson":
		p.Nick = "Blue"
		p.Userlevel = PermissionProducer
	case "Daniel Dala Tiderman":
		p.Nick = "Dala"
		p.Userlevel = PermissionProducer
	case "Lowe Thiderman":
		p.Nick = "thiderman"
		p.Userlevel = PermissionProducer
	case "Magnus Ulenius":
		p.Nick = "Goose"
		p.Userlevel = PermissionProducer
	case "Jonathan Gustafsson":
		p.Nick = "hest"
		p.Userlevel = PermissionProducer

	// Commentators
	case "Daniel McHugh":
		p.Nick = "Radcliffe"
		p.Userlevel = PermissionCommentator

	// Judges
	case "Daniele Sluijters":
		p.Nick = "Daenney"
		p.Userlevel = PermissionJudge

	// Other lovelies
	case "Agnes Skoog":
		p.Nick = "#swagnes"
	case "Mattias Aali Ahlström":
		p.Nick = "Aali"
	case "Sam Wise Ingberg":
		p.Nick = "Samselott"
	}
}

// UpdatePerson updates a person from a JoinRequest
func (p *Person) UpdatePerson(r *JoinRequest) {
	p.ID = r.ID
	p.Name = r.Name
	p.Nick = r.Nick
	p.ColorPreference = []string{r.Color}
}

// PreferredColor returns the preferred color
func (p *Person) PreferredColor() string {
	return p.ColorPreference[0]
}

// PersonFromSession returns the Person{} object attached to the session
func PersonFromSession(s *Server, r *http.Request) *Person {
	session, _ := CookieStore.Get(r, "session")
	id := session.Values["user"].(string)

	p := s.DB.GetPerson(id)
	return p
}
