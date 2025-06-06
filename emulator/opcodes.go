package emulator

import "fmt"

var ErrInvalidOpcode = fmt.Errorf("invalid opcode")

func (c *chip8) dispatch(opcode uint16) error {
	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			c.CLS()
		case 0x00EE:
			c.RET()
		}
	case 0x1000:
		c.JMP()
	case 0x6000:
		c.LD()
	case 0x7000:
		c.ADD()
	case 0xA000:
		c.LDI()
	case 0xD000:
		c.DRW()
	}

	return nil
}

// CLS clears the display.
func (c *chip8) CLS() {
	for i := range c.display {
		c.display[i] = false
	}
}

// RET returns from subroutine.
func (c *chip8) RET() {
	c.pc, c.stack = c.stack[len(c.stack)], c.stack[:len(c.stack)-1]
}

// JMP jumps to nnn.
func (c *chip8) JMP() {
	c.pc = c.nnn()
}

// LD loads nn into Vx.
func (c *chip8) LD() {
	c.v[c.x()] = c.nn()
}

// ADD adds nn to Vx.
func (c *chip8) ADD() {
	c.v[c.x()] += c.nn()
}

// LDI sets i to nnn.
func (c *chip8) LDI() {
	c.i = c.nnn()
}
