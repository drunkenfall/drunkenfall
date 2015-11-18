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
	assert.Equal(29, p.Score())
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
