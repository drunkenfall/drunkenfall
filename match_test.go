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
	p := Player{}

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, len(m.Players))
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")
	m.Players = []Player{Player{}, Player{}, Player{}, Player{}}
	p := Player{}

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")
	m.Started = time.Now()

	err := m.StartMatch()
	assert.NotNil(err)
}

func TestStartMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")

	err := m.StartMatch()
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndAlreadyEndedMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")
	m.Ended = time.Now()

	err := m.EndMatch()
	assert.NotNil(err)
}

func TestEndMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(tm, 1, "test")

	err := m.EndMatch()
	assert.Nil(err)
	assert.Equal(false, m.Ended.IsZero())
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	m := NewMatch(tm, 0, "test")
	m.AddPlayer(Player{name: "1"})
	m.AddPlayer(Player{name: "2"})
	m.AddPlayer(Player{name: "3"})
	m.AddPlayer(Player{name: "4"})
	ret := m.String()

	assert.Equal("Test 1: 1 / 2 / 3 / 4 - not started", ret)

	m2 := NewMatch(tm, 0, "final")
	m2.AddPlayer(Player{name: "a"})
	m2.AddPlayer(Player{name: "b"})
	m2.AddPlayer(Player{name: "c"})
	m2.AddPlayer(Player{name: "d"})
	m2.StartMatch()
	ret2 := m2.String()

	assert.Equal("Final: a / b / c / d - playing", ret2)
}
