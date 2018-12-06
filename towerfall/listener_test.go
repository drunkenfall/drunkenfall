package towerfall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testListener(t *testing.T) (*Listener, *Tournament, func()) {
	s, teardown := MockServer(t)
	db := s.DB

	l, err := NewListener(s.config, db)
	assert.NoError(t, err)

	tm := testTournament(t, s, 12)
	err = tm.StartTournament(nil)
	assert.NoError(t, err)

	return l, tm, teardown
}

// TODO(thiderman): Fix this and re-add it with a test that checks the database rather than just
// talking with the Match object
// func TestMessagesAreAddedToMatch(t *testing.T) {
// 	t.Run("Passing kind", func(t *testing.T) {
// 		l, tm, teardown := testListener(t)
// 		defer teardown()

// 		m, err := tm.CurrentMatch()
// 		assert.NoError(t, err)
// 		assert.Equal(t, 0, len(m.Messages))

// 		msg := `{"type":"round_start","data":{"arrows":[[1,0,0],[0,2,0],[0,0,3],[1,1,1]]}}`
// 		err = l.handle(tm, []byte(msg))
// 		assert.NoError(t, err)

// 		m, err = tm.CurrentMatch()
// 		assert.NoError(t, err)
// 		verifyMessages(t, m, 1)
// 	})

// 	t.Run("Unknown kind", func(t *testing.T) {
// 		l, tm, teardown := testListener(t)
// 		defer teardown()

// 		m, err := tm.CurrentMatch()
// 		assert.NoError(t, err)
// 		assert.Equal(t, 0, len(m.Messages))

// 		// https://open.spotify.com/track/4Uepu89yTm3wGtYFZz04Vf
// 		msg := `{"type":"black_queen_style","data":{}}`
// 		err = l.handle(tm, []byte(msg))
// 		assert.Error(t, err)

// 		m, err = tm.CurrentMatch()
// 		assert.NoError(t, err)
// 		verifyMessages(t, m, 0)
// 	})
// }

func verifyMessages(t *testing.T, m *Match, x int) {
	ret, err := globalDB.MatchMessages(m)
	assert.NoError(t, err)
	assert.Equal(t, x, len(ret))
}
