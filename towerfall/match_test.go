package towerfall

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/drunkenfall/drunkenfall/faking"
	"github.com/stretchr/testify/assert"
)

// MockMatch makes a mock Match{} with a dummy Tournament{}
func MockMatch(idx int, cat string) (m *Match) {
	tm := testTournament(8)
	tm.SetMatchPointers()
	err := tm.StartTournament(nil)
	if err != nil {
		log.Fatal(err)
	}

	offset := 0

	switch cat {
	case playoff:
		m = tm.Matches[offset+idx]
	case semi:
		offset = len(tm.Matches) - 3
		m = tm.Matches[offset+idx]
	case final:
		m = tm.Matches[len(tm.Matches)-1]
	default:
		log.Fatalf("Unknown match type: %s", cat)
	}

	return m
}

func testPlayer() Player {
	return *NewPlayer(testPerson(rand.Int()))
}

func testPerson(i int) *Person {
	return &Person{
		PersonID: fmt.Sprintf("%d: %s", i, faking.FakeName()),
		Name:     fmt.Sprintf("%d: %s", i, faking.FakeName()),
		Nick:     faking.FakeNick(),
		ColorPreference: []string{
			RandomColor(Colors),
			RandomColor(Colors),
		},
	}
}

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(0, "playoff")
	m.Players = []Player{}
	p := testPlayer()

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, len(m.Players))
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "playoff")

	p := testPlayer()

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "playoff")
	m.Started = time.Now()

	err := m.Start(nil)
	assert.NotNil(err)
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "playoff")

	err := m.Start(nil)
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndGivesShotToWinner(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "playoff")

	err := m.Start(nil)
	assert.Nil(err)
	m.Players[2].AddKills(10)
	// m.KillOrder = m.MakeKillOrder()

	err = m.End(nil)
	assert.Nil(err)
	assert.Equal(1, m.Players[2].Shots)
}

func TestEndAlreadyEndedMatch(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(1, "playoff")
	m.Ended = time.Now()

	err := m.End(nil)
	assert.NotNil(err)
}

func TestCommitSweepPlayer1(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{3, 0},
			{0, 0},
			{0, 0},
			{0, 0},
		},
		Shots: []bool{
			false,
			false,
			false,
			false,
		},
	}

	m.Commit(c)
	assert.Equal(1, m.Players[0].Sweeps)
}

func TestCommitDoubleKillPlayer2(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{0, 0},
			{2, 0},
			{0, 0},
			{0, 0},
		},
		Shots: []bool{
			false,
			false,
			false,
			false,
		},
	}

	m.Commit(c)
	assert.Equal(2, m.Players[1].Kills)
}

func TestCommitSweepAndSuicidePlayer3(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{0, 0},
			{0, 0},
			{3, -1},
			{0, 0},
		},
		Shots: []bool{
			false,
			false,
			true,
			false,
		},
	}

	m.Commit(c)
	assert.Equal(1, m.Players[2].Sweeps)
	assert.Equal(2, m.Players[2].Kills)
	assert.Equal(1, m.Players[2].Shots)
}

func TestCommitSuicidePlayer4(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{0, 0},
			{0, 0},
			{0, 0},
			{0, -1},
		},
		Shots: []bool{
			false,
			false,
			false,
			false,
		},
	}

	m.Commit(c)
	assert.Equal(1, m.Players[3].Self)
	assert.Equal(1, m.Players[3].Shots)
}

func TestCommitShotsForPlayer2and3(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{0, 0},
			{0, 0},
			{0, 0},
			{0, 0},
		},
		Shots: []bool{
			false,
			true,
			true,
			false,
		},
	}

	m.Commit(c)
	assert.Equal(1, m.Players[1].Shots)
	assert.Equal(1, m.Players[2].Shots)
}

func TestCommitSweepForPlayer1(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{3, 0},
			{0, 0},
			{0, 0},
			{0, 0},
		},
		// For the frontend it makes sense that a sweep marks a shot, therefore we
		// need to make sure that we don't add another shot.
		Shots: []bool{
			true,
			false,
			false,
			false,
		},
	}

	m.Commit(c)
	assert.Equal(3, m.Players[0].Kills)
	assert.Equal(1, m.Players[0].Shots)
}

