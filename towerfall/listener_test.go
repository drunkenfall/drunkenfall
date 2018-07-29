package towerfall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testListener(t *testing.T) (*Listener, *Tournament) {
	s := MockServer()
	db := s.DB

	l, err := NewListener(s.config, db)
	assert.NoError(t, err)

	tm := testTournament(12, s)
	err = tm.StartTournament(nil)
	assert.NoError(t, err)

	return l, tm
}

func TestMessagesAreAddedToMatch(t *testing.T) {
	t.Run("Passing kind", func(t *testing.T) {
		l, tm := testListener(t)
		assert.Equal(t, 0, len(tm.Matches[tm.Current].Messages))

		msg := `{"type":"round_start","data":{"arrows":[[1,0,0],[0,2,0],[0,0,3],[1,1,1]]}}`
		err := l.handle(tm, []byte(msg))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tm.Matches[tm.Current].Messages))
	})

	t.Run("Unknown kind", func(t *testing.T) {
		l, tm := testListener(t)
		assert.Equal(t, 0, len(tm.Matches[tm.Current].Messages))

		// https://open.spotify.com/track/4Uepu89yTm3wGtYFZz04Vf
		msg := `{"type":"black_queen_style","data":{}}`
		err := l.handle(tm, []byte(msg))
		assert.NoError(t, err)
		assert.Equal(t, 1, len(tm.Matches[tm.Current].Messages))
	})
}
