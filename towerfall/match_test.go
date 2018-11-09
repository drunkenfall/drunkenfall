package towerfall

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/drunkenfall/drunkenfall/faking"
	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
)

// People that have been used for a tournament. Used to make sure we
// don't randomly grab one we already have grabbed
var usedPeople []string

// MockMatch makes a mock Match{} with a dummy Tournament{}
func MockMatch(t *testing.T, s *Server, idx int, cat string) (m *Match) {
	tm := testTournament(t, s, 8)
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

func testPlayer(s *Server) Player {
	return *NewPlayer(testPerson(s))
}

func testPerson(s *Server) *Person {
	p, err := s.DB.GetRandomPerson(usedPeople)
	if err != nil {
		log.Fatal(err)
	}

	usedPeople = append(usedPeople, p.PersonID)
	return p
}

func randomPerson() *Person {
	return &Person{
		PersonID:       faking.FakeName(),
		Name:           faking.FakeName(),
		Nick:           faking.FakeNick(),
		PreferredColor: RandomColor(Colors),
	}
}

func TestAddPlayer(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	m.Players = []Player{}
	p := testPlayer(s)

	err := m.AddPlayer(p)
	assert.Nil(err)

	assert.Equal(1, len(m.Players))
}

func TestAddFifthPlayer(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, "playoff")

	p := testPlayer(s)

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, "playoff")
	m.Started = time.Now()

	err := m.Start(nil)
	assert.NotNil(err)
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, "playoff")

	err := m.Start(nil)
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndGivesShotToWinner(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, "playoff")

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
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, "playoff")
	m.Ended = time.Now()

	err := m.End(nil)
	assert.NotNil(err)
}

func TestCommitSweepPlayer1(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))
	_ = m.AddPlayer(testPlayer(s))

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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].Person.PreferredColor = "green"
	m.Tournament.Players[1].Person.PreferredColor = "green"
	m.Tournament.Players[2].Person.PreferredColor = "blue"
	m.Tournament.Players[3].Person.PreferredColor = "red"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].Person.PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].Person.PreferredColor)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].Person.PreferredColor)
	assert.Equal("red", m.Players[3].Color)
	assert.Equal("red", m.Players[3].Person.PreferredColor)

	// Just one event - nothing started or ended
	assert.Equal(1, len(m.Events))
}

func TestCorrectColorConflictsUserLevels(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].Person.PreferredColor = "green"
	m.Tournament.Players[0].Person.Userlevel = 10000
	m.Tournament.Players[1].Person.PreferredColor = "green"
	m.Tournament.Players[2].Person.PreferredColor = "red"
	m.Tournament.Players[2].Person.Userlevel = -10000
	m.Tournament.Players[3].Person.PreferredColor = "red"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].Person.PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].Person.PreferredColor)
	assert.NotEqual("red", m.Players[2].Color)
	assert.Equal("red", m.Players[2].Person.PreferredColor)
	assert.Equal("red", m.Players[3].Color)
	assert.Equal("red", m.Players[3].Person.PreferredColor)

	assert.Equal(2, len(m.Events))
}

// TODO(thiderman): I think the premise of this test is sound, but the
// execution is wrong since that's not how the things are actually
// executed. Should probably rewrite into all players starting with
// the same color.

// This test was needed since somehow the color were being kept
// func TestCorrectColorConflictsResetsToPreferredColor(t *testing.T) {
// 	assert := assert.New(t)

// 	tm := testTournament(t, s, 12)
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

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].Person.PreferredColor = "green"
	m.Tournament.Players[1].Person.PreferredColor = "green"
	m.Tournament.Players[2].Person.PreferredColor = "blue"
	m.Tournament.Players[3].Person.PreferredColor = "blue"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].Person.PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].Person.PreferredColor)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].Person.PreferredColor)
	assert.NotEqual("blue", m.Players[3].Color)
	assert.Equal("blue", m.Players[3].Person.PreferredColor)

	assert.Equal(2, len(m.Events))
}

