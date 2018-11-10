package towerfall

import (
	"fmt"
	"log"
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

		// XXX: If we don't add the person to the database anything that tries to
		// grab from it will fail. Backfilling from the semis is one of those
		// cases. That should be refactored away and this should be removed.
		tm.db.SavePerson(p)
	}

	return
}

func endPlayoffs(t *Tournament) error {
	for x := 0; x < len(t.Matches)-3; x++ {
		if err := t.Matches[x].Start(nil); err != nil {
			return err
		}
		if err := t.Matches[x].End(nil); err != nil {
			return err
		}
	}
	return nil
}

func endSemis(t *Tournament) error {
	offset := len(t.Matches) - 3
	if err := t.Matches[offset].Start(nil); err != nil {
		return err
	}
	if err := t.Matches[offset].End(nil); err != nil {
		return err
	}
	if err := t.Matches[offset+1].Start(nil); err != nil {
		return err
	}
	if err := t.Matches[offset+1].End(nil); err != nil {
		return err
	}
	return nil
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
		t.Run(fmt.Sprintf("Match %d", x+1), func(t *testing.T) { runTestMatch(t, tm, x) })
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
		if x%3 == 0 {
			p := testPerson(s)
			s := NewPlayer(p).Summary()
			t.Logf("Adding player: %s", p.Nick)

			err := tm.AddPlayer(&s)
			assert.NoError(t, err)
		}

		t.Run(fmt.Sprintf("Match %d", x+1), func(t *testing.T) { runTestMatch(t, tm, x) })
	}
}

func TestStartingTournamentWithFewerThan8PlayersFail(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 7)
	err := tm.StartTournament(nil)
	assert.NotNil(err)
}

func TestStartingTournamentWith8PlayersWorks(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 8)
	err := tm.StartTournament(nil)
	assert.NoError(err)
}

func TestStartingTournamentWith24PlayersWorks(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 24)
	err := tm.StartTournament(nil)
	assert.NoError(err)
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

func TestStartingGivesTheRightAmountOfPlayoffs(t *testing.T) {
	assert := assert.New(t)
	for x := 8; x <= 32; x++ {
		t.Run(fmt.Sprintf("With%d", x), func(t *testing.T) {
			s, teardown := MockServer(t)
			defer teardown()

			tm := testTournament(t, s, x)
			err := tm.StartTournament(nil)
			assert.NoError(err)

			if x == 8 {
				// A special case - we don't need any playoffs since we're ready
				// for semi-finals right away.
				assert.Equal(3, len(tm.Matches))
				return
			}

			// The -1 is to shift so that when we have a player count
			// divisible by four an extra match isn't started. E.g. when we
			// have 16 players we want 4 matches, but without the -1 a fifth
			// match would be added.
			y := (x - 8 - 1) / 4
			m := len(tm.Matches) - 3
			compare := 3 + y
			assert.Equal(
				compare,
				m,
				fmt.Sprintf("%d player tournament had %d matches, not %d", x, m, compare),
			)
		})
	}
}

func TestStartingTournamentSetsStartedTimestamp(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 8)

	tm.StartTournament(nil)
	assert.NotNil(tm.Started)
}

func TestStartingTournamentCreatesTenEvents(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 8)

	tm.StartTournament(nil)
	assert.Equal(1+8+1, len(tm.Events))
	assert.Equal(tm.Events[0].Kind, "new_tournament")
	assert.Equal(tm.Events[1].Kind, "player_join")
	assert.Equal(tm.Events[9].Kind, "start")
}

func TestPopulateMatchesPopulatesPlayoffsFor8Players(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 8)
	tm.StartTournament(nil)

	assert.Equal(4, len(tm.Matches[0].Players))
	assert.Equal(4, len(tm.Matches[1].Players))
}

func TestPopulateMatchesPopulatesAllMatchesFor24Players(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 24)
	tm.StartTournament(nil)

	assert.Equal(4, len(tm.Matches[0].Players))
	assert.Equal(4, len(tm.Matches[1].Players))
	assert.Equal(4, len(tm.Matches[2].Players))
	assert.Equal(4, len(tm.Matches[3].Players))
	assert.Equal(4, len(tm.Matches[4].Players))
	assert.Equal(4, len(tm.Matches[5].Players))
}

