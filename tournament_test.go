package main

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
)

// testTournament makes a test tournament with `count` players.
func testTournament(count int) (t *Tournament) {
	s := strconv.Itoa(count)
	db := MockDatabase()
	t, err := NewTournament("Tournament "+s, s, db)
	if err != nil {
		log.Fatal("tournament creation failed")
	}

	for i := 1; i <= count; i++ {
		name := strconv.Itoa(i)
		t.AddPlayer(name)
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

// func TestStartingTournamentWithMoreThan32PlayersFail(t *testing.T) {
// 	assert := assert.New(t)
// 	tm := testTournament(33)
// 	err := tm.StartTournament()
// 	assert.NotNil(err)
// }

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

func TestPopulateMatchesPopulatesTryoutsFor8Players(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(8)
	tm.StartTournament()

	assert.Equal(4, len(tm.Tryouts[0].Players))
	assert.Equal(4, len(tm.Tryouts[1].Players))
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
	assert.Equal(4, len(tm.Tryouts[2].Players))

	assert.Equal(4, tm.Tryouts[0].ActualPlayers())
	assert.Equal(4, tm.Tryouts[1].ActualPlayers())
	assert.Equal(2, tm.Tryouts[2].ActualPlayers())
}

func TestPopulateMatchesFillsAsMuchAsPossibleFor18Players(t *testing.T) {
	assert := assert.New(t)
	tm := testTournament(18)
	tm.StartTournament()

	assert.Equal(4, len(tm.Tryouts[0].Players))
	assert.Equal(4, len(tm.Tryouts[1].Players))
	assert.Equal(4, len(tm.Tryouts[2].Players))
	assert.Equal(4, len(tm.Tryouts[3].Players))
	assert.Equal(4, len(tm.Tryouts[4].Players))
	assert.Equal(4, len(tm.Tryouts[5].Players))

	assert.Equal(4, tm.Tryouts[0].ActualPlayers())
	assert.Equal(4, tm.Tryouts[1].ActualPlayers())
	assert.Equal(4, tm.Tryouts[2].ActualPlayers())
	assert.Equal(4, tm.Tryouts[3].ActualPlayers())
	assert.Equal(2, tm.Tryouts[4].ActualPlayers())
	assert.Equal(0, tm.Tryouts[5].ActualPlayers())
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

func TestUpdatePlayer(t *testing.T) {
	assert := assert.New(t)
	tm, _ := NewTournament("player test", "test", nil)
	tm.AddPlayer("winner")
	tm.AddPlayer("loser1")
	tm.AddPlayer("loser2")
	tm.AddPlayer("loser3")

	tm.Tryouts = []*Match{
		{
			Kind: "tryout",
			Players: []Player{
				Player{Name: "winner", Kills: 10},
				Player{Name: "loser1", Kills: 0, Shots: 1},
				Player{Name: "loser2", Kills: 0},
				Player{Name: "loser3", Kills: 0},
			},
		},
		{
			Kind: "tryout",
			Players: []Player{
				Player{Name: "winner", Kills: 10},
				Player{Name: "loser1", Kills: 0},
				Player{Name: "loser2", Kills: 0},
				Player{Name: "loser3", Kills: 0},
			},
		},
	}
	tm.Semis = []*Match{}
	tm.Final = &Match{}

	tm.UpdatePlayers()
	assert.Equal(20, tm.getPlayer("winner").Kills)
	assert.Equal(1, tm.getPlayer("loser1").Shots)
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
	winner := m.Players[3].Name
	silver := m.Players[2].Name

	m.End()

	assert.Equal(1, tm.Semis[0].ActualPlayers())
	assert.Equal(1, tm.Semis[1].ActualPlayers())
	assert.Equal(2, len(tm.Runnerups))

	assert.Equal(winner, tm.Semis[0].Players[0].Name)
	assert.Equal(silver, tm.Semis[1].Players[0].Name)
}

func TestEndComplete16PlayerTournamentKillsOnly(t *testing.T) {
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
	winner := m.Players[3].Name
	silver := m.Players[2].Name

	m.End()

	assert.Equal(1, tm.Semis[0].ActualPlayers())
	assert.Equal(1, tm.Semis[1].ActualPlayers())
	assert.Equal(2, len(tm.Runnerups))

	assert.Equal(winner, tm.Semis[0].Players[0].Name)
	assert.Equal(silver, tm.Semis[1].Players[0].Name)

	// Tryout 2
	m2, err2 := tm.NextMatch()
	assert.Nil(err2)

	m2.Start()

	m2.Players[0].AddKill(2)
	m2.Players[1].AddKill(10)
	m2.Players[2].AddKill(8)
	m2.Players[3].AddKill(4)
	winner2 := m2.Players[1].Name
	silver2 := m2.Players[2].Name

	m2.End()

	assert.Equal(2, tm.Semis[0].ActualPlayers())
	assert.Equal(2, tm.Semis[1].ActualPlayers())
	assert.Equal(4, len(tm.Runnerups))

	assert.Equal(winner2, tm.Semis[1].Players[1].Name)
	assert.Equal(silver2, tm.Semis[0].Players[1].Name)

	// Tryout 3
	m3, err3 := tm.NextMatch()
	assert.Nil(err3)

	m3.Start()

	m3.Players[0].AddKill(10)
	m3.Players[1].AddKill(3)
	m3.Players[2].AddKill(3)
	m3.Players[3].AddKill(5)
	winner3 := m3.Players[0].Name
	silver3 := m3.Players[3].Name

	m3.End()

	assert.Equal(3, tm.Semis[0].ActualPlayers())
	assert.Equal(3, tm.Semis[1].ActualPlayers())
	assert.Equal(6, len(tm.Runnerups))

	assert.Equal(winner3, tm.Semis[0].Players[2].Name)
	assert.Equal(silver3, tm.Semis[1].Players[2].Name)

	// Tryout 4
	m4, err4 := tm.NextMatch()
	assert.Nil(err4)

	m4.Start()

	m4.Players[0].AddKill(9)
	m4.Players[1].AddKill(10)
	m4.Players[2].AddKill(5)
	m4.Players[3].AddKill(5)
	winner4 := m4.Players[1].Name
	silver4 := m4.Players[0].Name

	m4.End()

	assert.Equal(4, tm.Semis[0].ActualPlayers())
	assert.Equal(4, tm.Semis[1].ActualPlayers())
	assert.Equal(8, len(tm.Runnerups))

	assert.Equal(winner4, tm.Semis[1].Players[3].Name)
	assert.Equal(silver4, tm.Semis[0].Players[3].Name)

	// Semi 1
	s1, serr1 := tm.NextMatch()
	assert.Nil(serr1)

	assert.Equal("semi", s1.Kind)

	s1.Start()

	s1.Players[0].AddKill(10)
	s1.Players[1].AddKill(7)
	s1.Players[2].AddKill(9)
	s1.Players[3].AddKill(8)
	winners1 := s1.Players[0].Name
	silvers1 := s1.Players[2].Name

	s1.End()

	assert.Equal(2, tm.Final.ActualPlayers())

	assert.Equal(winners1, tm.Final.Players[0].Name)
	assert.Equal(silvers1, tm.Final.Players[1].Name)

	// Semi 2
	s2, serr2 := tm.NextMatch()
	assert.Nil(serr2)

	assert.Equal("semi", s2.Kind)

	s2.Start()

	s2.Players[0].AddKill(8)
	s2.Players[1].AddKill(10)
	s2.Players[2].AddKill(8)
	s2.Players[3].AddKill(9)
	winners2 := s2.Players[1].Name
	silvers2 := s2.Players[3].Name

	s2.End()

	assert.Equal(4, tm.Final.ActualPlayers())

	assert.Equal(winners2, tm.Final.Players[2].Name)
	assert.Equal(silvers2, tm.Final.Players[3].Name)

	// Final!
	f, ferr := tm.NextMatch()
	assert.Nil(ferr)

	assert.Equal("final", f.Kind)

	f.Start()

	f.Players[0].AddKill(7)
	f.Players[1].AddKill(2)
	f.Players[2].AddKill(9)
	f.Players[3].AddKill(10)
	gold := f.Players[3].Name
	lowe := f.Players[2].Name
	bronze := f.Players[0].Name

	f.End()

	assert.Equal(gold, tm.Winners[0].Name)
	assert.Equal(lowe, tm.Winners[1].Name)
	assert.Equal(bronze, tm.Winners[2].Name)
}

func TestEndComplete19PlayerTournamentKillsOnly(t *testing.T) {
	// This primarily tests the runnerup population for the fifth match
	// and that only the winners are propagated when there are more
	// than 16 players.
	assert := assert.New(t)

	tm := testTournament(19)
	tm.StartTournament()

	assert.Equal(8, len(tm.Tryouts))

	// Tryout 1
	m, err := tm.NextMatch()
	assert.Nil(err)

	m.Start()

	m.Players[0].AddKill(5)
	m.Players[1].AddKill(6)
	m.Players[2].AddKill(7)
	m.Players[3].AddKill(10)
	winner := m.Players[3].Name

	m.End()

	assert.Equal(1, tm.Semis[0].ActualPlayers())
	assert.Equal(0, tm.Semis[1].ActualPlayers())
	assert.Equal(3, len(tm.Runnerups))

	assert.Equal(winner, tm.Semis[0].Players[0].Name)

	// Tryout 2
	m2, err2 := tm.NextMatch()
	assert.Nil(err2)

	m2.Start()

	m2.Players[0].AddKill(2)
	m2.Players[1].AddKill(10)
	m2.Players[2].AddKill(8)
	m2.Players[3].AddKill(4)
	winner2 := m2.Players[1].Name

	m2.End()

	assert.Equal(1, tm.Semis[0].ActualPlayers())
	assert.Equal(1, tm.Semis[1].ActualPlayers())
	assert.Equal(6, len(tm.Runnerups))

	assert.Equal(winner2, tm.Semis[1].Players[0].Name)

	// Tryout 3
	m3, err3 := tm.NextMatch()
	assert.Nil(err3)

	m3.Start()

	m3.Players[0].AddKill(10)
	m3.Players[1].AddKill(3)
	m3.Players[2].AddKill(3)
	m3.Players[3].AddKill(5)
	winner3 := m3.Players[0].Name

	m3.End()

	assert.Equal(2, tm.Semis[0].ActualPlayers())
	assert.Equal(1, tm.Semis[1].ActualPlayers())
	assert.Equal(9, len(tm.Runnerups))

	assert.Equal(winner3, tm.Semis[0].Players[1].Name)

	// Tryout 4
	m4, err4 := tm.NextMatch()
	assert.Nil(err4)

	m4.Start()

	m4.Players[0].AddKill(9)
	m4.Players[1].AddKill(10)
	m4.Players[2].AddKill(5)
	m4.Players[3].AddKill(5)
	winner4 := m4.Players[1].Name

	m4.End()

	assert.Equal(2, tm.Semis[0].ActualPlayers())
	assert.Equal(2, tm.Semis[1].ActualPlayers())
	assert.Equal(12, len(tm.Runnerups))

	assert.Equal(winner4, tm.Semis[1].Players[1].Name)

	// Tryout 5 / Runnerup 1
	m5, err5 := tm.NextMatch()
	assert.Nil(err5)
	assert.Equal("tryout", m5.Kind)

	m5.Start()
	// Given the 19 player match, there are 3 players that have yet to contend
	// and therefore we need to pick one of the runnerups.
	assert.Equal(4, m5.ActualPlayers())
	assert.Equal(12, len(tm.Runnerups))

	m5.Players[0].AddKill(8)
	m5.Players[1].AddKill(7)
	m5.Players[2].AddKill(2)
	m5.Players[3].AddKill(10)
	winner5 := m5.Players[3].Name

	m5.End()

	assert.Equal(3, tm.Semis[0].ActualPlayers())
	assert.Equal(2, tm.Semis[1].ActualPlayers())
	assert.Equal(14, len(tm.Runnerups))

	assert.Equal(winner5, tm.Semis[0].Players[2].Name)

	// Tryout 6 / Runnerup 2
	m6, err6 := tm.NextMatch()
	assert.Nil(err6)
	assert.Equal("tryout", m6.Kind)
	assert.Equal("tryout", m6.Kind)

	m6.Start()
	// Given the 19 player match, there are no new players.
	// As such, the backfill is completely from the runnerups.
	assert.Equal(4, m6.ActualPlayers())
	assert.Equal(14, len(tm.Runnerups))

	m6.Players[0].AddKill(10)
	m6.Players[1].AddKill(3)
	m6.Players[2].AddKill(1)
	m6.Players[3].AddKill(8)
	winner6 := m6.Players[0].Name

	m6.End()

	assert.Equal(3, tm.Semis[0].ActualPlayers())
	assert.Equal(3, tm.Semis[1].ActualPlayers())
	assert.Equal(13, len(tm.Runnerups))

	assert.Equal(winner6, tm.Semis[1].Players[2].Name)

	// Tryout 7 / Runnerup 3
	m7, err7 := tm.NextMatch()
	assert.Nil(err7)
	assert.Equal("tryout", m7.Kind)

	m7.Start()

	assert.Equal(4, m7.ActualPlayers())
	assert.Equal(13, len(tm.Runnerups))

	m7.Players[0].AddKill(7)
	m7.Players[1].AddKill(9)
	m7.Players[2].AddKill(10)
	m7.Players[3].AddKill(8)
	winner7 := m7.Players[2].Name

	m7.End()

	assert.Equal(4, tm.Semis[0].ActualPlayers())
	assert.Equal(3, tm.Semis[1].ActualPlayers())
	assert.Equal(12, len(tm.Runnerups))

	assert.Equal(winner7, tm.Semis[0].Players[3].Name)

	// Tryout 8 / Runnerup 4
	m8, err8 := tm.NextMatch()
	assert.Nil(err8)
	assert.Equal("tryout", m8.Kind)

	m8.Start()

	assert.Equal(4, m8.ActualPlayers())
	assert.Equal(12, len(tm.Runnerups))

	m8.Players[0].AddKill(10)
	m8.Players[1].AddKill(5)
	m8.Players[2].AddKill(6)
	m8.Players[3].AddKill(9)
	winner8 := m8.Players[0].Name

	m8.End()

	assert.Equal(4, tm.Semis[0].ActualPlayers())
	assert.Equal(4, tm.Semis[1].ActualPlayers())
	assert.Equal(11, len(tm.Runnerups))

	assert.Equal(winner8, tm.Semis[1].Players[3].Name)

	// Semi 1
	s1, serr1 := tm.NextMatch()
	assert.Nil(serr1)

	assert.Equal("semi", s1.Kind)

	s1.Start()

	s1.Players[0].AddKill(10)
	s1.Players[1].AddKill(7)
	s1.Players[2].AddKill(9)
	s1.Players[3].AddKill(8)
	winners1 := s1.Players[0].Name
	silvers1 := s1.Players[2].Name

	s1.End()

	assert.Equal(2, tm.Final.ActualPlayers())

	assert.Equal(winners1, tm.Final.Players[0].Name)
	assert.Equal(silvers1, tm.Final.Players[1].Name)

	// Semi 2
	s2, serr2 := tm.NextMatch()
	assert.Nil(serr2)

	assert.Equal("semi", s2.Kind)

	s2.Start()

	s2.Players[0].AddKill(8)
	s2.Players[1].AddKill(10)
	s2.Players[2].AddKill(8)
	s2.Players[3].AddKill(9)
	winners2 := s2.Players[1].Name
	silvers2 := s2.Players[3].Name

	s2.End()

	assert.Equal(4, tm.Final.ActualPlayers())

	assert.Equal(winners2, tm.Final.Players[2].Name)
	assert.Equal(silvers2, tm.Final.Players[3].Name)

	// Final!
	f, ferr := tm.NextMatch()
	assert.Nil(ferr)

	assert.Equal("final", f.Kind)

	f.Start()

	f.Players[0].AddKill(7)
	f.Players[1].AddKill(2)
	f.Players[2].AddKill(9)
	f.Players[3].AddKill(10)
	gold := f.Players[3].Name
	lowe := f.Players[2].Name
	bronze := f.Players[0].Name

	f.End()

	assert.Equal(gold, tm.Winners[0].Name)
	assert.Equal(lowe, tm.Winners[1].Name)
	assert.Equal(bronze, tm.Winners[2].Name)
}
