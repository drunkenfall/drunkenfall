package towerfall

import (
	"log"
	"testing"
	"time"

	"github.com/drunkenfall/drunkenfall/faking"
	"github.com/stretchr/testify/assert"
)

// People that have been used for a tournament. Used to make sure we
// don't randomly grab one we already have grabbed
var usedPeople []string

// MockMatch makes a mock Match{} with a dummy Tournament{}
func MockMatch(t *testing.T, s *Server, idx int, cat string) (m *Match) {
	tm := testTournament(t, s, minPlayers)
	err := tm.StartTournament(nil)
	if err != nil {
		log.Fatal(err)
	}

	offset := 0

	switch cat {
	case qualifying:
		m = tm.Matches[offset+idx]
	case playoff:
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 1, qualifying)

	p := testPlayer(s)

	err := m.AddPlayer(p)
	assert.NotNil(err)
	assert.Equal(4, len(m.Players))
}

func TestStartAlreadyStartedMatch(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, qualifying)
	m.Started = time.Now()

	err := m.Start(nil)
	assert.NotNil(err)
}

func TestStart(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, qualifying)

	err := m.Start(nil)
	assert.Nil(err)
	assert.Equal(false, m.Started.IsZero())
}

func TestEndGivesShotToWinner(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 1, qualifying)

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

	m := MockMatch(t, s, 1, qualifying)
	m.Ended = time.Now()

	err := m.End(nil)
	assert.NotNil(err)
}

func TestCommitSweepPlayer1(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

	m := MockMatch(t, s, 0, qualifying)
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

func TestCorrectColorConflicts(t *testing.T) {

}

func TestCorrectColorConflictsNoScores(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].getPerson().PreferredColor = "green"
	m.Tournament.Players[1].getPerson().PreferredColor = "green"
	m.Tournament.Players[2].getPerson().PreferredColor = "blue"
	m.Tournament.Players[3].getPerson().PreferredColor = "red"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].getPerson().PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].getPerson().PreferredColor)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].getPerson().PreferredColor)
	assert.Equal("red", m.Players[3].Color)
	assert.Equal("red", m.Players[3].getPerson().PreferredColor)
}

func TestCorrectColorConflictsUserLevels(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].getPerson().PreferredColor = "green"
	m.Tournament.Players[0].getPerson().Userlevel = 10000
	m.Tournament.Players[1].getPerson().PreferredColor = "green"
	m.Tournament.Players[2].getPerson().PreferredColor = "red"
	m.Tournament.Players[2].getPerson().Userlevel = -10000
	m.Tournament.Players[3].getPerson().PreferredColor = "red"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].getPerson().PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].getPerson().PreferredColor)
	assert.NotEqual("red", m.Players[2].Color)
	assert.Equal("red", m.Players[2].getPerson().PreferredColor)
	assert.Equal("red", m.Players[3].Color)
	assert.Equal("red", m.Players[3].getPerson().PreferredColor)

}