func TestRunnerupInsertion(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 23)
	tm.StartTournament(nil)

	m, err := tm.NextMatch()
	assert.NoError(err)
	m.Start(nil)
	m.Players[0].AddKills(10)
	m.End(nil)

	assert.Equal(m.Players[1].Person.PersonID, tm.Runnerups[0].PersonID)
	assert.Equal(m.Players[2].Person.PersonID, tm.Runnerups[1].PersonID)
	assert.Equal(m.Players[3].Person.PersonID, tm.Runnerups[2].PersonID)
}

func TestNextMatchNoMatchesAreStartedWithPlayoffs(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	tm.StartTournament(nil)

	m, err := tm.NextMatch()
	assert.NoError(err)
	assert.Equal(0, m.Index)
	assert.Equal("playoff", m.Kind)

	m.Start(nil)
	m.End(nil)

	m, err = tm.NextMatch()
	assert.NoError(err)
	assert.Equal(1, m.Index)
	assert.Equal("playoff", m.Kind)
	assert.Equal(CurrentMatch(1), tm.Current+1)
}

func TestNextMatchNoMatchesAreStartedWithPlayoffsDone(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	tm.StartTournament(nil)
	err := endPlayoffs(tm)
	assert.NoError(err)

	m, err := tm.NextMatch()
	assert.NoError(err)
	assert.Equal(4, m.Index)
	assert.Equal("semi", m.Kind)
	assert.Equal(CurrentMatch(tm.MatchIndex(tm.Semi(0))), tm.Current+1)
}

func TestNextMatchNoMatchesAreStartedWithPlayoffsAndSemisDone(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	tm.StartTournament(nil)
	endPlayoffs(tm)
	endSemis(tm)

	m, err := tm.NextMatch()
	assert.NoError(err)
	assert.Equal(6, m.Index)
	assert.Equal("final", m.Kind)
	assert.Equal(CurrentMatch(tm.MatchIndex(tm.Final())), tm.Current+1)
}

func TestNextMatchEverythingDone(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	tm.StartTournament(nil)
	endPlayoffs(tm)
	endSemis(tm)
	tm.Final().Start(nil)
	tm.Final().End(nil)

	_, err := tm.NextMatch()
	assert.NotNil(err)
}

// func TestUpdatePlayer(t *testing.T) {
// 	assert := assert.New(t)
// 	s, teardown := MockServer(t)
// 	defer teardown()

// 	tm := testTournament(t, s, 8)
// 	tm.StartTournament(nil)
// 	m, err := tm.NextMatch()
// 	assert.NoError(err)

// 	p, err := tm.getTournamentPlayerObject(m.Players[3].Person)
// 	t.Logf("%+v", p)

// 	m.Start(nil)

// 	m.Players[0].AddKills(5)
// 	m.Players[1].AddKills(6)
// 	m.Players[2].AddKills(7)
// 	m.Players[3].AddKills(10)

// 	m.End(nil)

// 	p, err = tm.getTournamentPlayerObject(m.Players[3].Person)
// 	t.Logf("%+v", p)
// 	t.Log(len(tm.Matches))
// 	assert.NoError(err)
// 	assert.Equal(10, p.Kills)

// 	p, err = tm.getTournamentPlayerObject(m.Players[2].Person)
// 	t.Logf("%+v", p)
// 	assert.NoError(err)
// 	assert.Equal(7, p.Kills)

// 	p, err = tm.getTournamentPlayerObject(m.Players[1].Person)
// 	t.Logf("%+v", p)
// 	assert.NoError(err)
// 	assert.Equal(6, p.Kills)

// 	p, err = tm.getTournamentPlayerObject(m.Players[0].Person)
// 	t.Logf("%+v", p)
// 	assert.NoError(err)
// 	assert.Equal(5, p.Kills)
// }

func TestEnd4MatchPlayoffsPlacesWinnerAndSecondIntoSemisAndRestIntoRunnerups(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	tm.StartTournament(nil)
	m, err := tm.NextMatch()
	assert.NoError(err)

	m.Start(nil)

	m.Players[0].AddKills(5)
	m.Players[1].AddKills(6)
	m.Players[2].AddKills(7)
	m.Players[3].AddKills(10)
	winner := m.Players[3].Name()
	silver := m.Players[2].Name()

	m.End(nil)

	assert.Equal(1, len(tm.Semi(0).Players))
	assert.Equal(1, len(tm.Semi(1).Players))
	assert.Equal(2, len(tm.Runnerups))

	assert.Equal(winner, tm.Semi(0).Players[0].Name())
	assert.Equal(silver, tm.Semi(1).Players[0].Name())
}

