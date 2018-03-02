package towerfall

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

// MockServer returns a Server{} a with clean test Database{}
func MockServer(arg ...string) *Server {
	var fn string

	os.Mkdir("test/", 0755)
	if len(arg) != 0 {
		fn = "test/" + arg[0] // Use existing
	} else {
		fn = "test/test.db"
	}

	os.Remove(fn) // Clean it out
	db, err := NewDatabase(fn)
	if err != nil {
		log.Fatal(err)
	}
	db.LoadTournaments()

	s := NewServer(db)
	db.Server = s

	return s
}

func TestSaveTournament(t *testing.T) {
	assert := assert.New(t)
	fn := "persist.db"
	s := MockServer(fn)
	db := s.DB

	id := "1241234"
	tm, err := NewTournament("hehe", id, "", time.Now().Add(time.Hour), nil, s)
	assert.Nil(err)

	db.SaveTournament(tm)
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

func TestGetCurrentTournament(t *testing.T) {
	s := MockServer()
	db := s.DB

	_, err := NewTournament("not started", "not", "", time.Now().Add(time.Hour), nil, s)
	tm2, err := NewTournament("started", "go", "", time.Now().Add(time.Hour), nil, s)

	for i := 1; i <= 8; i++ {
		p := testPerson(i)
		err := tm2.AddPlayer(NewPlayer(p))
		if err != nil {
			log.Fatal(err)
		}
		tm2.db.SavePerson(p)
	}

	err = tm2.StartTournament(nil)
	assert.NoError(t, err)

	t.Run("Get", func(t *testing.T) {
		tm3, err := db.GetCurrentTournament()
		assert.NoError(t, err)
		assert.Equal(t, tm3.ID, tm2.ID)
	})
}
