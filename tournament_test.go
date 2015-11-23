package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
)

// testTournament makes a test tournament with `count` players.
func testTournament(count int) (t *Tournament) {
	t, err := NewTournament()
	if err != nil {
		log.Fatal("tournament creation failed")
	}

	for i := 1; i <= count; i++ {
		name := strconv.Itoa(i)
		p := Player{name: name}
		t.Players[name] = p
	}
	return
}

func endTryouts(t *Tournament) {
	for x := range t.Tryouts {
		t.Tryouts[x].Start()
		t.Tryouts[x].End()
	}
}

func endSemis(t *Tournament) {
	for x := range t.Semis {
		t.Semis[x].Start()
		t.Semis[x].End()
	}
}

func TestStartingTournamentWithFewerThan8PlayersFail(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(7)
	err := tm.StartTournament()
	assert.NotNil(err)
}
func TestStartingTournamentWith8PlayersWorks(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(8)
	err := tm.StartTournament()
	assert.Nil(err)
}

func TestStartingTournamentWithMoreThan24PlayersFail(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(25)
	err := tm.StartTournament()
	assert.NotNil(err)
}

func TestStartingTournamentWith24PlayersWorks(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(24)
	err := tm.StartTournament()
	assert.Nil(err)
}

func TestStartingTournamentSetsStartedTimestamp(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(8)

	tm.StartTournament()
	assert.NotNil(tm.Started)
}

func assertMatches(assert *assert.Assertions, players, matches, semis int) {
	tm := testTournament(players)
	tm.StartTournament()

	assert.Equal(len(tm.Tryouts), matches)
	assert.Equal(len(tm.Semis), semis)
}

func TestStartingTournamentParticipantsToMatches(t *testing.T) {
	// Assert that a given number of participants results in a number of matches.
	// See docstring for GenerateMatches()

	a := assert.New(t)
	assertMatches(a, 8, 0, 2)
	assertMatches(a, 9, 4, 2)
	assertMatches(a, 10, 4, 2)
	assertMatches(a, 11, 4, 2)
	assertMatches(a, 12, 4, 2)
	assertMatches(a, 13, 4, 2)
	assertMatches(a, 14, 4, 2)
	assertMatches(a, 15, 4, 2)
	assertMatches(a, 16, 4, 2)
	assertMatches(a, 17, 7, 2)
	assertMatches(a, 18, 7, 2)
	assertMatches(a, 19, 7, 2)
	assertMatches(a, 20, 7, 2)
	assertMatches(a, 21, 7, 2)
	assertMatches(a, 22, 7, 2)
	assertMatches(a, 23, 7, 2)
	assertMatches(a, 24, 7, 2)
}

func TestPopulateMatchesPopulatesSemisFor8Players(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(8)
	tm.StartTournament()

	assert.Equal(4, len(tm.Semis[0].Players))
	assert.Equal(4, len(tm.Semis[1].Players))
}

func TestPopulateMatchesPopulatesAllMatchesFor24Players(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(24)
	tm.StartTournament()

	assert.Equal(4, len(tm.Tryouts[0].Players))
	assert.Equal(4, len(tm.Tryouts[1].Players))
	assert.Equal(4, len(tm.Tryouts[2].Players))
	assert.Equal(4, len(tm.Tryouts[3].Players))
	assert.Equal(4, len(tm.Tryouts[4].Players))
	assert.Equal(4, len(tm.Tryouts[5].Players))
}

func TestPopulateMatchesFillsAsMuchAsPossibleFor10Players(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(10)
	tm.StartTournament()

	assert.Equal(4, len(tm.Tryouts[0].Players))
	assert.Equal(4, len(tm.Tryouts[1].Players))
	assert.Equal(2, len(tm.Tryouts[2].Players))
}

func TestPopulateMatchesFillsAsMuchAsPossibleFor18Players(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(18)
	tm.StartTournament()

	assert.Equal(4, len(tm.Tryouts[0].Players))
	assert.Equal(4, len(tm.Tryouts[1].Players))
	assert.Equal(4, len(tm.Tryouts[2].Players))
	assert.Equal(4, len(tm.Tryouts[3].Players))
	assert.Equal(2, len(tm.Tryouts[4].Players))
	assert.Equal(0, len(tm.Tryouts[5].Players))
}

func TestNextMatchNoMatchesAreStartedWithTryouts(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(16)
	tm.StartTournament()

	m, err := tm.NextMatch()
	assert.Nil(err)
	assert.Equal(0, m.Index)
	assert.Equal("tryout", m.Kind)

	m.Start()
	m.End()

	m, err = tm.NextMatch()
	assert.Nil(err)
	assert.Equal(1, m.Index)
	assert.Equal("tryout", m.Kind)
}

func TestNextMatchNoMatchesAreStartedWithTryoutsDone(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(16)
	tm.StartTournament()
	endTryouts(tm)

	m, err := tm.NextMatch()
	assert.Nil(err)
	assert.Equal(0, m.Index)
	assert.Equal("semi", m.Kind)
}

func TestNextMatchNoMatchesAreStartedWithTryoutsAndSemisDone(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(16)
	tm.StartTournament()
	endTryouts(tm)
	endSemis(tm)

	m, err := tm.NextMatch()
	assert.Nil(err)
	assert.Equal(0, m.Index)
	assert.Equal("final", m.Kind)
}

func TestNextMatchEverythingDone(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(16)
	tm.StartTournament()
	endTryouts(tm)
	endSemis(tm)
	tm.Final.Start()
	tm.Final.End()

	_, err := tm.NextMatch()
	assert.NotNil(err)
}

