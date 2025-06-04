package emulator

import (
	"encoding/binary"
	"time"
)

const rom_start = 0x200

type chip8 struct {
	// 4kb of Internal memory
	mem [4 * 1024]byte

	// Monochromatic display buffer
	display [64 * 32]bool

	// Program counter/instruction pointer
	pc uint16

	// Index register
	i uint16

	// Stack
	stack []uint16

	// Timers, decremented every 60hz independently of the system clock
	delayTimer byte
	soundTimer byte
	timerClock *time.Ticker

	// General-purpose variable registers
	v [16]byte

	// Clock signal, typically set at 60hz
	clock *time.Ticker

	// Shutdown channel, for asynchronous operation
	done chan bool
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
