package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch()
	p := Player{}

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, len(m.Players))
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	m := NewMatch()
	m.Players = []Player{Player{}, Player{}, Player{}, Player{}}
	p := Player{}

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}
