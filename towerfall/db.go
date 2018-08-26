package towerfall

import (
	"encoding/json"
	"log"
)

type Database interface {
	LoadTournaments() error
	LoadPeople() error

	SaveTournament(t *Tournament) error
	OverwriteTournament(t *Tournament) error

	SavePerson(p *Person) error
	GetPerson(id string) (*Person, error)
	GetSafePerson(id string) *Person
	DisablePerson(id string) error

	GetCurrentTournament() (*Tournament, error)
	ClearTestTournaments() error

	Close() error
}

// LoadTournament loads a tournament from persisted JSON data
func LoadTournament(data []byte, db *BoltDatabase) (t *Tournament, e error) {
	t = &Tournament{}
	err := json.Unmarshal(data, t)
	if err != nil {
		log.Print(err)
		return t, err
	}

	t.db = db
	t.server = db.Server

	t.SetMatchPointers()
	return
}
