package main

import (
	"encoding/json"
	"fmt"
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
}

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
	case "Ida Andreasson":
		p.Nick = "Queen Obscene"

	case "Daniel Dala Tiderman":
	case "Lowe Thiderman":
		p.Nick = "OP"

	case "Magnus Ulenius":
		p.Nick = "Goose"
	case "Agnes Skoog":
		p.Nick = "#swagnes"
	case "Jonathan Gustafsson":
		p.Nick = "hestxk"
	}
}

// UpdatePerson updates a person from a JoinRequest
func (p *Person) UpdatePerson(r *JoinRequest) {
	p.ID = r.ID
	p.Name = r.Name
	p.Nick = r.Nick
	p.ColorPreference = []string{r.Color}
}
