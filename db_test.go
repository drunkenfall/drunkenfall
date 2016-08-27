package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// MockServer returns a Server{} a with clean test Database{}
func MockServer(arg ...string) *Server {
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

	s := NewServer(db)

	return s
}

func TestSaveTournament(t *testing.T) {
	assert := assert.New(t)
	fn := "persist.db"
	db := MockServer(fn).DB

	id := "1241234"
	tm := Tournament{Name: "hehe", ID: id}

	db.SaveTournament(&tm)
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