func TestCorrectColorConflictsNoScoresDoubleConflict(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].getPerson().PreferredColor = "green"
	m.Tournament.Players[1].getPerson().PreferredColor = "green"
	m.Tournament.Players[2].getPerson().PreferredColor = "blue"
	m.Tournament.Players[3].getPerson().PreferredColor = "blue"

	assert.NoError(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.NoError(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.NoError(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.NoError(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.Equal("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].getPerson().PreferredColor)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].getPerson().PreferredColor)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].getPerson().PreferredColor)
	assert.NotEqual("blue", m.Players[3].Color)
	assert.Equal("blue", m.Players[3].getPerson().PreferredColor)
}

func TestCorrectColorConflictsWithScoresDoubleConflict(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].getPerson().PreferredColor = "green"
	m.Tournament.Players[0].getPerson().Nick = "GreenCorrected"

	m.Tournament.Players[1].TotalScore = 3
	m.Tournament.Players[1].getPerson().PreferredColor = "green"

	m.Tournament.Players[2].getPerson().PreferredColor = "blue"
	m.Tournament.Players[2].getPerson().Nick = "BlueCorrected"

	m.Tournament.Players[3].TotalScore = 3
	m.Tournament.Players[3].getPerson().PreferredColor = "blue"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.NotEqual("green", m.Players[0].Color)
	assert.Equal("green", m.Players[0].getPerson().PreferredColor)
	assert.Equal("green", m.Players[1].Color)
	assert.Equal("green", m.Players[1].getPerson().PreferredColor)
	assert.NotEqual("blue", m.Players[2].Color)
	assert.Equal("blue", m.Players[2].getPerson().PreferredColor)
	assert.Equal("blue", m.Players[3].Color)
	assert.Equal("blue", m.Players[3].getPerson().PreferredColor)

}

func TestCorrectColorConflictsWithScoresTripleConflict(t *testing.T) {
	assert := assert.New(t)

	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, "final")
	m.Players = make([]Player, 0)

	m.Tournament.Players[0].getPerson().PreferredColor = "green"
	m.Tournament.Players[0].getPerson().Nick = "Green1Corrected"

	m.Tournament.Players[1].TotalScore = 3
	m.Tournament.Players[1].getPerson().PreferredColor = "green"
	m.Tournament.Players[1].getPerson().Nick = "Green2Corrected"

	m.Tournament.Players[2].getPerson().PreferredColor = "blue"
	m.Tournament.Players[2].getPerson().Nick = "BlueKeep"

	m.Tournament.Players[3].TotalScore = 10
	m.Tournament.Players[3].getPerson().PreferredColor = "green"
	m.Tournament.Players[3].getPerson().Nick = "GreenKeep"

	assert.Nil(m.AddPlayer(m.Tournament.Players[0].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[1].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[2].Player()))
	assert.Nil(m.AddPlayer(m.Tournament.Players[3].Player()))

	assert.NotEqual("green", m.Players[0].Color)
	assert.NotEqual("green", m.Players[1].Color)
	assert.Equal("blue", m.Players[2].Color)
	assert.Equal("green", m.Players[3].Color)

}

func TestMakeKillOrder(t *testing.T) {
	assert := assert.New(t)
	s, teardown := MockServer(t)
	defer teardown()

	m := MockMatch(t, s, 0, qualifying)

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

		m := MockMatch(t, s, 0, qualifying)

		km := KillMessage{1, 2, rArrow}
		err := m.Kill(km)
		assert.NoError(t, err)
	})

	t.Run("Environment kill", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

		km := KillMessage{1, EnvironmentKill, rExplosion}
		err := m.Kill(km)
		assert.NoError(t, err)
		assert.Equal(t, 1, m.Players[1].Self)
	})

	t.Run("Suicide", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

		km := KillMessage{1, 1, rCurse}
		err := m.Kill(km)
		assert.NoError(t, err)
		assert.Equal(t, 1, m.Players[1].Self)
	})
}

func TestLavaOrb(t *testing.T) {
	t.Run("Enable", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

		lm := LavaOrbMessage{0, true}
		err := m.LavaOrb(lm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Lava)
	})

	t.Run("Disable", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

		lm := LavaOrbMessage{0, false}
		err := m.LavaOrb(lm)
		assert.NoError(t, err)
		assert.Equal(t, false, m.Players[0].State.Lava)
	})
}

func TestShield(t *testing.T) {
	t.Run("Gain", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

		sm := ShieldMessage{0, true}
		err := m.ShieldUpdate(sm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Shield)
	})

	t.Run("Break", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)
		m.Players[0].State.Shield = true

		sm := ShieldMessage{0, false}
		err := m.ShieldUpdate(sm)
		assert.NoError(t, err)
		assert.Equal(t, false, m.Players[0].State.Shield)
	})
}

func TestWings(t *testing.T) {
	t.Run("Grow", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

		wm := WingsMessage{0, true}
		err := m.WingsUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, true, m.Players[0].State.Wings)
	})

	t.Run("Lose", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)
		m.Players[0].State.Wings = true

		wm := WingsMessage{0, false}
		err := m.WingsUpdate(wm)
		assert.NoError(t, err)
		assert.Equal(t, false, m.Players[0].State.Wings)
	})
}

func TestArrow(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		s, teardown := MockServer(t)
		defer teardown()

		m := MockMatch(t, s, 0, qualifying)

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