func TestCommitStoredOnMatch(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{0, 0}, {1, 0}, {1, 0}, {1, 0},
		},
		// For the frontend it makes sense that a sweep marks a shot, therefore we
		// need to make sure that we don't add another shot.
		Shots: []bool{true, false, false, false},
	}

	assert.Equal(0, len(m.Rounds))
	m.Commit(c)
	assert.Equal(1, len(m.Rounds))
	assert.Equal(1, m.Rounds[0].Kills[1][0])
	assert.Equal(1, m.Rounds[0].Kills[2][0])
	assert.Equal(1, m.Rounds[0].Kills[3][0])
	assert.Equal(true, m.Rounds[0].Shots[0])
}

func TestCommitWithOnlyShotsNotStoredOnMatch(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "playoff")
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())
	_ = m.AddPlayer(testPlayer())

	c := Round{
		Kills: [][]int{
			{0, 0}, {0, 0}, {0, 0}, {0, 0},
		},
		Shots: []bool{true, false, false, false},
	}

	assert.Equal(0, len(m.Rounds))
	m.Commit(c)
	assert.Equal(0, len(m.Rounds))
	assert.Equal(1, m.Players[0].Shots)
}

func TestCorrectColorConflictsNoScores(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].PreferredColor = "green"
	m.Tournament.Players[1].PreferredColor = "green"
	m.Tournament.Players[2].PreferredColor = "blue"
	m.Tournament.Players[3].PreferredColor = "red"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3]))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].PreferredColor)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].PreferredColor)
	assert.Equal("red", m.Players[3].Color)
	assert.Equal("red", m.Players[3].PreferredColor)

	// Just one event - nothing started or ended
	assert.Equal(1, len(m.Events))
}

func TestCorrectColorConflictsUserLevels(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].PreferredColor = "green"
	m.Tournament.Players[0].Person.Userlevel = 10000
	m.Tournament.Players[1].PreferredColor = "green"
	m.Tournament.Players[2].PreferredColor = "red"
	m.Tournament.Players[2].Person.Userlevel = -10000
	m.Tournament.Players[3].PreferredColor = "red"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3]))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].PreferredColor)
	assert.NotEqual("red", m.Players[2].Color)
	assert.Equal("red", m.Players[2].PreferredColor)
	assert.Equal("red", m.Players[3].Color)
	assert.Equal("red", m.Players[3].PreferredColor)

	assert.Equal(2, len(m.Events))
}

// TODO(thiderman): I think the premise of this test is sound, but the
// execution is wrong since that's not how the things are actually
// executed. Should probably rewrite into all players starting with
// the same color.

// This test was needed since somehow the color were being kept
// func TestCorrectColorConflictsResetsToPreferredColor(t *testing.T) {
// 	assert := assert.New(t)

// 	tm := testTournament(12)
// 	tm.StartTournament(nil)
// 	m := tm.Matches[0]
// 	m2 := tm.Matches[1]
// 	m.Players = make([]Player, 0)
// 	m2.Players = make([]Player, 0)

// 	m.Tournament.Players[0].PreferredColor = "green"
// 	m.Tournament.Players[1].PreferredColor = "green"
// 	m.Tournament.Players[2].PreferredColor = "green"
// 	m.Tournament.Players[3].PreferredColor = "green"
// 	m.Tournament.Players[4].PreferredColor = "green"

// 	assert.Nil(m.AddPlayer(m.Tournament.Players[0]))
// 	assert.Nil(m.AddPlayer(m.Tournament.Players[1]))
// 	assert.Nil(m.AddPlayer(m.Tournament.Players[2]))
// 	assert.Nil(m.AddPlayer(m.Tournament.Players[3]))

// 	assert.Nil(m2.AddPlayer(m.Tournament.Players[4]))
// 	assert.Nil(m2.AddPlayer(m.Tournament.Players[1]))
// 	assert.Nil(m2.AddPlayer(m.Tournament.Players[2]))
// 	assert.Nil(m2.AddPlayer(m.Tournament.Players[3]))

// 	assert.Equal("green", m.Players[0].Color)
// 	assert.Equal("green", m.Players[0].PreferredColor)
// 	assert.NotEqual("green", m.Players[1].Color)
// 	assert.Equal("green", m.Players[1].PreferredColor)
// 	assert.NotEqual("green", m.Players[2].Color)
// 	assert.Equal("green", m.Players[2].PreferredColor)
// 	assert.NotEqual("green", m.Players[3].Color)

