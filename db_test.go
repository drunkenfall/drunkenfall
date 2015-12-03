package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// TestDatabase returns a clean test database
func MockDatabase(arg ...string) *Database {
	var fn string
	if len(arg) != 0 {
		fn = "test/" + arg[0]
	} else {
		fn = "test/test.db"
	}

	os.Remove(fn)
	os.Mkdir("test/", 0700)

	db, err := NewDatabase(fn)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func TestPersist(t *testing.T) {
	assert := assert.New(t)
	fn := "persist.db"
	db := MockDatabase(fn)

	id := "1241234"
	tm := Tournament{Name: "hehe", ID: id}

	db.Persist(&tm)
	db.Close()

	boltd, err := bolt.Open("test/"+fn, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	ct := Tournament{}
	boltd.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TournamentKey)
		if b == nil {
			t.Fatal("bucket not created")
		}

		data := b.Get([]byte(id))
		err := json.Unmarshal(data, &ct)
		if err != nil {
			t.Fatal(err)
		}

		return nil
	})

	assert.Equal(ct.Name, tm.Name)
	assert.Equal(ct.ID, tm.ID)
}
