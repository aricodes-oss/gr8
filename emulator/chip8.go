package emulator

import (
	"encoding/binary"
	"image"
	"math/rand/v2"
	"time"

	"github.com/gammazero/deque"
)

const DISPLAY_WIDTH = 64
const DISPLAY_HEIGHT = 32
const DISPLAY_SIZE = DISPLAY_WIDTH * DISPLAY_HEIGHT

const DEFAULT_CLOCK_SPEED = 16 * time.Millisecond // 60hz
const DEFAULT_IPF = 700                           // 700 instructions per frame

const MEM_SIZE = 4 * 1024 // 4kb
const ROM_START = uint16(0x200)

type chip8 struct {
	// 4kb of Internal memory
	mem [MEM_SIZE]byte

	// Monochromatic display buffer
	display [DISPLAY_SIZE]bool

	// Program counter/instruction pointer
	pc uint16

	// Index register
	i uint16

	// Stack
	stack []uint16

	// Timers, decremented every 60hz
	delayTimer byte
	soundTimer byte

	// General-purpose variable registers
	v [16]byte

	// Clock signal, typically set at 60fps
	clock *time.Ticker

	// Instructions to process per frame
	ipf int

	// Live keypad state and snapshot (16 keys)
	keypad,
  frameKeys keypad



	// Shutdown channel, for asynchronous operation
	done chan bool

	// RNG generator
	rng *rand.Rand

	// Frame buffer
	frameBuf deque.Deque[*image.RGBA]
}

// opcode returns the full 2-byte instruction
func (c *chip8) opcode() uint16 {
	return binary.BigEndian.Uint16(c.mem[c.pc : c.pc+2])
}

// The second opcode nibble
func (c *chip8) x() uint8 {
	return uint8((c.opcode() & 0x0F00) >> 8)
}

// Register V[x]
func (c *chip8) vx() uint8 {
	return c.v[c.x()]
}

// The third opcode nibble
func (c *chip8) y() uint8 {
	return uint8((c.opcode() & 0x00F0) >> 4)
}

// Register V[y]
func (c *chip8) vy() uint8 {
	return c.v[c.y()]
}

// Register V[0xF]
func (c *chip8) vf() uint8 {
	return c.v[0xF]
}

// The 4th nibble
func (c *chip8) n() uint8 {
	return uint8(c.opcode() & 0x000F)
}

// Second byte (3rd+4th nibble)
func (c *chip8) nn() uint8 {
	return uint8(c.opcode() & 0x00FF)
}

// 2nd+3rd+4th nibble
func (c *chip8) nnn() uint16 {
	return uint16(c.opcode() & 0x0FFF)
}