// 	assert.Equal("green", m2.Players[0].Color)
// 	assert.Equal("green", m2.Players[0].PreferredColor)
// 	assert.NotEqual("green", m2.Players[1].Color)
// 	assert.Equal("green", m2.Players[1].PreferredColor)
// 	assert.NotEqual("green", m2.Players[2].Color)
// 	assert.Equal("green", m2.Players[2].PreferredColor)
// 	assert.NotEqual("green", m2.Players[3].Color)

// 	assert.Equal(3, len(m.Events))
// 	assert.Equal(3, len(m2.Events))
// }

func TestCorrectColorConflictsNoScoresDoubleConflict(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].PreferredColor = "green"
	m.Tournament.Players[1].PreferredColor = "green"
	m.Tournament.Players[2].PreferredColor = "blue"
	m.Tournament.Players[3].PreferredColor = "blue"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3]))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].PreferredColor)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].PreferredColor)
	assert.NotEqual("blue", m.Players[3].Color)
	assert.Equal("blue", m.Players[3].PreferredColor)

	assert.Equal(2, len(m.Events))
}

func TestCorrectColorConflictsWithScoresDoubleConflict(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].PreferredColor = "green"
	m.Tournament.Players[0].Person.Nick = "GreenCorrected"

	m.Tournament.Players[1].Kills = 3
	m.Tournament.Players[1].PreferredColor = "green"

	m.Tournament.Players[2].PreferredColor = "blue"
	m.Tournament.Players[2].Person.Nick = "BlueCorrected"

	m.Tournament.Players[3].Kills = 3
	m.Tournament.Players[3].PreferredColor = "blue"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3]))

	assert.NotEqual("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].PreferredColor)
	assert.Equal("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].PreferredColor)
	assert.NotEqual("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].PreferredColor)
	assert.Equal("blue", m.Players[3].Color)
	assert.Equal("blue", m.Players[3].PreferredColor)

	assert.Equal(2, len(m.Events))
}

func TestCorrectColorConflictsWithScoresTripleConflict(t *testing.T) {
	assert := assert.New(t)

	m := MockMatch(0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].PreferredColor = "green"
	m.Tournament.Players[0].Person.Nick = "Green1Corrected"

	m.Tournament.Players[1].Kills = 3
	m.Tournament.Players[1].PreferredColor = "green"
	m.Tournament.Players[1].Person.Nick = "Green2Corrected"

	m.Tournament.Players[2].PreferredColor = "blue"
	m.Tournament.Players[2].Person.Nick = "BlueKeep"

	m.Tournament.Players[3].Kills = 10
	m.Tournament.Players[3].PreferredColor = "green"
	m.Tournament.Players[3].Person.Nick = "GreenKeep"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2]))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3]))

	assert.NotEqual("green", m.Players[0].Color)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("green", m.Players[3].Color)

	assert.Equal(2, len(m.Events))
}

func TestMakeKillOrder(t *testing.T) {
	assert := assert.New(t)
	m := MockMatch(0, "playoff")

	m.Players[0].AddKills(1)
	m.Players[1].AddKills(4)
	m.Players[2].AddKills(5)
	m.Players[3].AddKills(10)

	ko := m.MakeKillOrder()

	// As long as the order is reversed, this test is proven.
	// ...just like above. <3
	assert.Equal(ko[0], 3)
	assert.Equal(ko[1], 2)
	assert.Equal(ko[2], 1)
	assert.Equal(ko[3], 0)
}

func TestRoundIsShotUpdate(t *testing.T) {
	t.Run("Has shots", func(t *testing.T) {
		r := Round{
			Kills: [][]int{
				{0, 0}, {0, 0}, {0, 0}, {0, 0},
			},
			Shots: []bool{true, false, false, false},
		}
		assert.Equal(t, true, r.IsShotUpdate())
	})

	t.Run("Has kills", func(t *testing.T) {
		r := Round{
			Kills: [][]int{
				{0, 0}, {1, 0}, {1, 0}, {0, 0},
			},
			Shots: []bool{true, false, false, false},
		}
		assert.Equal(t, false, r.IsShotUpdate())
	})

	t.Run("Is completely empty", func(t *testing.T) {
		r := Round{
			Kills: [][]int{
				{0, 0}, {0, 0}, {0, 0}, {0, 0},
			},
			Shots: []bool{false, false, false, false},
		}
		assert.Equal(t, false, r.IsShotUpdate())
	})
}

