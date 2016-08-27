package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// MockMatch makes a mock Match{} with a dummy Tournament{}
func MockMatch(idx int, cat string) *Match {
	s := MockServer()
	tm, _ := NewTournament("test", "t", s)
	tm.SetMatchPointers()
	return NewMatch(tm, idx, cat)
}

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "test")

	assert.Equal(4, len(m.Players))
	assert.Equal(0, m.ActualPlayers())

	p := Player{Name: "I exist"}

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, m.ActualPlayers())
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "test")

	m.Players = []Player{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
		{Name: "d"},
	}
	p := Player{}

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "test")
	m.Started = time.Now()

	err := m.Start()
	assert.NotNil(err)
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "test")
	m.Players = []Player{
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
	}

	err := m.Start()
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndGivesShotToWinner(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "tryout")
	// TODO(thiderman): This is terrible, but it works for now :(
	m.Tournament = nil
	m.Players = []Player{
		{Name: "1"},
		{Name: "2"},
		{Name: "3"},
		{Name: "4"},
	}

	err := m.Start()
	assert.Nil(err)
	m.Players[2].AddKill(10)

	err = m.End()
	assert.Nil(err)
	assert.Equal(1, m.Players[0].Shots)
}

func TestEndAlreadyEndedMatch(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "test")
	m.Ended = time.Now()

	err := m.End()
	assert.NotNil(err)
}

func TestEnd(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "test")

	err := m.End()
	assert.Nil(err)
	assert.Equal(false, m.Ended.IsZero())
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "test")
	_ = m.AddPlayer(Player{Name: "1"})
	_ = m.AddPlayer(Player{Name: "2"})
	_ = m.AddPlayer(Player{Name: "3"})
	_ = m.AddPlayer(Player{Name: "4"})
	ret := m.String()

	assert.Equal("<Test 1: 1 / 2 / 3 / 4 - not started>", ret)

	m2 := MockMatch(0, "final")
	_ = m2.AddPlayer(Player{Name: "a"})
	_ = m2.AddPlayer(Player{Name: "b"})
	_ = m2.AddPlayer(Player{Name: "c"})
	_ = m2.AddPlayer(Player{Name: "d"})
	_ = m2.Start()
	ret2 := m2.String()

	assert.Equal("<Final: a / b / c / d - playing>", ret2)
}

func TestCommitSweepPlayer1(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "test")
	_ = m.AddPlayer(Player{Name: "1"})
	_ = m.AddPlayer(Player{Name: "2"})
	_ = m.AddPlayer(Player{Name: "3"})
	_ = m.AddPlayer(Player{Name: "4"})

	scores := [][]int{
		[]int{3, 0},
		[]int{0, 0},
		[]int{0, 0},
		[]int{0, 0},
	}
	shots := []bool{
		false,
		false,
		false,
		false,
	}

	m.Commit(scores, shots)
	assert.Equal(1, m.Players[0].Sweeps)
}

func TestCommitDoubleKillPlayer2(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "test")
	_ = m.AddPlayer(Player{Name: "1"})
	_ = m.AddPlayer(Player{Name: "2"})
	_ = m.AddPlayer(Player{Name: "3"})
	_ = m.AddPlayer(Player{Name: "4"})

	scores := [][]int{
		[]int{0, 0},
		[]int{2, 0},
		[]int{0, 0},
		[]int{0, 0},
	}
	shots := []bool{
		false,
		false,
		false,
		false,
	}

	m.Commit(scores, shots)
	assert.Equal(2, m.Players[1].Kills)
}

func TestCommitSweepAndSuicidePlayer3(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "test")
	_ = m.AddPlayer(Player{Name: "1"})
	_ = m.AddPlayer(Player{Name: "2"})
	_ = m.AddPlayer(Player{Name: "3"})
	_ = m.AddPlayer(Player{Name: "4"})

	scores := [][]int{
		[]int{0, 0},
		[]int{0, 0},
		[]int{3, 1},
		[]int{0, 0},
	}
	shots := []bool{
		false,
		false,
		false,
		false,
	}

	m.Commit(scores, shots)
	assert.Equal(1, m.Players[2].Sweeps)
	assert.Equal(2, m.Players[2].Kills)
	assert.Equal(1, m.Players[2].Shots)
}

