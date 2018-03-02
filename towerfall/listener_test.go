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
