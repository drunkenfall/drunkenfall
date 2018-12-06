package towerfall

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinalMultiplier(t *testing.T) {
	assert.Equal(t, 2.5, FinalMultiplier(15))
	assert.Equal(t, 2.5, FinalMultiplier(16))

	assert.Equal(t, 2.8940625000000004, FinalMultiplier(19))
	assert.Equal(t, 3.3502391015625, FinalMultiplier(22))
	assert.Equal(t, 3.8783205399462894, FinalMultiplier(25))
	assert.Equal(t, 4.489640815055323, FinalMultiplier(28))
}
