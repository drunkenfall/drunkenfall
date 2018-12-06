package towerfall

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// This disables websocket pushing during the tests
	broadcasting = false
}

func TestMain(m *testing.M) {
	// This makes sure that the test output just has the filenames, making it
	// easier for tools that parses it to find where the log output happened.
	log.SetFlags(log.Lshortfile)
	os.Exit(m.Run())
}

func percentTrue(n int) bool {
	return rand.Intn(100) <= n
}

func runTestMatch(t *testing.T, tm *Tournament, index int, checkCounts bool) {
	t.Run("Play", func(t *testing.T) {
		nm, err := tm.NextMatch()
		assert.NoError(t, err)
		if !assert.NotNil(t, nm) {
			t.Fatal()
		}

		assert.Equal(t, index, nm.Index)

		err = nm.Autoplay()
		assert.NoError(t, err)
	})

	if checkCounts {
		t.Run("Match counts for players are set", func(t *testing.T) {

			m := 4 * (index + 1)
			players := len(tm.Players)
			minMatches := m / players
			played := m % players

			ps := []*PlayerSummary{}
			err := tm.db.DB.Model(&ps).Where("tournament_id = ? AND matches = ?", tm.ID, minMatches+1).Select()
			assert.NoError(t, err)

			assert.Equal(t, played, len(ps))
		})
	}

	done, err := globalDB.QualifyingMatchesDone(tm)
	assert.NoError(t, err)

	if !done {
		t.Run("Runnerup order", func(t *testing.T) {
			rups, err := tm.GetRunnerups()
			assert.NoError(t, err)
			assert.Equal(t, 4, len(rups))
		})
	}

}

// testTournament ma s,kes a test tournament with `count` players.
func testTournament(t *testing.T, server *Server, count int) (tm *Tournament) {
	usedPeople = []string{}
	s := strconv.Itoa(count)
	tm, err := NewTournament(t.Name(), s, "cover", time.Now().Add(time.Hour), nil, server)
	if err != nil {
		t.Fatalf("tournament creation failed")
	}

	for i := 1; i <= count; i++ {
		p := testPerson(server)
		s := NewPlayer(p).Summary()

		err := tm.AddPlayer(&s)
		if err != nil {
			t.Fatalf("adding player failed: %+v", err)
			return
		}
	}

	err = tm.Persist()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func TestQualifyingFlowNoNewJoiners(t *testing.T) {
	players := 19
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, players)

	t.Run("Starting", func(t *testing.T) {
		err := tm.StartTournament(nil)
		assert.NoError(t, err)
	})

	t.Run("Matches set", func(t *testing.T) {
		assert.Equal(t, 2, len(tm.Matches))
	})

	t.Run("Players set in matches", func(t *testing.T) {
		assert.Equal(t, 4, len(tm.Matches[0].Players))
		assert.Equal(t, 4, len(tm.Matches[1].Players))
	})

	for x := 0; x < 20; x++ {
		t.Run(fmt.Sprintf("Match %d", x+1), func(t *testing.T) { runTestMatch(t, tm, x, true) })
	}
}

func TestQualifyingFlowWithLateJoiners(t *testing.T) {
	players := 19
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, players)

	t.Run("Starting", func(t *testing.T) {
		err := tm.StartTournament(nil)
		assert.NoError(t, err)
	})

	t.Run("Matches set", func(t *testing.T) {
		assert.Equal(t, 2, len(tm.Matches))
	})

	t.Run("Players set in matches", func(t *testing.T) {
		assert.Equal(t, 4, len(tm.Matches[0].Players))
		assert.Equal(t, 4, len(tm.Matches[1].Players))
	})

	for x := 0; x < 100; x++ {
		if percentTrue(40) {
			for x := 0; x < rand.Intn(4); x++ {

				p := testPerson(s)
				s := NewPlayer(p).Summary()
				t.Logf("Adding player: %s", p.Nick)

				err := tm.AddPlayer(&s)
				assert.NoError(t, err)
			}
		}

		t.Run(fmt.Sprintf("Match %d", x+1), func(t *testing.T) { runTestMatch(t, tm, x, false) })
	}
}

