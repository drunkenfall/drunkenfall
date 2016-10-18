package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func loadmig3Tournaments(db, ro *bolt.DB) *mig3prevTournament {
	id := "moon"
	orig := &mig3prevTournament{}
	tx, err := db.Begin(false)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	b := tx.Bucket(TournamentKey)
	bs := b.Get([]byte(id))
	_ = json.Unmarshal(bs, orig)

	return orig
}

func TestMigration3(t *testing.T) {
	assert := assert.New(t)

	db, ro := migrationDB(3)
	err := MigrateMatchScoreOrder(db)
	assert.Nil(err)

	v, err := getVersion(db)
	assert.Nil(err)
	assert.Equal(4, v)

	mig := loadmig3Tournaments(db, ro)

	// Manually checked in numbers from the matches
	assert.Equal([]int{1, 0, 2, 3}, mig.Tryouts[0].ScoreOrder)
	assert.Equal([]int{0, 3, 2, 1}, mig.Tryouts[1].ScoreOrder)
	assert.Equal([]int{2, 1, 3, 0}, mig.Tryouts[7].ScoreOrder)

	assert.Equal([]int{0, 2, 3, 1}, mig.Semis[0].ScoreOrder)
	assert.Equal([]int{1, 3, 0, 2}, mig.Semis[1].ScoreOrder)

	assert.Equal([]int{0, 3, 2, 1}, mig.Final.ScoreOrder)

	db.Close()
	ro.Close()
}
