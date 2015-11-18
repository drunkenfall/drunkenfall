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

func assertParticipants(assert *assert.Assertions, players, matches int) {
	tm := testTournament(players)
	tm.StartTournament()

	assert.Equal(len(tm.Tryouts), matches)
}

func TestStartingTournamentParticipantsToMatches(t *testing.T) {
	// Assert that a given number of participants results in a number of matches.
	// See docstring for GenerateTryouts()

	a := assert.New(t)
	assertParticipants(a, 8, 2)
	assertParticipants(a, 9, 4)
	assertParticipants(a, 10, 4)
	assertParticipants(a, 11, 4)
	assertParticipants(a, 12, 4)
	assertParticipants(a, 13, 4)
	assertParticipants(a, 14, 4)
	assertParticipants(a, 15, 4)
	assertParticipants(a, 16, 4)
	assertParticipants(a, 17, 6)
	assertParticipants(a, 18, 6)
	assertParticipants(a, 19, 6)
	assertParticipants(a, 20, 6)
	assertParticipants(a, 21, 6)
	assertParticipants(a, 22, 6)
	assertParticipants(a, 23, 6)
	assertParticipants(a, 24, 6)
}