func TestEndComplete16PlayerTournamentKillsOnly(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 16)
	tm.StartTournament(nil)

	// Playoff 1 (same as test above)
	m, err := tm.NextMatch()
	assert.NoError(err)

	m.Start(nil)

	m.Players[0].AddKills(5)
	m.Players[1].AddKills(6)
	m.Players[2].AddKills(7)
	m.Players[3].AddKills(10)
	winner := m.Players[3].Name()
	silver := m.Players[2].Name()

	m.End(nil)

	assert.Equal(1, len(tm.Semi(0).Players))
	assert.Equal(1, len(tm.Semi(1).Players))
	assert.Equal(2, len(tm.Runnerups))

	assert.Equal(winner, tm.Semi(0).Players[0].Name())
	assert.Equal(silver, tm.Semi(1).Players[0].Name())

	// Playoff 2
	m2, err2 := tm.NextMatch()
	assert.NoError(err2)

	err2 = m2.Start(nil)
	assert.NoError(err2)

	m2.Players[0].AddKills(2)
	m2.Players[1].AddKills(10)
	m2.Players[2].AddKills(8)
	m2.Players[3].AddKills(4)
	winner2 := m2.Players[1].Name()
	silver2 := m2.Players[2].Name()

	m2.End(nil)

	assert.Equal(2, len(tm.Semi(0).Players))
	assert.Equal(2, len(tm.Semi(1).Players))
	assert.Equal(4, len(tm.Runnerups))

	assert.Equal(winner2, tm.Semi(1).Players[1].Name())
	assert.Equal(silver2, tm.Semi(0).Players[1].Name())

	// Playoff 3
	m3, err3 := tm.NextMatch()
	assert.NoError(err3)

	m3.Start(nil)

	m3.Players[0].AddKills(10)
	m3.Players[1].AddKills(3)
	m3.Players[2].AddKills(3)
	m3.Players[3].AddKills(5)
	winner3 := m3.Players[0].Name()
	silver3 := m3.Players[3].Name()

	m3.End(nil)

	assert.Equal(3, len(tm.Semi(0).Players))
	assert.Equal(3, len(tm.Semi(1).Players))
	assert.Equal(6, len(tm.Runnerups))

	assert.Equal(winner3, tm.Semi(0).Players[2].Name())
	assert.Equal(silver3, tm.Semi(1).Players[2].Name())

	// Playoff 4
	m4, err4 := tm.NextMatch()
	assert.NoError(err4)

	m4.Start(nil)

	m4.Players[0].AddKills(9)
	m4.Players[1].AddKills(10)
	m4.Players[2].AddKills(5)
	m4.Players[3].AddKills(5)
	winner4 := m4.Players[1].Name()
	silver4 := m4.Players[0].Name()

	assert.NoError(m4.End(nil))

	assert.Equal(4, len(tm.Semi(0).Players))
	assert.Equal(4, len(tm.Semi(1).Players))
	assert.Equal(8, len(tm.Runnerups))

	assert.Equal(winner4, tm.Semi(1).Players[3].Name())
	assert.Equal(silver4, tm.Semi(0).Players[3].Name())

	// Semi 1
	s1, serr1 := tm.NextMatch()
	assert.NoError(serr1)

	assert.Equal("semi", s1.Kind)

	s1.Start(nil)

	s1.Players[0].AddKills(10)
	s1.Players[1].AddKills(7)
	s1.Players[2].AddKills(9)
	s1.Players[3].AddKills(8)
	winners1 := s1.Players[0].Name()
	silvers1 := s1.Players[2].Name()

	s1.End(nil)

	assert.Equal(2, len(tm.Final().Players))

	assert.Equal(winners1, tm.Final().Players[0].Name())
	assert.Equal(silvers1, tm.Final().Players[1].Name())

	// Semi 2
	s2, serr2 := tm.NextMatch()
	assert.NoError(serr2)

	assert.Equal("semi", s2.Kind)

	s2.Start(nil)

	s2.Players[0].AddKills(8)
	s2.Players[1].AddKills(10)
	s2.Players[2].AddKills(8)
	s2.Players[3].AddKills(9)
	winners2 := s2.Players[1].Name()
	silvers2 := s2.Players[3].Name()

	s2.End(nil)

	assert.Equal(4, len(tm.Final().Players))

	assert.Equal(winners2, tm.Final().Players[2].Name())
	assert.Equal(silvers2, tm.Final().Players[3].Name())

	// Final!
	f, ferr := tm.NextMatch()
	assert.NoError(ferr)

	assert.Equal("final", f.Kind)

	f.Start(nil)

	f.Players[0].AddKills(7)
	f.Players[1].AddKills(2)
	f.Players[2].AddKills(9)
	f.Players[3].AddKills(20)
	gold := f.Players[3].Name()
	lowe := f.Players[2].Name()
	bronze := f.Players[0].Name()

	f.End(nil)

	assert.Equal(gold, tm.Winners[0].Name())
	assert.Equal(lowe, tm.Winners[1].Name())
	assert.Equal(bronze, tm.Winners[2].Name())
}

