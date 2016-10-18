package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// loadmig2Tournaments
func loadmig2Tournaments(db, ro *bolt.DB) (*mig2prevTournament, *mig2curTournament) {
	id := "moon"
	orig := &mig2prevTournament{}
	tx, err := ro.Begin(false)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	b := tx.Bucket(TournamentKey)
	bs := b.Get([]byte(id))
	_ = json.Unmarshal(bs, orig)

	mig := &mig2curTournament{}
	tx, err = db.Begin(false)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	b = tx.Bucket(TournamentKey)
	bs = b.Get([]byte(id))
	_ = json.Unmarshal(bs, mig)

	return orig, mig
}

func TestMigration2(t *testing.T) {
	assert := assert.New(t)

	db, ro := migrationDB(1)
	err := MigrateTournamentRunnerupStringPerson(db)
	assert.Nil(err)

	v, err := getVersion(db)
	assert.Nil(err)
	assert.Equal(3, v)
	orig, mig := loadmig2Tournaments(db, ro)

	assert.Equal(len(orig.Runnerups), len(mig.Runnerups))

	// Check all the runnerups
	for x, n := range orig.Runnerups {
		assert.Equal(n, mig.Runnerups[x].Nick)
	}

	db.Close()
	ro.Close()
}
