package emulator

import (
	"testing"
)

func TestFontLoads(t *testing.T) {
	c, assert := setup(t)

	assert.ElementsMatch(c.mem[FONT_START:FONT_START+len(FONT)], FONT)
}
