package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreWithShots(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddShot()

	assert.Equal(3, p.Score())
	p.AddShot()
	assert.Equal(6, p.Score())
}

func TestScoreWithSweeps(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddSweep()

	assert.Equal(14, p.Score())
	p.AddSweep()
	assert.Equal(28, p.Score())
}

func TestScoreWithKills(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddKill()

	assert.Equal(2, p.Score())
	p.AddKill()
	assert.Equal(4, p.Score())
}

func TestScoreWithSelfs(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddSelf()

	assert.Equal(4, p.Score())
	p.AddSelf()
	assert.Equal(8, p.Score())
}

func TestScoreWithExplosions(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddExplosion()

	assert.Equal(6, p.Score())
	p.AddExplosion()
	assert.Equal(12, p.Score())
}

func TestScoreWithAll(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddShot()
	p.AddSweep()
	p.AddKill()
	p.AddSelf()
	p.AddExplosion()
	assert.Equal(27, p.Score())
}

func TestAddShot(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()

	p.AddShot()
	assert.Equal(1, p.shots)
	p.AddShot()
	assert.Equal(2, p.shots)
}

func TestRemoveShot(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.shots = 1

	p.RemoveShot()
	assert.Equal(0, p.shots)
	p.RemoveShot()
	assert.Equal(0, p.shots)
}

func TestAddSweep(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()

	p.AddSweep()
	assert.Equal(1, p.sweeps)
	assert.Equal(1, p.shots)
	assert.Equal(3, p.kills)

	p.AddSweep()
	assert.Equal(2, p.sweeps)
	assert.Equal(2, p.shots)
	assert.Equal(6, p.kills)
}

func TestRemoveSweep(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.sweeps = 1
	p.shots = 1
	p.kills = 3

	p.RemoveSweep()
	assert.Equal(0, p.sweeps)
	assert.Equal(0, p.shots)
	assert.Equal(0, p.kills)
	p.RemoveSweep()
	assert.Equal(0, p.sweeps)
	assert.Equal(0, p.shots)
	assert.Equal(0, p.kills)
}

func TestAddKill(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()

	p.AddKill()
	assert.Equal(1, p.kills)
	p.AddKill()
	assert.Equal(2, p.kills)
}

func TestRemoveKill(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.kills = 1

	p.RemoveKill()
	assert.Equal(0, p.kills)
	p.RemoveKill()
	assert.Equal(0, p.kills)
}

func TestAddSelf(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()

	p.AddSelf()
	assert.Equal(1, p.self)
	assert.Equal(1, p.shots)
	p.AddSelf()
	assert.Equal(2, p.self)
	assert.Equal(2, p.shots)
}

func TestRemoveSelf(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.self = 1
	p.shots = 1

	p.RemoveSelf()
	assert.Equal(0, p.self)
	assert.Equal(0, p.shots)
	p.RemoveSelf()
	assert.Equal(0, p.self)
	assert.Equal(0, p.shots)
}

func TestAddExplosion(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()

	p.AddExplosion()
	assert.Equal(1, p.explosions)
	assert.Equal(1, p.shots)
	assert.Equal(1, p.kills)
	p.AddExplosion()
	assert.Equal(2, p.explosions)
	assert.Equal(2, p.shots)
	assert.Equal(2, p.kills)
}

func TestRemoveExplosion(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.explosions = 1

	p.RemoveExplosion()
	assert.Equal(0, p.explosions)
	p.RemoveExplosion()
	assert.Equal(0, p.explosions)
}

func TestSortPlayers(t *testing.T) {
	assert := assert.New(t)
	p1 := NewPlayer() // 2 points
	p2 := NewPlayer() // 14 points
	p3 := NewPlayer() // 7 points

	p1.AddKill()
	p2.AddSweep()
	p3.AddShot()
	p3.AddSelf()

	ps := []Player{p1, p2, p3}
	ret := SortByScore(ps)

	assert.Equal(14, ret[0].Score())
	assert.Equal(7, ret[1].Score())
	assert.Equal(2, ret[2].Score())
}

func TestSortRunnerups(t *testing.T) {
	assert := assert.New(t)
	p1 := NewPlayer() // 10 points, 1 match
	p2 := NewPlayer() // 20 points, 2 matches
	p3 := NewPlayer() // 16 points, 1 match

	p1.kills = 5
	p1.matches = 1
	p1.name = "second"
	p2.kills = 10
	p2.matches = 2
	p2.name = "last"
	p3.kills = 8
	p3.matches = 1
	p3.name = "first"

	ps := []Player{p1, p2, p3}
	ret := SortByRunnerup(ps)

	assert.Equal("first", ret[0].name)
	assert.Equal("second", ret[1].name)
	assert.Equal("last", ret[2].name)
}
