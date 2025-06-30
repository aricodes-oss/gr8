package emulator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPress(t *testing.T) {
	assert := assert.New(t)
	k := keypad(0x0)

	assert.Equal(keypad(0x0), k)
	k = k.Press(1).(keypad)
	assert.Equal(keypad(2), k)
}

func TestRelease(t *testing.T) {
	assert := assert.New(t)
	k := keypad(0).Press(1)

	assert.Equal(keypad(2), k)
	k = k.Release(1).(keypad)
	assert.Equal(keypad(0x0), k)
}

func TestPressed(t *testing.T) {
	assert := assert.New(t)
	k := keypad(0).Press(1)

	assert.Equal(keypad(2), k)
	assert.True(k.Pressed(1))
}
