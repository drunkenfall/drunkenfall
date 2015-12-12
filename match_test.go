package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// These tests do not care about the tournaments,
// but the reference is required in NewMatch()
var tm *Tournament

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")

	assert.Equal(4, len(m.Players))
	assert.Equal(0, m.ActualPlayers())

	p := Player{Name: "I exist"}

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, m.ActualPlayers())
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")

	m.Players = []Player{Player{Name: "a"}, Player{Name: "b"}, Player{Name: "c"}, Player{Name: "d"}}
	p := Player{}

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")
	m.Started = time.Now()

	err := m.Start()
	assert.NotNil(err)
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")
	m.Players = []Player{
		Player{Name: "1"},
		Player{Name: "2"},
		Player{Name: "3"},
		Player{Name: "4"},
	}

	err := m.Start()
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndGivesShotToWinner(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "tryout")
	m.Players = []Player{
		Player{Name: "1"},
		Player{Name: "2"},
		Player{Name: "3"},
		Player{Name: "4"},
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
	m := NewMatch(tm, 1, "test")
	m.Ended = time.Now()

	err := m.End()
	assert.NotNil(err)
}

func TestEnd(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")

	err := m.End()
	assert.Nil(err)
	assert.Equal(false, m.Ended.IsZero())
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	m := NewMatch(tm, 0, "test")
	m.AddPlayer(Player{Name: "1"})
	m.AddPlayer(Player{Name: "2"})
	m.AddPlayer(Player{Name: "3"})
	m.AddPlayer(Player{Name: "4"})
	ret := m.String()

	assert.Equal("<Test 1: 1 / 2 / 3 / 4 - not started>", ret)

	m2 := NewMatch(tm, 0, "final")
	m2.AddPlayer(Player{Name: "a"})
	m2.AddPlayer(Player{Name: "b"})
	m2.AddPlayer(Player{Name: "c"})
	m2.AddPlayer(Player{Name: "d"})
	m2.Start()
	ret2 := m2.String()

	assert.Equal("<Final: a / b / c / d - playing>", ret2)
}
