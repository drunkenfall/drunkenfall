package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(1, "test")
	p := Player{}

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, len(m.Players))
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(1, "test")
	m.Players = []Player{Player{}, Player{}, Player{}, Player{}}
	p := Player{}

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(1, "test")
	m.Started = time.Now()

	err := m.StartMatch()
	assert.NotNil(err)
}

func TestStartMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(1, "test")

	err := m.StartMatch()
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndAlreadyEndedMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(1, "test")
	m.Ended = time.Now()

	err := m.EndMatch()
	assert.NotNil(err)
}

func TestEndMatch(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch(1, "test")

	err := m.EndMatch()
	assert.Nil(err)
	assert.Equal(false, m.Ended.IsZero())
}
