package towerfall

import (
	"encoding/json"
	"log"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

var globalDB *Database

// LoadTournament loads a tournament from persisted JSON data
func LoadTournament(data []byte, s *Server) (t *Tournament, err error) {
	t = &Tournament{}
	err = json.Unmarshal(data, t)
	if err != nil {
		log.Print(err)
		return t, err
	}

	t.server = s
	t.db = s.DB

	t.SetMatchPointers()
	return
}

// NewDatabase sets up the database reader and writer
func NewDatabase(c *Config) (*Database, error) {
	log, _ := zap.NewDevelopment()

	pg := pg.Connect(&pg.Options{
		User:     c.DbUser,
		Database: c.DbName,
	})

	db := &Database{
		DB:  pg,
		log: log,
	}

	globalDB = db

	return db, nil
}