func TestKill(t *testing.T) {
	t.Run("Kill by other player", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		km := KillMessage{1, 2, rArrow}
		err := m.Kill(km)
		assert.NoError(t, err)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "kill", m.Events[ev].Kind)
		assert.Equal(t, rArrow, m.Events[ev].Items["cause"])
	})

	t.Run("Environment kill", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		km := KillMessage{1, EnvironmentKill, rExplosion}
		err := m.Kill(km)
		assert.NoError(t, err)
		assert.Equal(t, 1, m.Players[1].Self)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "kill_environ", m.Events[ev].Kind)
		assert.Equal(t, rExplosion, m.Events[ev].Items["cause"])
	})

	t.Run("Suicide", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		km := KillMessage{1, 1, rCurse}
		err := m.Kill(km)
		assert.NoError(t, err)
		assert.Equal(t, 1, m.Players[1].Self)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "suicide", m.Events[ev].Kind)
		assert.Equal(t, rCurse, m.Events[ev].Items["cause"])
	})
}

func TestLavaOrb(t *testing.T) {
	t.Run("Enable", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		lm := LavaOrbMessage{0, true}
		err := m.LavaOrb(lm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Lava)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "lava", m.Events[ev].Kind)
	})

	t.Run("Disable", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		lm := LavaOrbMessage{0, false}
		err := m.LavaOrb(lm)
		assert.NoError(t, err)
		assert.Equal(t, false, m.Players[0].State.Lava)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "lava_off", m.Events[ev].Kind)
	})
}

func TestShield(t *testing.T) {
	t.Run("Gain", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		sm := ShieldMessage{0, true}
		err := m.ShieldUpdate(sm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Shield)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "shield", m.Events[ev].Kind)
	})

	t.Run("Break", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)
		m.Players[0].State.Shield = true

		sm := ShieldMessage{0, false}
		err := m.ShieldUpdate(sm)
		assert.NoError(t, err)
		assert.Equal(t, false, m.Players[0].State.Shield)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "shield_off", m.Events[ev].Kind)
	})
}

func TestWings(t *testing.T) {
	t.Run("Grow", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)

		wm := WingsMessage{0, true}
		err := m.WingsUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Wings)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "wings", m.Events[ev].Kind)
	})

	t.Run("Lose", func(t *testing.T) {
		m := MockMatch(0, playoff)
		ev := len(m.Events)
		m.Players[0].State.Wings = true

		wm := WingsMessage{0, false}
		err := m.WingsUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, false, m.Players[0].State.Wings)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "wings_off", m.Events[ev].Kind)
	})
}

func TestArrow(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		m := MockMatch(0, playoff)

		wm := ArrowMessage{3, Arrows{aNormal, aPrism, aPrism}}
		err := m.ArrowUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, Arrows{aNormal, aPrism, aPrism}, m.Players[3].State.Arrows)
	})
}

func TestStartRound(t *testing.T) {
	tm := testTournament(12)
	err := tm.StartTournament(nil)
	assert.NoError(t, err)
	m := tm.Matches[0]

	// These are absurd, of course...
	sr := StartRoundMessage{[]Arrows{
		[]int{aNormal, aBomb, aNormal},
		[]int{aSuperBomb, aBolt, aPrism},
		[]int{aNormal, aNormal, aNormal},
		[]int{aBomb, aBomb, aBomb},
	}}
	t.Run("Arrows are set", func(t *testing.T) {
		err := m.StartRound(sr)
		assert.NoError(t, err)

		assert.Equal(t, Arrows{aNormal, aBomb, aNormal}, m.Players[0].State.Arrows)
		assert.Equal(t, Arrows{aSuperBomb, aBolt, aPrism}, m.Players[1].State.Arrows)
		assert.Equal(t, Arrows{aNormal, aNormal, aNormal}, m.Players[2].State.Arrows)
		assert.Equal(t, Arrows{aBomb, aBomb, aBomb}, m.Players[3].State.Arrows)
	})
}