func TestFullTournament(t *testing.T) {
	players := 19
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, players)

	t.Run("Starting", func(t *testing.T) {
		err := tm.StartTournament(nil)
		if !assert.NoError(t, err) {
			t.Fatal(err)
		}
	})

	t.Run("Matches set", func(t *testing.T) {
		assert.Equal(t, 2, len(tm.Matches))
	})

	t.Run("Players set in matches", func(t *testing.T) {
		assert.Equal(t, 4, len(tm.Matches[0].Players))
		assert.Equal(t, 4, len(tm.Matches[1].Players))
	})

	matches := 16
	for x := 0; x < matches; x++ {
		t.Run(fmt.Sprintf("Match %d", x+1), func(t *testing.T) { runTestMatch(t, tm, x, false) })
	}

	t.Run("Schedule qualifying end", func(t *testing.T) {
		err := tm.EndQualifyingRounds(time.Now())
		assert.NoError(t, err)
	})

	// Run the last two matches
	runTestMatch(t, tm, matches, false)
	runTestMatch(t, tm, matches+1, false)

	t.Run("All qualifying matches ended", func(t *testing.T) {
		ret, err := tm.db.QualifyingMatchesDone(tm)
		assert.NoError(t, err)
		assert.True(t, ret)
	})

	t.Run("Four playoffs and a funeral", func(t *testing.T) {
		playoffs, err := tm.db.GetMatches(tm, playoff)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(playoffs))

		t.Run("Players scheduled for playoffs", func(t *testing.T) {
			for _, m := range playoffs {
				log.Printf("%+v", m)
				assert.Equal(t, 4, len(m.Players))
				assert.Equal(t, 20, m.Length)
				assert.Equal(t, "A", m.Ruleset)
			}
		})

		t.Run("Final scheduled", func(t *testing.T) {
			finals, err := tm.db.GetMatches(tm, final)
			assert.NoError(t, err)

			if !assert.Equal(t, 1, len(finals)) {
				t.Fatal("final not set")
			}

			f := finals[0]
			assert.Equal(t, f.Level, "cataclysm")
			assert.Equal(t, f.Ruleset, "B")
			assert.Equal(t, f.Length, 20)
		})
	})

	t.Run("First playoff", func(t *testing.T) {
		runTestMatch(t, tm, matches+2, false)

		finals, err := tm.db.GetMatches(tm, final)
		assert.NoError(t, err)
		if !assert.Equal(t, 1, len(finals)) {
			t.Fatal("no final set")
		}

		t.Run("One player sent to finals", func(t *testing.T) {
			assert.Equal(t, 1, len(finals[0].Players))
		})
	})

	t.Run("Second playoff", func(t *testing.T) {
		runTestMatch(t, tm, matches+3, false)

		finals, err := tm.db.GetMatches(tm, final)
		assert.NoError(t, err)
		if !assert.Equal(t, 1, len(finals)) {
			t.Fatal("no final set")
		}

		t.Run("Two players sent to finals", func(t *testing.T) {
			assert.Equal(t, 2, len(finals[0].Players))
		})
	})

	t.Run("Third playoff", func(t *testing.T) {
		runTestMatch(t, tm, matches+4, false)

		finals, err := tm.db.GetMatches(tm, final)
		assert.NoError(t, err)
		if !assert.Equal(t, 1, len(finals)) {
			t.Fatal("no final set")
		}

		t.Run("Three players sent to finals", func(t *testing.T) {
			assert.Equal(t, 3, len(finals[0].Players))
		})
	})

	t.Run("Last playoff", func(t *testing.T) {
		runTestMatch(t, tm, matches+5, false)

		finals, err := tm.db.GetMatches(tm, final)
		assert.NoError(t, err)
		if !assert.Equal(t, 1, len(finals)) {
			t.Fatal("no final set")
		}

		t.Run("All players sent to finals", func(t *testing.T) {
			assert.Equal(t, 4, len(finals[0].Players))
		})
	})

	t.Run("Fhfppfphiinhnallehheh", func(t *testing.T) {
		runTestMatch(t, tm, matches+6, false)
	})

	t.Run("Tournament end state", func(t *testing.T) {
		t.Run("End is set", func(t *testing.T) {
			assert.NotZero(t, tm.Ended)
		})
	})
}

func TestDoubleStartIsForbidden(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	err := tm.StartTournament(nil)
	assert.NoError(err)
	err = tm.StartTournament(nil)
	assert.EqualError(err, "tournament is already running")
}

func TestPlayoffPlayerDistribution(t *testing.T) {
	players := []*PlayerSummary{}
	for x := 16; x > 0; x-- {
		players = append(players, &PlayerSummary{ID: uint(x), SkillScore: x})
	}

	buckets, err := DividePlayoffPlayers(players)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(buckets))

	assert.Equal(t, 16, buckets[0][0].SkillScore)
	assert.Equal(t, 15, buckets[1][0].SkillScore)
	assert.Equal(t, 14, buckets[2][0].SkillScore)
	assert.Equal(t, 13, buckets[3][0].SkillScore)

	assert.Equal(t, 12, buckets[0][1].SkillScore)
	assert.Equal(t, 11, buckets[1][1].SkillScore)
	assert.Equal(t, 10, buckets[2][1].SkillScore)
	assert.Equal(t, 9, buckets[3][1].SkillScore)

	assert.Equal(t, 8, buckets[0][2].SkillScore)
	assert.Equal(t, 7, buckets[1][2].SkillScore)
	assert.Equal(t, 6, buckets[2][2].SkillScore)
	assert.Equal(t, 5, buckets[3][2].SkillScore)

	assert.Equal(t, 4, buckets[0][3].SkillScore)
	assert.Equal(t, 3, buckets[1][3].SkillScore)
	assert.Equal(t, 2, buckets[2][3].SkillScore)
	assert.Equal(t, 1, buckets[3][3].SkillScore)
}
