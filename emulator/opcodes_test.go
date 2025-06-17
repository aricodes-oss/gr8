package emulator

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func opcodeTest(t *testing.T, instructions []byte) (*chip8, *assert.Assertions) {
	c, assert := setup(t)
	c.LoadBuffer(bytes.NewReader(instructions))
	return c, assert
}

// 0x00E0
func TestCLS(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x00, 0xE0})
	// Turn on all pixels
	for i := range DISPLAY_SIZE {
		c.display[i] = true
	}

	assert.Equal(true, c.display[0])
	c.Cycle()

	// Ensure all pixels are cleared
	assert.ElementsMatch(make([]bool, DISPLAY_SIZE), c.display)
}

// 0x00EE
func TestRET(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x00, 0xEE})
	newAddr := uint16(0xBEEF)

	c.stack = append(c.stack, newAddr)
	assert.Equal(ROM_START, c.pc)
	assert.Greater(len(c.stack), 0)

	c.Cycle()
	assert.Equal(newAddr, c.pc)
	assert.Equal(len(c.stack), 0)
}

// 0x1nnn
func TestJMP(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x1B, 0xEE})
	newAddr := c.nnn()

	assert.Equal(uint16(ROM_START), c.pc)
	assert.NotEqual(newAddr, c.pc)
	c.Cycle()
	assert.Equal(newAddr, c.pc)
}

// 0x2nnn
func TestCALL(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x2B, 0xEE})
	newAddr := c.nnn()

	assert.Equal(0, len(c.stack))
	c.Cycle()
	assert.Equal(ROM_START+2, c.stack[0])
	assert.Equal(newAddr, c.pc)
	assert.Equal(1, len(c.stack))
}

// 0x3xnn
func TestSEVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x30, 0x22})
	c.v[c.x()] = 0x22

	c.Cycle()
	assert.Equal(ROM_START+4, c.pc)
}

// 0x4xnn
func TestSNEVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x40, 0x22})
	c.v[c.x()] = 0x33

	c.Cycle()
	assert.Equal(ROM_START+4, c.pc)
}

// 0x5xy0
func TestSEVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x50, 0x10})
	c.v[c.x()] = 0x33
	c.v[c.y()] = 0x33

	c.Cycle()
	assert.Equal(ROM_START+4, c.pc)
}

// 0x6xnn
func TestLD(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x60, 0xBE})

	assert.Equal(uint8(0), c.v[0])
	c.Cycle()
	assert.Equal(uint8(0xBE), c.v[0])
}

// 0x7xnn
func TestADD(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x70, 0xBE})

	assert.Equal(uint8(0), c.v[0])
	c.Cycle()
	assert.Equal(uint8(0xBE), c.v[0])
}

// 0x8xy0
func TestLDVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x10})
	c.v[c.x()] = 255

	initialX := c.vy()

	assert.NotEqual(c.vx(), c.vy())
	c.Cycle()
	assert.Equal(c.v[0], c.v[1])

	// Check that X was copied to Y
	assert.Equal(c.v[0], initialX)
}

// 0x8xy1
func TestORVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x11})
	c.v[c.x()] = 0b01010101
	c.v[c.y()] = 0b10101010

	assert.NotEqual(0xFF, c.vx())
	c.Cycle()
	assert.Equal(uint8(0xFF), c.v[0])
}

// 0x8xy2
func TestANDVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x12})
	c.v[c.x()] = 0b01010101
	c.v[c.y()] = 0b10101010

	assert.NotEqual(0xFF, c.vx())
	c.Cycle()
	assert.Equal(uint8(0x00), c.v[0])
}

// 0x8xy3
func TestXORVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x13})
	c.v[c.x()] = 0b01010101
	c.v[c.y()] = 0b10101010

	assert.NotEqual(0xFF, c.vx())
	c.Cycle()
	assert.Equal(uint8(0xFF), c.v[0])
}

// 0x8xy4
func TestADDVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x14})
	c.v[c.x()] = 0b01010101
	c.v[c.y()] = 0b10101010

	assert.NotEqual(0xFF, c.vx())
	c.Cycle()
	assert.Equal(uint8(0xFF), c.v[0])
	assert.Equal(uint8(0), c.vf())
}

// 0x8xy5
func TestSUBVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x15})
	c.v[c.x()] = 0b10101010
	c.v[c.y()] = 0b01010101

	assert.NotEqual(0xFF, c.vx())
	c.Cycle()
	assert.Equal(uint8(0x55), c.v[0])
	assert.Equal(uint8(0x1), c.vf())
}

// 0x8xy6
func TestSHRVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x16})
	var initial uint8 = 0b01010101
	c.v[c.x()] = initial

	assert.Equal(uint8(0), c.vf())
	c.Cycle()

	assert.Equal(initial>>1, c.v[0])
}

// 0x8xy7
func TestSUBNVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x17})
	c.v[c.x()] = 0b01010101
	c.v[c.y()] = 0b10101010
	expected := c.vy() - c.vx()

	c.Cycle()
	assert.Equal(uint8(0x1), c.vf())
	assert.Equal(uint8(expected), c.v[0])
}

