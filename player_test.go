package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunnerupScoreWithShots(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddShot()

	assert.Equal(3, p.RunnerupScore())
	p.AddShot()
	assert.Equal(6, p.RunnerupScore())
}

func TestRunnerupScoreWithSweeps(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddSweep()

	assert.Equal(14, p.RunnerupScore())
	p.AddSweep()
	assert.Equal(28, p.RunnerupScore())
}

func TestRunnerupScoreWithKills(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddKill()

	assert.Equal(2, p.RunnerupScore())
	p.AddKill()
	assert.Equal(4, p.RunnerupScore())
}

func TestRunnerupScoreWithSelfs(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddSelf()

	assert.Equal(4, p.RunnerupScore())
	p.AddSelf()
	assert.Equal(8, p.RunnerupScore())
}

func TestRunnerupScoreWithExplosions(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddExplosion()

	assert.Equal(6, p.RunnerupScore())
	p.AddExplosion()
	assert.Equal(12, p.RunnerupScore())
}

func TestRunnerupScoreWithAll(t *testing.T) {
	assert := assert.New(t)
	p := NewPlayer()
	p.AddShot()
	p.AddSweep()
	p.AddKill()
	p.AddSelf()
	p.AddExplosion()
	assert.Equal(29, p.RunnerupScore())
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
