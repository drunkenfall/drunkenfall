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
		p := Player{name: strconv.Itoa(i)}
		t.Players = append(t.Players, p)
	}
	return
}

func endTryouts(t *Tournament) {
	for x := range t.Tryouts {
		t.Tryouts[x].StartMatch()
		t.Tryouts[x].EndMatch()
	}
}

func endSemis(t *Tournament) {
	for x := range t.Semis {
		t.Semis[x].StartMatch()
		t.Semis[x].EndMatch()
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
	tm.Final.StartMatch()
	tm.Final.EndMatch()

	_, err := tm.NextMatch()
	assert.NotNil(err)
}
