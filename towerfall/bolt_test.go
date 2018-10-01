package towerfall

// func TestSaveTournament(t *testing.T) {
// 	assert := assert.New(t)
// 	fn := "persist.db"
// 	s, teardown := MockServer(t)
// 	defer teardown()
// 	db := s.DB

// 	id := "1241234"
// 	tm, err := NewTournament("hehe", id, "", time.Now().Add(time.Hour), nil, s)
// 	assert.Nil(err)

// 	db.SaveTournament(tm)
// 	db.Close()

// 	boltd, err := bolt.Open("test/"+fn, 0600, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ct := Tournament{}
// 	boltd.View(func(tx *bolt.Tx) error {
// 		b := tx.Bucket(TournamentKey)
// 		if b == nil {
// 			t.Fatal("bucket not created")
// 		}

// 		data := b.Get([]byte(id))
// 		err := json.Unmarshal(data, &ct)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		return nil
// 	})

// 	assert.Equal(ct.Name, tm.Name)
// 	assert.Equal(ct.Slug, tm.Slug)
// }

// func TestGetCurrentTournament(t *testing.T) {
// 	s, teardown := MockServer(t)
// 	defer teardown()
// 	db := s.DB

// 	_, err := NewTournament("not started", "not", "", time.Now().Add(time.Hour), nil, s)
// 	tm2, err := NewTournament("started", "go", "", time.Now().Add(time.Hour), nil, s)

// 	for i := 1; i <= 8; i++ {
// 		p := testPerson(i)
// 		s := NewPlayer(p).Summary()
// 		err := tm2.AddPlayer(&s)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		tm2.db.SavePerson(p)
// 	}

// 	err = tm2.StartTournament(nil)
// 	assert.NoError(t, err)

// 	t.Run("Get", func(t *testing.T) {
// 		tm3, err := db.GetCurrentTournament(s)
// 		assert.NoError(t, err)
// 		assert.Equal(t, tm3.Slug, tm2.Slug)
// 	})
// }
