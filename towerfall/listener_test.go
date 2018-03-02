package towerfall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKillMessage(t *testing.T) {
	s := MockServer()
	db := s.DB
	l, err := NewListener(db, 12345)
	assert.NoError(t, err)

	t.Run("Normal handling", func(t *testing.T) {
		msg := `{"type":"kill","data":{"killer":1,"player":2,"cause":0}}`
		l.handle(msg)
	})
}

func TestStartRound(t *testing.T) {
	s := MockServer()
	db := s.DB
	l, err := NewListener(db, 12345)
	assert.NoError(t, err)
	tm := testTournament(12, s)
	err = tm.StartTournament(nil)
	assert.NoError(t, err)

	t.Run("Normal handling", func(t *testing.T) {
		msg := `{"type":"round_start","data":{"arrows":[[0,0,0],[0,0,0]]}}`
		err := l.handle(msg)
		assert.NoError(t, err)
	})
}
