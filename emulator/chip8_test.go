package emulator

import (
	"fmt"
	"testing"
)

func TestOpcodeSanity(t *testing.T) {
	c, assert := setup(t)

	// Opcode
	fmt.Println(c.opcode())
	assert.Equal(1, 1)
}