func TestEnd4MatchTryoutsPlacesWinnerAndSecondIntoSemisAndRestIntoRunnerups(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(16)
	tm.StartTournament()
	m, err := tm.NextMatch()
	assert.Nil(err)

	m.Start()

	m.Players[0].AddKill(5)
	m.Players[1].AddKill(6)
	m.Players[2].AddKill(7)
	m.Players[3].AddKill(10)
	winner := m.Players[3].name
	silver := m.Players[2].name

	m.End()

	assert.Equal(1, len(tm.Semis[0].Players))
	assert.Equal(1, len(tm.Semis[1].Players))
	assert.Equal(2, len(tm.Runnerups))

	assert.Equal(winner, tm.Semis[0].Players[0].name)
	assert.Equal(silver, tm.Semis[1].Players[0].name)
}

func TestEndComplete16PlayerTournament(t *testing.T) {
	assert := assert.New(t)

	tm := testTournament(16)
	tm.StartTournament()

	// Tryout 1 (same as test above)
	m, err := tm.NextMatch()
	assert.Nil(err)

	m.Start()

	m.Players[0].AddKill(5)
	m.Players[1].AddKill(6)
	m.Players[2].AddKill(7)
	m.Players[3].AddKill(10)
	winner := m.Players[3].name
	silver := m.Players[2].name

	m.End()

	assert.Equal(1, len(tm.Semis[0].Players))
	assert.Equal(1, len(tm.Semis[1].Players))
	assert.Equal(2, len(tm.Runnerups))

	assert.Equal(winner, tm.Semis[0].Players[0].name)
	assert.Equal(silver, tm.Semis[1].Players[0].name)

	// Tryout 2
	m2, err2 := tm.NextMatch()
	assert.Nil(err2)

	m2.Start()

	m2.Players[0].AddKill(2)
	m2.Players[1].AddKill(10)
	m2.Players[2].AddKill(8)
	m2.Players[3].AddKill(4)
	winner2 := m2.Players[1].name
	silver2 := m2.Players[2].name

	m2.End()

	assert.Equal(2, len(tm.Semis[0].Players))
	assert.Equal(2, len(tm.Semis[1].Players))
	assert.Equal(4, len(tm.Runnerups))

	assert.Equal(winner2, tm.Semis[1].Players[1].name)
	assert.Equal(silver2, tm.Semis[0].Players[1].name)

	// Tryout 3
	m3, err3 := tm.NextMatch()
	assert.Nil(err3)

	m3.Start()

	m3.Players[0].AddKill(10)
	m3.Players[1].AddKill(3)
	m3.Players[2].AddKill(3)
	m3.Players[3].AddKill(5)
	winner3 := m3.Players[0].name
	silver3 := m3.Players[3].name

	m3.End()

	assert.Equal(3, len(tm.Semis[0].Players))
	assert.Equal(3, len(tm.Semis[1].Players))
	assert.Equal(6, len(tm.Runnerups))

	assert.Equal(winner3, tm.Semis[0].Players[2].name)
	assert.Equal(silver3, tm.Semis[1].Players[2].name)

	// Tryout 4
	m4, err4 := tm.NextMatch()
	assert.Nil(err4)

	m4.Start()

	m4.Players[0].AddKill(9)
	m4.Players[1].AddKill(10)
	m4.Players[2].AddKill(5)
	m4.Players[3].AddKill(5)
	winner4 := m4.Players[1].name
	silver4 := m4.Players[0].name

	m4.End()

	assert.Equal(4, len(tm.Semis[0].Players))
	assert.Equal(4, len(tm.Semis[1].Players))
	assert.Equal(8, len(tm.Runnerups))

	assert.Equal(winner4, tm.Semis[1].Players[3].name)
	assert.Equal(silver4, tm.Semis[0].Players[3].name)

	// Semi 1
	s1, serr1 := tm.NextMatch()
	assert.Nil(serr1)

	assert.Equal("semi", s1.Kind)

	s1.Start()

	s1.Players[0].AddKill(10)
	s1.Players[1].AddKill(7)
	s1.Players[2].AddKill(9)
	s1.Players[3].AddKill(8)
	winners1 := s1.Players[0].name
	silvers1 := s1.Players[2].name

	s1.End()

	assert.Equal(2, len(tm.Final.Players))

	assert.Equal(winners1, tm.Final.Players[0].name)
	assert.Equal(silvers1, tm.Final.Players[1].name)

	// Semi 2
	s2, serr2 := tm.NextMatch()
	assert.Nil(serr2)

	assert.Equal("semi", s2.Kind)

	s2.Start()

	s2.Players[0].AddKill(8)
	s2.Players[1].AddKill(10)
	s2.Players[2].AddKill(8)
	s2.Players[3].AddKill(9)
	winners2 := s2.Players[1].name
	silvers2 := s2.Players[3].name

	s2.End()

	assert.Equal(4, len(tm.Final.Players))

	assert.Equal(winners2, tm.Final.Players[2].name)
	assert.Equal(silvers2, tm.Final.Players[3].name)

	// Final!
	f, ferr := tm.NextMatch()
	assert.Nil(ferr)

	assert.Equal("final", f.Kind)

	f.Start()

	f.Players[0].AddKill(7)
	f.Players[1].AddKill(2)
	f.Players[2].AddKill(9)
	f.Players[3].AddKill(10)
	gold := f.Players[3].name
	lowe := f.Players[2].name
	bronze := f.Players[0].name

	f.End()

	t.Log(tm.Winners)
	assert.Equal(gold, tm.Winners[0].name)
	assert.Equal(lowe, tm.Winners[1].name)
	assert.Equal(bronze, tm.Winners[2].name)
}