// 0x8xyE
func TestSHLVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x80, 0x1E})
	var initial uint8 = 0b01010101
	c.v[c.x()] = initial

	assert.Equal(uint8(0), c.vf())
	c.Cycle()
	assert.Equal(uint8(0), c.vf())
	assert.Equal(initial<<1, c.v[0])
}

// 0x9xy0
func TestSNEVxVy(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0x90, 0x10})
	c.v[c.x()] = 0
	c.v[c.y()] = 1

	c.Cycle()
	assert.Equal(ROM_START+4, c.pc)
}

// 0xAnnn
func TestLDI(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xA1, 0x23})

	assert.Equal(c.i, uint16(0))
	c.Cycle()
	assert.Equal(c.i, uint16(0x0123))
}

// 0xBnnn
func TestJPV(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xB1, 0x23})
	c.v[0] = 0x0001
	expected := uint16(c.v[0]) + c.nnn()

	assert.Equal(ROM_START, c.pc)
	c.Cycle()
	assert.Equal(expected, c.pc)
}

// 0xCxnn
func TestRNDVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xC0, 0x23})
	rng := rand.New(rand.NewPCG(0, 1))
	c.rng = rand.New(rand.NewPCG(0, 1))

	expected := uint8(rng.Uint32()) & c.nn()
	c.Cycle()
	assert.Equal(expected, c.v[0])
}

// 0xDxyn
func TestDRW(t *testing.T) {
	var trueElements = func(l []bool) int {
		total := 0

		for _, el := range l {
			if el {
				total += 1
			}
		}

		return total
	}

	spriteSize := 2
	sprite := make([]byte, spriteSize)
	for idx := range sprite {
		sprite[idx] = 0xFF
	}
	c, assert := opcodeTest(t, append([]byte{0xD0, 0x02}, sprite...))
	c.i = ROM_START + 2

	expected := make([]bool, spriteSize*8)
	for idx := range expected {
		expected[idx] = true
	}
	expected = append(expected, make([]bool, DISPLAY_SIZE-(spriteSize*8))...)

	assert.ElementsMatch(make([]bool, DISPLAY_SIZE), c.display)
	c.Cycle()
	assert.ElementsMatch(expected, c.display)
	fmt.Println(trueElements(expected))
}

// 0xEx9E
func TestSKPVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xE0, 0x9E})

	c.Press(1)
	c.v[c.x()] = 1

	assert.Equal(ROM_START, c.pc)
	c.Cycle()
	assert.Equal(ROM_START+4, c.pc)
}

// 0xExA1
func TestSKNPVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xE0, 0xA1})

	c.Press(1)
	c.v[c.x()] = 1

	assert.Equal(ROM_START, c.pc)
	c.Cycle()
	assert.Equal(ROM_START+2, c.pc)
}

// 0xFx07
func TestLDVxDT(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x07})
	c.delayTimer = 42

	assert.NotEqual(c.vx(), c.delayTimer)
	c.Cycle()
	assert.Equal(c.v[0], c.delayTimer)
}

// 0xFx0A
func TestLDVxK(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x0A})

	// At ROM_START, no keys pressed, Vx empty
	assert.Equal(uint8(0), c.vx())
	assert.Equal(ROM_START, c.pc)
	assert.Equal(keypad(0), c.keypad)
	c.Cycle()

	// Same check again
	assert.Equal(uint8(0), c.vx())
	assert.Equal(ROM_START, c.pc)

	// Press a key, try again
	c.Press(2)
	c.Cycle()

	// Past the first instruction, V0 == 1
	assert.Equal(ROM_START+2, c.pc)
	assert.Equal(uint8(1), c.v[0])
}

// 0xFx15
func TestLDDTVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x15})
	c.delayTimer = 42

	assert.NotEqual(c.vx(), c.delayTimer)
	c.Cycle()
	assert.Equal(c.vx(), c.delayTimer)
}

// 0xFx18
func TestLDSTVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x18})
	c.soundTimer = 42

	assert.NotEqual(c.vx(), c.soundTimer)
	c.Cycle()
	assert.Equal(c.vx(), c.soundTimer)
}

// 0xFx1E
func TestADDIVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x1E})
	c.v[c.x()] = 255

	c.Cycle()
	assert.Equal(uint16(c.v[0]), c.i)
}

// 0xFx29
func TestLDFVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x29})
	c.v[c.x()] = 2
	assert.Equal(c.i, uint16(0))

	c.Cycle()
	assert.Equal(c.i, uint16(FONT_START+(2*5)))
}

// 0xFx33
func TestLDBVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF0, 0x33})
	c.v[c.x()] = 123

	c.Cycle()
	assert.ElementsMatch(c.mem[0:3], []byte{1, 2, 3})
}

// 0xFx55
func TestLDIVx(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF2, 0x55})
	c.i = ROM_START + 16
	copy(c.v[:], []byte{1, 2, 3})

	c.Cycle()
	assert.ElementsMatch(c.mem[c.i:c.i+3], c.v[:3])
}

// 0xFx65
func TestLDVxI(t *testing.T) {
	c, assert := opcodeTest(t, []byte{0xF2, 0x65})
	c.i = ROM_START + 16
	copy(c.mem[c.i:], []byte{1, 2, 3})

	c.Cycle()
	assert.ElementsMatch(c.v[:3], []byte{1, 2, 3})
}
