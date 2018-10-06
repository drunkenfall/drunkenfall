package towerfall

import (
	"encoding/json"
	"log"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

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
	var pg *gorm.DB
	var err error

	log, _ := zap.NewDevelopment()

	connStr := c.DbPostgresConn
	pg, err = gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatal("postgres open error", zap.Error(err))
	}

	db := &Database{
		DB:  pg,
		log: log,
	}

	return db, nil
}
