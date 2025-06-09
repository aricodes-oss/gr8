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

	assert.ElementsMatch(c.mem[FONT_START:FONT_START+len(FONT)], FONT)
}