// func TestEndComplete19PlayerTournamentKillsOnly(t *testing.T) {
// 	// This primarily tests the runnerup population for the fifth match
// 	// and that only the winners are propagated when there are more
// 	// than 16 players.
// 	assert := assert.New(t)

// 	s, teardown := MockServer(t)
// 	defer teardown()

// 	tm := testTournament(t, s, 19)
// 	tm.StartTournament(nil)

// 	// There should be 5 playoffs (and the predefineds)
// 	assert.Equal(5+3, len(tm.Matches))

// 	// Playoff 1
// 	m, err := tm.NextMatch()
// 	assert.NoError(err)

// 	m.Start(nil)

// 	m.Players[0].AddKills(5)
// 	m.Players[1].AddKills(6)
// 	m.Players[2].AddKills(7)
// 	m.Players[3].AddKills(10)
// 	winner := m.Players[3].Name()

// 	m.End(nil)

// 	assert.Equal(1, len(tm.Semi(0).Players))
// 	assert.Equal(0, len(tm.Semi(1).Players))
// 	assert.Equal(3, len(tm.Runnerups))

// 	assert.Equal(winner, tm.Semi(0).Players[0].Name())

// 	// Playoff 2
// 	m2, err2 := tm.NextMatch()
// 	assert.NoError(err2)

// 	m2.Start(nil)

// 	m2.Players[0].AddKills(2)
// 	m2.Players[1].AddKills(10)
// 	m2.Players[2].AddKills(8)
// 	m2.Players[3].AddKills(4)
// 	winner2 := m2.Players[1].Name()

// 	m2.End(nil)

// 	assert.Equal(1, len(tm.Semi(0).Players))
// 	assert.Equal(1, len(tm.Semi(1).Players))
// 	assert.Equal(6, len(tm.Runnerups))

// 	assert.Equal(winner2, tm.Semi(1).Players[0].Name())

// 	// Playoff 3
// 	m3, err3 := tm.NextMatch()
// 	assert.NoError(err3)

// 	m3.Start(nil)

// 	m3.Players[0].AddKills(10)
// 	m3.Players[1].AddKills(3)
// 	m3.Players[2].AddKills(3)
// 	m3.Players[3].AddKills(5)
// 	winner3 := m3.Players[0].Name()

// 	m3.End(nil)

// 	assert.Equal(2, len(tm.Semi(0).Players))
// 	assert.Equal(1, len(tm.Semi(1).Players))
// 	assert.Equal(9, len(tm.Runnerups))

// 	assert.Equal(winner3, tm.Semi(0).Players[1].Name())

// 	// Playoff 4
// 	m4, err4 := tm.NextMatch()
// 	assert.NoError(err4)

// 	m4.Start(nil)

// 	m4.Players[0].AddKills(9)
// 	m4.Players[1].AddKills(10)
// 	m4.Players[2].AddKills(5)
// 	m4.Players[3].AddKills(5)
// 	winner4 := m4.Players[1].Name()

// 	m4.End(nil)

// 	assert.Equal(2, len(tm.Semi(0).Players))
// 	assert.Equal(2, len(tm.Semi(1).Players))
// 	assert.Equal(12, len(tm.Runnerups))

// 	assert.Equal(winner4, tm.Semi(1).Players[1].Name())

// 	// Playoff 5 / Runnerup 1
// 	m5, err5 := tm.NextMatch()
// 	assert.NoError(err5)
// 	assert.Equal("playoff", m5.Kind)

