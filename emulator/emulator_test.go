package emulator

import (
	"bytes"
	"gr8/roms"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var rom = roms.Chip8Logo

func TestRomLoads(t *testing.T) {
	c, assert := setup(t)
	assert.ElementsMatch(c.mem[ROM_START:int(ROM_START)+len(rom)], roms.Chip8Logo)
}

func TestTimerDecrements(t *testing.T) {
	c, assert := setup(t)

	// Set the timers
	initial := uint8(60)
	expected := initial - 1
	c.delayTimer = initial
	c.soundTimer = initial

	// Start the emulator in the background
	go c.Run()

	// Sleep for one full timer tick then stop emulation
	time.Sleep(DEFAULT_CLOCK_SPEED + 5*time.Millisecond)
	c.Stop()

	// Check to make sure the timers decremented properly
	assert.Equal(c.delayTimer, expected)
	assert.Equal(c.soundTimer, expected)
}

func setup(t *testing.T) (*chip8, *assert.Assertions) {
	assert := assert.New(t)
	emu, err := NewEmulatorFromBuf(bytes.NewReader(rom), DEFAULT_CLOCK_SPEED)
	if err != nil {
		t.Fatal(err)
		return nil, assert
	}
	c := emu.(*chip8)
	return c, assert
}
