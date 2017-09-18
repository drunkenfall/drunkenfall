package migrations

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/drunkenfall/drunkenfall/towerfall"
	"github.com/stretchr/testify/assert"
)

func loadmig4Tournaments(db, ro *bolt.DB) (*mig4curTournament, *mig4prevTournament) {
	id := "moon"
	orig := &mig4prevTournament{}
	tx, err := ro.Begin(false)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	b := tx.Bucket(towerfall.TournamentKey)
	bs := b.Get([]byte(id))
	_ = json.Unmarshal(bs, orig)

	mig := &mig4curTournament{}
	tx, err = db.Begin(false)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	b = tx.Bucket(towerfall.TournamentKey)
	bs = b.Get([]byte(id))
	_ = json.Unmarshal(bs, mig)

	return mig, orig
}

func TestMigration4(t *testing.T) {
	assert := assert.New(t)

	db, ro := migrationDB(4)
	defer db.Close()
	defer ro.Close()

	err := MigrateMatchCommitToRound(db)
	assert.Nil(err)

	v, err := getVersion(db)
	assert.Nil(err)
	assert.Equal(5, v)

	mig, orig := loadmig4Tournaments(db, ro)

	assert.Equal(len(mig.Final.Rounds), len(orig.Final.Commits))
	for x, m := range mig.Tryouts {
		o := orig.Tryouts[x]
		assert.Equal(len(m.Rounds), len(o.Commits))
	}
	for x, m := range mig.Semis {
		o := orig.Semis[x]
		assert.Equal(len(m.Rounds), len(o.Commits))
	}
}