// 	m5.Start(nil)
// 	// Given the 19 player match, there are 3 players that have yet to contend
// 	// and therefore we need to pick one of the runnerups.
// 	assert.Equal(4, len(m5.Players))
// 	assert.Equal(12, len(tm.Runnerups))

// 	m5.Players[0].AddKills(8)
// 	m5.Players[1].AddKills(7)
// 	m5.Players[2].AddKills(2)
// 	m5.Players[3].AddKills(10)
// 	winner5 := m5.Players[3].Name()

// 	m5.End(nil)

// 	assert.Equal(winner5, tm.Semi(0).Players[2].Name())

// 	// We need to backfill the players, and since that is a judge action we need
// 	// to simulate that
// 	err = tm.BackfillSemis(nil, []string{
// 		tm.Runnerups[0].PersonID,
// 		tm.Runnerups[1].PersonID,
// 		tm.Runnerups[2].PersonID,
// 	})

// 	assert.NoError(err)

// 	assert.Equal(4, len(tm.Semi(0).Players))
// 	assert.Equal(4, len(tm.Semi(1).Players))
// 	assert.Equal(11, len(tm.Runnerups))

// 	// Semi 1
// 	s1, serr1 := tm.NextMatch()
// 	assert.NoError(serr1)

// 	assert.Equal("semi", s1.Kind)

// 	s1.Start(nil)

// 	s1.Players[0].AddKills(10)
// 	s1.Players[1].AddKills(7)
// 	s1.Players[2].AddKills(9)
// 	s1.Players[3].AddKills(8)
// 	winners1 := s1.Players[0].Name()
// 	silvers1 := s1.Players[2].Name()

// 	s1.End(nil)

// 	assert.Equal(2, len(tm.Final().Players))

// 	assert.Equal(winners1, tm.Final().Players[0].Name())
// 	assert.Equal(silvers1, tm.Final().Players[1].Name())

// 	// Semi 2
// 	s2, serr2 := tm.NextMatch()
// 	assert.NoError(serr2)

// 	assert.Equal("semi", s2.Kind)

// 	s2.Start(nil)

// 	s2.Players[0].AddKills(8)
// 	s2.Players[1].AddKills(10)
// 	s2.Players[2].AddKills(8)
// 	s2.Players[3].AddKills(9)
// 	winners2 := s2.Players[1].Name()
// 	silvers2 := s2.Players[3].Name()

// 	s2.End(nil)

// 	assert.Equal(4, len(tm.Final().Players))

// 	assert.Equal(winners2, tm.Final().Players[2].Name())
// 	assert.Equal(silvers2, tm.Final().Players[3].Name())

// 	// Final!
// 	f, ferr := tm.NextMatch()
// 	assert.NoError(ferr)

// 	assert.Equal("final", f.Kind)

// 	f.Start(nil)

// 	f.Players[0].AddKills(7)
// 	f.Players[1].AddKills(2)
// 	f.Players[2].AddKills(9)
// 	f.Players[3].AddKills(10)

// 	gold := f.Players[3].Name()
// 	lowe := f.Players[2].Name()
// 	bronze := f.Players[0].Name()

// 	f.End(nil)

// 	assert.Equal(gold, tm.Winners[0].Name())
// 	assert.Equal(lowe, tm.Winners[1].Name())
// 	assert.Equal(bronze, tm.Winners[2].Name())
// }

func runTestMatch(t *testing.T, tm *Tournament, index int) {
	t.Run("Play", func(t *testing.T) {
		nm, err := tm.NextMatch()
		assert.Equal(t, index, nm.Index)
		assert.NoError(t, err)

		err = nm.Autoplay()
		assert.NoError(t, err)
	})

	// t.Run("Match counts for players are set", func(t *testing.T) {
	// 	m := 4 * (index + 1)
	// 	players := len(tm.Players)
	// 	minMatches := m / players
	// 	played := m % players

	// 	ps := []*PlayerSummary{}
	// 	err := tm.db.DB.Model(&ps).Where("tournament_id = ? AND matches = ?", tm.ID, minMatches+1).Select()
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, played, len(ps))
	// })

	t.Run("Runnerup order", func(t *testing.T) {
		rups, err := tm.GetRunnerups()
		assert.NoError(t, err)
		assert.Equal(t, 4, len(rups))
	})
}
