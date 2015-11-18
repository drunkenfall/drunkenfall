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
	assertMatches(a, 17, 6, 2)
	assertMatches(a, 18, 6, 2)
	assertMatches(a, 19, 6, 2)
	assertMatches(a, 20, 6, 2)
	assertMatches(a, 21, 6, 2)
	assertMatches(a, 22, 6, 2)
	assertMatches(a, 23, 6, 2)
	assertMatches(a, 24, 6, 2)
}
