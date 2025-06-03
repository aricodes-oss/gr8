package emulator

import (
	"bytes"
	"gr8/roms"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFontLoads(t *testing.T) {
	assert := assert.New(t)
	emu, _ := NewEmulatorFromBuf(bytes.NewReader(roms.Chip8Logo), DEFAULT_CLOCK_SPEED)
	c := emu.(*chip8)

	assert.ElementsMatch(c.mem[font_start:font_start+len(FONT)], FONT)
}