func TestCorrectColorConflictsWithScoresDoubleConflict(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].Person.PreferredColor = "green"
	m.Tournament.Players[0].Person.Nick = "GreenCorrected"

	m.Tournament.Players[1].TotalScore = 3
	m.Tournament.Players[1].Person.PreferredColor = "green"

	m.Tournament.Players[2].Person.PreferredColor = "blue"
	m.Tournament.Players[2].Person.Nick = "BlueCorrected"

	m.Tournament.Players[3].TotalScore = 3
	m.Tournament.Players[3].Person.PreferredColor = "blue"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.NotEqual("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].Person.PreferredColor)
	assert.Equal("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].Person.PreferredColor)
	assert.NotEqual("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].Person.PreferredColor)
	assert.Equal("blue", m.Players[3].Color)
	assert.Equal("blue", m.Players[3].Person.PreferredColor)

	assert.Equal(2, len(m.Events))
}

func TestCorrectColorConflictsWithScoresTripleConflict(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].Person.PreferredColor = "green"
	m.Tournament.Players[0].Person.Nick = "Green1Corrected"

	m.Tournament.Players[1].TotalScore = 3
	m.Tournament.Players[1].Person.PreferredColor = "green"
	m.Tournament.Players[1].Person.Nick = "Green2Corrected"

	m.Tournament.Players[2].Person.PreferredColor = "blue"
	m.Tournament.Players[2].Person.Nick = "BlueKeep"

	m.Tournament.Players[3].TotalScore = 10
	m.Tournament.Players[3].Person.PreferredColor = "green"
	m.Tournament.Players[3].Person.Nick = "GreenKeep"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.NotEqual("green", m.Players[0].Color)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("green", m.Players[3].Color)

	assert.Equal(2, len(m.Events))
}

func TestMakeKillOrder(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "playoff")

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
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
		ev := len(m.Events)

		km := KillMessage{1, 2, rArrow}
		err := m.Kill(km)
		assert.NoError(t, err)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "kill", m.Events[ev].Kind)
		assert.Equal(t, rArrow, m.Events[ev].Items["cause"])
	})

	t.Run("Environment kill", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
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
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
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
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
		ev := len(m.Events)

		lm := LavaOrbMessage{0, true}
		err := m.LavaOrb(lm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Lava)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "lava", m.Events[ev].Kind)
	})

	t.Run("Disable", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
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
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
		ev := len(m.Events)

		sm := ShieldMessage{0, true}
		err := m.ShieldUpdate(sm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Shield)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "shield", m.Events[ev].Kind)
	})

	t.Run("Break", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
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
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
		ev := len(m.Events)

		wm := WingsMessage{0, true}
		err := m.WingsUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Wings)
		assert.Equal(t, ev+1, len(m.Events))
		assert.Equal(t, "wings", m.Events[ev].Kind)
	})

	t.Run("Lose", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)
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
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, playoff)

		wm := ArrowMessage{3, Arrows{aNormal, aPrism, aPrism}}
		err := m.ArrowUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, Arrows{aNormal, aPrism, aPrism}, m.Players[3].State.Arrows)
	})
}

func TestStartRound(t *testing.T) {
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 12)
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