func TestCommitSuicidePlayer4(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "test")
	_ = m.AddPlayer(Player{Name: "1"})
	_ = m.AddPlayer(Player{Name: "2"})
	_ = m.AddPlayer(Player{Name: "3"})
	_ = m.AddPlayer(Player{Name: "4"})

	scores := [][]int{
		[]int{0, 0},
		[]int{0, 0},
		[]int{0, 0},
		[]int{0, 1},
	}
	shots := []bool{
		false,
		false,
		false,
		false,
	}

	m.Commit(scores, shots)
	assert.Equal(1, m.Players[3].Self)
	assert.Equal(1, m.Players[3].Shots)
}

func TestCommitShotsForPlayer2and3(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "test")
	_ = m.AddPlayer(Player{Name: "1"})
	_ = m.AddPlayer(Player{Name: "2"})
	_ = m.AddPlayer(Player{Name: "3"})
	_ = m.AddPlayer(Player{Name: "4"})

	scores := [][]int{
		[]int{0, 0},
		[]int{0, 0},
		[]int{0, 0},
		[]int{0, 0},
	}
	shots := []bool{
		false,
		true,
		true,
		false,
	}

	m.Commit(scores, shots)
	assert.Equal(1, m.Players[1].Shots)
	assert.Equal(1, m.Players[2].Shots)
}

func TestCorrectColorConflictsNoScores(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	_ = m.AddPlayer(Player{Name: "a", PreferredColor: "green"})
	_ = m.AddPlayer(Player{Name: "b", PreferredColor: "green"})
	_ = m.AddPlayer(Player{Name: "c", PreferredColor: "blue"})
	_ = m.AddPlayer(Player{Name: "d", PreferredColor: "pink"})

	err := m.CorrectColorConflicts()
	assert.Nil(err)

	assert.Equal("green", m.Players[0].PreferredColor)
	assert.NotEqual("green", m.Players[1].PreferredColor)
}

func TestCorrectColorConflictsNoScoresDoubleConflict(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	_ = m.AddPlayer(Player{Name: "a", PreferredColor: "green"})
	_ = m.AddPlayer(Player{Name: "b", PreferredColor: "green"})
	_ = m.AddPlayer(Player{Name: "c", PreferredColor: "blue"})
	_ = m.AddPlayer(Player{Name: "d", PreferredColor: "blue"})

	err := m.CorrectColorConflicts()
	assert.Nil(err)

	assert.Equal("green", m.Players[0].PreferredColor)
	assert.NotEqual("green", m.Players[1].PreferredColor)
	assert.Equal("blue", m.Players[2].PreferredColor)
	assert.NotEqual("blue", m.Players[3].PreferredColor)
}

func TestCorrectColorConflictPlayerTwoHasHigherScore(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	_ = m.AddPlayer(Player{Name: "a", PreferredColor: "green"})
	_ = m.AddPlayer(Player{Name: "b", PreferredColor: "green"})
	_ = m.AddPlayer(Player{Name: "c", PreferredColor: "blue"})
	_ = m.AddPlayer(Player{Name: "d", PreferredColor: "cyan"})

	// Add some score to player 2 so that it has preference over green.
	m.Players[1].AddKill(3)

	err := m.CorrectColorConflicts()
	assert.Nil(err)

	// TODO(thiderman): Fix this
	// assert.NotEqual("green", m.Players[0].PreferredColor)
	// assert.Equal("green", m.Players[1].PreferredColor)
	// assert.Equal("blue", m.Players[2].PreferredColor)
	// assert.Equal("cyan", m.Players[3].PreferredColor)
}