func TestHandleMessage(t *testing.T) {
	s, teardown := MockServer(t)
	defer teardown()

	tm := testTournament(t, s, 12)
	err := tm.StartTournament(nil)
	assert.NoError(t, err)

	t.Run("Kill", func(t *testing.T) {
		m := tm.Matches[0]
		t.Run("Player on player", func(t *testing.T) {
			msg := Message{
				Type: inKill,
				Data: KillMessage{0, 1, rArrow},
			}

			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, false, m.Players[0].State.Alive)
			assert.Equal(t, 1, m.Players[0].State.Killer)
		})

		t.Run("Suicide", func(t *testing.T) {
			msg := Message{
				Type: inKill,
				Data: KillMessage{2, 2, rArrow},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, false, m.Players[2].State.Alive)
			assert.Equal(t, 2, m.Players[2].State.Killer)
		})

		t.Run("Environment kill", func(t *testing.T) {
			msg := Message{
				Type: inKill,
				Data: KillMessage{3, EnvironmentKill, rSpikeBall},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, false, m.Players[3].State.Alive)
			assert.Equal(t, EnvironmentKill, m.Players[3].State.Killer)
		})

	})

	t.Run("Round start", func(t *testing.T) {
		m := tm.Matches[0]
		t.Run("Reset", func(t *testing.T) {
			msg := Message{
				Type: inRoundStart,
				Data: StartRoundMessage{
					Arrows: []Arrows{
						{aNormal, aNormal, aNormal},
						{aNormal, aNormal, aNormal},
						{aNormal, aNormal, aNormal},
						{aNormal, aNormal, aNormal, aBomb},
					},
				},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}

			def := Arrows{aNormal, aNormal, aNormal}
			assert.Equal(t, def, m.Players[0].State.Arrows)
			assert.Equal(t, def, m.Players[1].State.Arrows)
			assert.Equal(t, def, m.Players[2].State.Arrows)
			assert.Equal(t, Arrows{aNormal, aNormal, aNormal, aBomb}, m.Players[3].State.Arrows)

			assert.Equal(t, true, m.Players[0].State.Alive)
			assert.Equal(t, true, m.Players[1].State.Alive)
			assert.Equal(t, true, m.Players[2].State.Alive)
			assert.Equal(t, true, m.Players[3].State.Alive)
		})
	})

	t.Run("Round end", func(t *testing.T) {
		m := tm.Matches[0]
		assert.Equal(t, 0, len(m.Rounds))

		t.Run("End", func(t *testing.T) {
			msg := Message{
				Type: inRoundEnd,
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, 1, len(m.Rounds))
		})
	})

	t.Run("Player updates", func(t *testing.T) {
		m := tm.Matches[0]

		// TODO(thiderman): Since we test without websockets right now,
		// the effects of these cannot be unit tested. However, we can
		// test that handleMessage() does its thing.

		t.Run("Shot", func(t *testing.T) {
			msg := Message{
				Type: inShot,
				Data: ArrowMessage{
					Player: 0,
					Arrows: Arrows{aBomb, aNormal, aNormal},
				},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Pickup", func(t *testing.T) {
			msg := Message{
				Type: inPickup,
				Data: ArrowMessage{
					Player: 0,
					Arrows: Arrows{aNormal, aNormal},
				},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Wings", func(t *testing.T) {
			msg := Message{
				Type: inWings,
				Data: WingsMessage{
					Player: 0,
					State:  true,
				},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}
		})

		t.Run("Lava", func(t *testing.T) {
			msg := Message{
				Type: inOrbLava,
				Data: LavaOrbMessage{
					Player: 2,
					State:  true,
				},
			}
			err := m.handleMessage(msg)
			if err != nil {
				t.Fatal(err)
			}
		})

	})
}

func replay(t *testing.T, db *pg.DB, m *Match, id int) {
	t.Helper()

	sm := Match{}
	// Grab the first match from it
	err := db.Model(&sm).Where("id = ?", id).First()
	assert.NoError(t, err)

	// Load all the messages from it
	msgs := []*Message{}
	err = db.Model(&msgs).Where("match_id = ?", sm.ID).Order("id").Select()
	if err != nil {
		t.Fatal(err)
	}

	for _, msg := range msgs {
		err = json.Unmarshal([]byte(msg.JSON), &msg.Data)
		if err != nil {
			t.Fatal(err)
		}

		msg.ID = 0 // Reset so that we can insert new ones for this match
		err = m.handleMessage(*msg)
		if !assert.NoError(t, err) {
			log.Printf("ERROR: %+v", err)
			t.Fatal(err)
		}
	}
}

func TestReplayLockStock(t *testing.T) {
	//  id | shots | sweeps | kills | self | match_score | total_score
	// ----+-------+--------+-------+------+-------------+-------------
	//  86 |     4 |      0 |     1 |    4 |           0 |           0
	//  86 |     2 |      0 |    10 |    1 |           0 |           0
	//  86 |     1 |      1 |     9 |    0 |           0 |           0
	//  86 |     2 |      0 |     6 |    2 |           0 |           0
	//  87 |     3 |      1 |     6 |    2 |           0 |           0
	//  87 |     5 |      2 |    10 |    2 |           0 |           0
	//  87 |     2 |      0 |     8 |    2 |           0 |           0
	//  87 |     1 |      0 |     6 |    1 |           0 |           0
	//  88 |     1 |      0 |     8 |    1 |           0 |           0
	//  88 |     4 |      0 |    10 |    3 |           0 |           0
	//  88 |     2 |      0 |     6 |    2 |           0 |           0
	//  88 |     6 |      2 |     9 |    4 |           0 |           0
	//  89 |     2 |      0 |     8 |    2 |           0 |           0
	//  89 |     1 |      0 |     8 |    1 |           0 |           0
	//  89 |     1 |      0 |    10 |    0 |           0 |          28
	//  89 |     1 |      0 |     8 |    1 |           0 |          20
	//  90 |     2 |      0 |     7 |    2 |           0 |           0
	//  90 |     2 |      0 |    11 |    1 |           0 |          24
	//  90 |     2 |      1 |    10 |    1 |           0 |           0
	//  90 |     0 |      0 |     6 |    0 |           0 |           0
	//  91 |     0 |      0 |     5 |    0 |           0 |           0
	//  91 |     0 |      0 |     9 |    0 |           0 |           0
	//  91 |     2 |      0 |    10 |    1 |           0 |          50
	//  91 |     3 |      1 |     7 |    2 |           0 |          28
	//  92 |     1 |      0 |    20 |    0 |           0 |          24
	//  92 |     4 |      1 |    15 |    3 |           0 |           0
	//  92 |     5 |      2 |    17 |    3 |           0 |          50
	//  92 |     3 |      0 |    10 |    3 |           0 |           0
	// (28 rows)

	s, teardown := MockServer(t)
	defer teardown()

	db := s.DB.DB
	tm := testTournament(t, s, 0)

	ids := []string{
		"2316298658399263",
		"10160764154475012",
		"10154542569541289",
		"10160309716400054",
		"10160969636925217",
		"10153943465786915",
		"1279099058796903",
		"10156829856779238",
		"10153964695568099",
		"1308346815842239",
		"10156062370636832",
		"10155017292347007",
		"10155849790765189",
		"10153861129901027",
	}

	for _, id := range ids {
		p, _ := s.DB.GetPerson(id)
		log.Printf("%s: %s", p.Nick, p.PreferredColor)
		s := NewPlayer(p).Summary()
		tm.AddPlayer(&s)
	}

	err := tm.StartTournament(nil)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	t.Run("Tryout1", func(t *testing.T) {
		m := tm.Matches[0]

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 86)
		})

		t.Run("Messages", func(t *testing.T) {
			assert.Equal(t, 572, len(m.Messages))
			check, err := db.Model(&Message{}).Where("match_id = ?", m.ID).Count()
			assert.NoError(t, err)
			assert.Equal(t, 572, check)
		})

		t.Run("Rounds and commits", func(t *testing.T) {
			assert.Equal(t, 12, len(m.Rounds))
			check, err := db.Model(&Commit{}).Where("match_id = ?", m.ID).Count()
			assert.NoError(t, err)
			assert.Equal(t, 12, check)
		})

		t.Run("Match state", func(t *testing.T) {
			assert.NotZero(t, m.Started)
			assert.NotZero(t, m.Ended)
		})

		// Finally check that the player stats are stored and that they
		// are the same
		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			t.Run("P1: 4 shots, 1 kills, 4 selfs", func(t *testing.T) {
				assert.Equal(t, 4, ps[0].Shots)
				assert.Equal(t, 0, ps[0].Sweeps)
				assert.Equal(t, 1, ps[0].Kills)
				assert.Equal(t, 4, ps[0].Self)
				assert.Equal(t, scoreFourth, ps[0].MatchScore)
				assert.Equal(t, 210, ps[0].Score())
			})

			t.Run("P2: 2 shots, 10 kills, 1 self - winner", func(t *testing.T) {
				assert.Equal(t, 2, ps[1].Shots)
				assert.Equal(t, 0, ps[1].Sweeps)
				assert.Equal(t, 10, ps[1].Kills)
				assert.Equal(t, 1, ps[1].Self)
				assert.Equal(t, scoreWinner, ps[1].MatchScore)
				assert.Equal(t, 3675, ps[1].Score())
			})

			t.Run("P3: 1 shot, 1 sweep, 9 kills, no selfs", func(t *testing.T) {
				assert.Equal(t, 1, ps[2].Shots)
				assert.Equal(t, 1, ps[2].Sweeps)
				assert.Equal(t, 9, ps[2].Kills)
				assert.Equal(t, 0, ps[2].Self)
				assert.Equal(t, scoreSecond, ps[2].MatchScore)
				assert.Equal(t, 3052, ps[2].Score())
			})

			t.Run("P4: 2 shots, 6 kills, 2 selfs", func(t *testing.T) {
				assert.Equal(t, 2, ps[3].Shots)
				assert.Equal(t, 0, ps[3].Sweeps)
				assert.Equal(t, 6, ps[3].Kills)
				assert.Equal(t, 2, ps[3].Self)
				assert.Equal(t, scoreThird, ps[3].MatchScore)
				assert.Equal(t, 882, ps[3].Score())
			})
		})

	})

	t.Run("Tryout2", func(t *testing.T) {
		m := tm.Matches[1]

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 87)
		})

		t.Run("Messages", func(t *testing.T) {
			assert.Equal(t, 628, len(m.Messages))
			check, err := db.Model(&Message{}).Where("match_id = ?", m.ID).Count()
			assert.NoError(t, err)
			assert.Equal(t, 628, check)
		})

		t.Run("Rounds and commits", func(t *testing.T) {
			assert.Equal(t, 13, len(m.Rounds))
			check, err := db.Model(&Commit{}).Where("match_id = ?", m.ID).Count()
			assert.NoError(t, err)
			assert.Equal(t, 13, check)
		})

		t.Run("Match state", func(t *testing.T) {
			assert.NotZero(t, m.Started)
			assert.NotZero(t, m.Ended)
		})

		// Finally check that the player stats are stored and that they
		// are the same
		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			t.Run("P1", func(t *testing.T) {
				assert.Equal(t, 3, ps[0].Shots)
				assert.Equal(t, 1, ps[0].Sweeps)
				assert.Equal(t, 6, ps[0].Kills)
				assert.Equal(t, 2, ps[0].Self)
				assert.Equal(t, scoreThird, ps[0].MatchScore)
				assert.Equal(t, 1561, ps[0].Score())
			})

			t.Run("P2", func(t *testing.T) {
				assert.Equal(t, 5, ps[1].Shots)
				assert.Equal(t, 2, ps[1].Sweeps)
				assert.Equal(t, 10, ps[1].Kills)
				assert.Equal(t, 2, ps[1].Self)
				assert.Equal(t, scoreWinner, ps[1].MatchScore)
				assert.Equal(t, 4788, ps[1].Score())
			})

			t.Run("P3", func(t *testing.T) {
				assert.Equal(t, 2, ps[2].Shots)
				assert.Equal(t, 0, ps[2].Sweeps)
				assert.Equal(t, 8, ps[2].Kills)
				assert.Equal(t, 2, ps[2].Self)
				assert.Equal(t, scoreSecond, ps[2].MatchScore)
				assert.Equal(t, 1736, ps[2].Score())
			})

			t.Run("P4", func(t *testing.T) {
				assert.Equal(t, 1, ps[3].Shots)
				assert.Equal(t, 0, ps[3].Sweeps)
				assert.Equal(t, 6, ps[3].Kills)
				assert.Equal(t, 1, ps[3].Self)
				assert.Equal(t, scoreFourth, ps[3].MatchScore)
				assert.Equal(t, 847, ps[3].Score())
			})
		})

	})

	t.Run("Tryout3", func(t *testing.T) {
		m := tm.Matches[2]

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 88)
		})

		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			// 88 |     1 |      0 |     8 |    1 |           0 |           0
			// 88 |     4 |      0 |    10 |    3 |           0 |           0
			// 88 |     2 |      0 |     6 |    2 |           0 |           0
			// 88 |     6 |      2 |     9 |    4 |           0 |           0
			t.Run("P1", func(t *testing.T) {
				assert.Equal(t, 1, ps[0].Shots)
				assert.Equal(t, 0, ps[0].Sweeps)
				assert.Equal(t, 8, ps[0].Kills)
				assert.Equal(t, 1, ps[0].Self)
				assert.Equal(t, scoreThird, ps[0].MatchScore)
				assert.Equal(t, 1421, ps[0].Score())
			})

			t.Run("P2", func(t *testing.T) {
				assert.Equal(t, 4, ps[1].Shots)
				assert.Equal(t, 0, ps[1].Sweeps)
				assert.Equal(t, 10, ps[1].Kills)
				assert.Equal(t, 3, ps[1].Self)
				assert.Equal(t, scoreWinner, ps[1].MatchScore)
				assert.Equal(t, 3185, ps[1].Score())
			})

			t.Run("P3", func(t *testing.T) {
				assert.Equal(t, 2, ps[2].Shots)
				assert.Equal(t, 0, ps[2].Sweeps)
				assert.Equal(t, 6, ps[2].Kills)
				assert.Equal(t, 2, ps[2].Self)
				assert.Equal(t, scoreFourth, ps[2].MatchScore)
				assert.Equal(t, 602, ps[2].Score())
			})

			t.Run("P4", func(t *testing.T) {
				assert.Equal(t, 6, ps[3].Shots)
				assert.Equal(t, 2, ps[3].Sweeps)
				assert.Equal(t, 9, ps[3].Kills)
				assert.Equal(t, 4, ps[3].Self)
				assert.Equal(t, scoreSecond, ps[3].MatchScore)
				assert.Equal(t, 2751, ps[3].Score())
			})
		})

	})

	t.Run("Tryout4", func(t *testing.T) {
		m := tm.Matches[3]

		// Fake runnerups since that doesn't work right now
		m.AddPlayer(tm.Matches[0].Players[3])
		m.AddPlayer(tm.Matches[1].Players[0])

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 89)
		})

		t.Run("Kill order", func(t *testing.T) {
			assert.Equal(t, []int{2, 1, 3, 0}, m.MakeKillOrder())
		})

		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			// 89 |     2 |      0 |     8 |    2 |           0 |           0
			// 89 |     1 |      0 |     8 |    1 |           0 |           0
			// 89 |     1 |      0 |    10 |    0 |           0 |          28
			// 89 |     1 |      0 |     8 |    1 |           0 |          20
			t.Run("P1", func(t *testing.T) {
				assert.Equal(t, 2, ps[0].Shots)
				assert.Equal(t, 0, ps[0].Sweeps)
				assert.Equal(t, 8, ps[0].Kills)
				assert.Equal(t, 2, ps[0].Self)
				assert.Equal(t, scoreFourth, ps[0].MatchScore)
				assert.Equal(t, 896, ps[0].Score())
			})

			t.Run("P2", func(t *testing.T) {
				assert.Equal(t, 1, ps[1].Shots)
				assert.Equal(t, 0, ps[1].Sweeps)
				assert.Equal(t, 8, ps[1].Kills)
				assert.Equal(t, 1, ps[1].Self)
				assert.Equal(t, scoreSecond, ps[1].MatchScore)
				assert.Equal(t, 1981, ps[1].Score())
			})

			t.Run("P3", func(t *testing.T) {
				assert.Equal(t, 1, ps[2].Shots)
				assert.Equal(t, 0, ps[2].Sweeps)
				assert.Equal(t, 10, ps[2].Kills)
				assert.Equal(t, 0, ps[2].Self)
				assert.Equal(t, scoreWinner, ps[2].MatchScore)
				assert.Equal(t, 3920, ps[2].Score())
			})

			t.Run("P4", func(t *testing.T) {
				assert.Equal(t, 1, ps[3].Shots)
				assert.Equal(t, 0, ps[3].Sweeps)
				assert.Equal(t, 8, ps[3].Kills)
				assert.Equal(t, 1, ps[3].Self)
				assert.Equal(t, scoreThird, ps[3].MatchScore)
				assert.Equal(t, 1421, ps[3].Score())
			})
		})

	})

	t.Run("Semi1", func(t *testing.T) {
		m := tm.Matches[4]

		t.Run("Players moved correctly", func(t *testing.T) {
			assert.Equal(t, tm.Matches[0].Players[1].PersonID, m.Players[0].PersonID)
			assert.Equal(t, tm.Matches[1].Players[2].PersonID, m.Players[1].PersonID)
			assert.Equal(t, tm.Matches[2].Players[1].PersonID, m.Players[2].PersonID)
			assert.Equal(t, tm.Matches[3].Players[1].PersonID, m.Players[3].PersonID)
		})

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 90)
		})

		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			//  90 |     2 |      0 |     7 |    2 |           0 |           0
			//  90 |     2 |      0 |    11 |    1 |           0 |          24
			//  90 |     2 |      1 |    10 |    1 |           0 |           0
			//  90 |     0 |      0 |     6 |    0 |           0 |           0
			t.Run("P1", func(t *testing.T) {
				assert.Equal(t, 2, ps[0].Shots)
				assert.Equal(t, 0, ps[0].Sweeps)
				assert.Equal(t, 7, ps[0].Kills)
				assert.Equal(t, 2, ps[0].Self)
				assert.Equal(t, scoreThird, ps[0].MatchScore)
				assert.Equal(t, 1029, ps[0].Score())
			})

			t.Run("P2", func(t *testing.T) {
				assert.Equal(t, 2, ps[1].Shots)
				assert.Equal(t, 0, ps[1].Sweeps)
				assert.Equal(t, 11, ps[1].Kills)
				assert.Equal(t, 1, ps[1].Self)
				assert.Equal(t, scoreWinner, ps[1].MatchScore)
				assert.Equal(t, 3822, ps[1].Score())
			})

			t.Run("P3", func(t *testing.T) {
				assert.Equal(t, 2, ps[2].Shots)
				assert.Equal(t, 1, ps[2].Sweeps)
				assert.Equal(t, 10, ps[2].Kills)
				assert.Equal(t, 1, ps[2].Self)
				assert.Equal(t, scoreSecond, ps[2].MatchScore)
				assert.Equal(t, 2954, ps[2].Score())
			})

			t.Run("P4", func(t *testing.T) {
				assert.Equal(t, 0, ps[3].Shots)
				assert.Equal(t, 0, ps[3].Sweeps)
				assert.Equal(t, 6, ps[3].Kills)
				assert.Equal(t, 0, ps[3].Self)
				assert.Equal(t, scoreFourth, ps[3].MatchScore)
				assert.Equal(t, 1092, ps[3].Score())
			})
		})

	})

	t.Run("Semi2", func(t *testing.T) {
		m := tm.Matches[5]

		t.Run("Players moved correctly", func(t *testing.T) {
			assert.Equal(t, tm.Matches[0].Players[2].PersonID, m.Players[0].PersonID)
			assert.Equal(t, tm.Matches[1].Players[1].PersonID, m.Players[1].PersonID)
			assert.Equal(t, tm.Matches[2].Players[3].PersonID, m.Players[2].PersonID)
			assert.Equal(t, tm.Matches[3].Players[2].PersonID, m.Players[3].PersonID)
		})

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 91)
		})

		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			//  91 |     0 |      0 |     5 |    0 |           0 |           0
			//  91 |     0 |      0 |     9 |    0 |           0 |           0
			//  91 |     2 |      0 |    10 |    1 |           0 |          50
			//  91 |     3 |      1 |     7 |    2 |           0 |          28
			t.Run("P1", func(t *testing.T) {
				assert.Equal(t, 0, ps[0].Shots)
				assert.Equal(t, 0, ps[0].Sweeps)
				assert.Equal(t, 5, ps[0].Kills)
				assert.Equal(t, 0, ps[0].Self)
				assert.Equal(t, scoreFourth, ps[0].MatchScore)
				assert.Equal(t, 945, ps[0].Score())
			})

			t.Run("P2", func(t *testing.T) {
				assert.Equal(t, 0, ps[1].Shots)
				assert.Equal(t, 0, ps[1].Sweeps)
				assert.Equal(t, 9, ps[1].Kills)
				assert.Equal(t, 0, ps[1].Self)
				assert.Equal(t, scoreSecond, ps[1].MatchScore)
				assert.Equal(t, 2373, ps[1].Score())
			})

			t.Run("P3", func(t *testing.T) {
				assert.Equal(t, 2, ps[2].Shots)
				assert.Equal(t, 0, ps[2].Sweeps)
				assert.Equal(t, 10, ps[2].Kills)
				assert.Equal(t, 1, ps[2].Self)
				assert.Equal(t, scoreWinner, ps[2].MatchScore)
				assert.Equal(t, 3675, ps[2].Score())
			})

			t.Run("P4", func(t *testing.T) {
				assert.Equal(t, 3, ps[3].Shots)
				assert.Equal(t, 1, ps[3].Sweeps)
				assert.Equal(t, 7, ps[3].Kills)
				assert.Equal(t, 2, ps[3].Self)
				assert.Equal(t, scoreThird, ps[3].MatchScore)
				assert.Equal(t, 1708, ps[3].Score())
			})
		})

	})
	t.Run("Fhphihnalhehehh", func(t *testing.T) {
		m := tm.Matches[6]

		t.Run("Players moved correctly", func(t *testing.T) {
			assert.Equal(t, tm.Matches[4].Players[1].PersonID, m.Players[0].PersonID)
			assert.Equal(t, tm.Matches[4].Players[2].PersonID, m.Players[1].PersonID)
			assert.Equal(t, tm.Matches[5].Players[2].PersonID, m.Players[2].PersonID)
			assert.Equal(t, tm.Matches[5].Players[1].PersonID, m.Players[3].PersonID)
		})

		t.Run("Replay all messages", func(t *testing.T) {
			replay(t, db, m, 92)
		})

		t.Run("Scores", func(t *testing.T) {
			ps := []*Player{}
			db.Model(&ps).Where("match_id = ?", m.ID).Order("id").Select()
			if !assert.Equal(t, 4, len(ps)) {
				return
			}

			//  92 |     1 |      0 |    20 |    0 |           0 |          24
			//  92 |     4 |      1 |    15 |    3 |           0 |           0
			//  92 |     5 |      2 |    17 |    3 |           0 |          50
			//  92 |     3 |      0 |    10 |    3 |           0 |           0
			t.Run("P1", func(t *testing.T) {
				assert.Equal(t, 1, ps[0].Shots)
				assert.Equal(t, 0, ps[0].Sweeps)
				assert.Equal(t, 20, ps[0].Kills)
				assert.Equal(t, 0, ps[0].Self)
				assert.Equal(t, int(float64(scoreWinner)*finalMultiplier), ps[0].MatchScore)
				assert.Equal(t, 9065, ps[0].Score())
			})

			t.Run("P2", func(t *testing.T) {
				assert.Equal(t, 4, ps[1].Shots)
				assert.Equal(t, 1, ps[1].Sweeps)
				assert.Equal(t, 15, ps[1].Kills)
				assert.Equal(t, 3, ps[1].Self)
				assert.Equal(t, int(float64(scoreThird)*finalMultiplier), ps[1].MatchScore)
				assert.Equal(t, 3374, ps[1].Score())
			})

			t.Run("P3", func(t *testing.T) {
				assert.Equal(t, 5, ps[2].Shots)
				assert.Equal(t, 2, ps[2].Sweeps)
				assert.Equal(t, 17, ps[2].Kills)
				assert.Equal(t, 3, ps[2].Self)
				assert.Equal(t, int(float64(scoreSecond)*finalMultiplier), ps[2].MatchScore)
				assert.Equal(t, 5747, ps[2].Score())
			})

			t.Run("P4", func(t *testing.T) {
				assert.Equal(t, 3, ps[3].Shots)
				assert.Equal(t, 0, ps[3].Sweeps)
				assert.Equal(t, 10, ps[3].Kills)
				assert.Equal(t, 3, ps[3].Self)
				assert.Equal(t, int(float64(scoreFourth)*finalMultiplier), ps[3].MatchScore)
				assert.Equal(t, 1260, ps[3].Score())
			})
		})

	})

	t.Run("End state", func(t *testing.T) {
		assert.NotZero(t, tm.Started)
		assert.NotZero(t, tm.Ended)
		assert.Equal(t, 130, tm.db.persistcalls)
	})

	t.Run("End scores", func(t *testing.T) {
		ps := []*PlayerSummary{}
		db.Model(&ps).Where("tournament_id = ?", tm.ID).Order("total_score DESC").Select()
		if !assert.Equal(t, len(tm.Players), len(ps)) {
			return
		}

		assert.Equal(t, 14623, ps[0].TotalScore)
		assert.Equal(t, 12173, ps[1].TotalScore)
		assert.Equal(t, 9513, ps[2].TotalScore)
		assert.Equal(t, 8421, ps[3].TotalScore)
		assert.Equal(t, 6510, ps[4].TotalScore)
		assert.Equal(t, 4704, ps[5].TotalScore)
		assert.Equal(t, 3997, ps[6].TotalScore)
		assert.Equal(t, 3073, ps[7].TotalScore)
		assert.Equal(t, 2982, ps[8].TotalScore)
		assert.Equal(t, 1421, ps[9].TotalScore)
		assert.Equal(t, 896, ps[10].TotalScore)
		assert.Equal(t, 847, ps[11].TotalScore)
		assert.Equal(t, 602, ps[12].TotalScore)
		assert.Equal(t, 210, ps[13].TotalScore)

	})

	t.Run("Averages", func(t *testing.T) {
		ps := []*PlayerSummary{}
		db.Model(&ps).Where("tournament_id = ?", tm.ID).Order("total_score DESC").Select()
		if !assert.Equal(t, len(tm.Players), len(ps)) {
			return
		}

		assert.Equal(t, 4874, ps[0].SkillScore)
		assert.Equal(t, 4057, ps[1].SkillScore)
		assert.Equal(t, 3171, ps[2].SkillScore)
		assert.Equal(t, 2807, ps[3].SkillScore)
		assert.Equal(t, 2170, ps[4].SkillScore)
		assert.Equal(t, 2352, ps[5].SkillScore)
		assert.Equal(t, 1998, ps[6].SkillScore)
		assert.Equal(t, 1536, ps[7].SkillScore)
		assert.Equal(t, 1491, ps[8].SkillScore)
		assert.Equal(t, 1421, ps[9].SkillScore)
		assert.Equal(t, 896, ps[10].SkillScore)
		assert.Equal(t, 847, ps[11].SkillScore)
		assert.Equal(t, 602, ps[12].SkillScore)
		assert.Equal(t, 210, ps[13].SkillScore)

	})
}
