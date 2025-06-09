package emulator

import (
	"fmt"
)

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
	case 0x2000:
		c.CALL()
	case 0x3000:
		c.SEVx()
	case 0x4000:
		c.SNEVx()
	case 0x5000:
		c.SEVxVy()
	case 0x6000:
		c.LD()
	case 0x7000:
		c.ADD()
	case 0x8000:
		switch opcode & 0xF00F {
		case 0x8000:
			c.LDVxVy()
		case 0x8001:
			c.ORVxVy()
		case 0x8002:
			c.ANDVxVy()
		case 0x8003:
			c.XORVxVy()
		case 0x8004:
			c.ADDVxVy()
		case 0x8005:
			c.SUBVxVy()
		case 0x8006:
			c.SHRVx()
		case 0x8007:
			c.SUBNVxVy()
		case 0x800E:
			c.SHLVx()
		}
	case 0x9000:
		c.SNEVxVy()
	case 0xA000:
		c.LDI()
	case 0xB000:
		c.JPV()
	case 0xC000:
		c.RNDVx()
	case 0xD000:
		c.DRW()
	case 0xE000:
		switch opcode & 0xF0FF {
		case 0xE09E:
			c.SKPVx()
		case 0xE0A1:
			c.SKNPVx()
		}
	case 0xF000:
		switch opcode & 0xF0FF {
		case 0xF007:
			c.LDVxDT()
		case 0xF00A:
			c.LDVxK()
		case 0xF015:
			c.LDDTVx()
		case 0xF018:
			c.LDSTVx()
		case 0xF01E:
			c.ADDIVx()
		case 0xF029:
			c.LDFVx()
		case 0xF033:
			c.LDBVx()
		case 0xF055:
			c.LDIVx()
		case 0xF065:
			c.LDVxI()
		}
	default:
		return fmt.Errorf("parsing opcodes: %w (%X)", ErrInvalidOpcode, opcode)
	}

	return nil
}

// CLS clears the display.
func (c *chip8) CLS() {
	for i := range DISPLAY_SIZE {
		c.display[i] = false
	}
}

// RET returns from subroutine.
func (c *chip8) RET() {
	c.pc, c.stack = c.stack[len(c.stack)-1]-2, c.stack[:len(c.stack)]
}

// JMP jumps to nnn.
func (c *chip8) JMP() {
	c.pc = c.nnn() - 2
}

// CALL calls subroutine at nnn.
func (c *chip8) CALL() {
	c.stack = append(c.stack, c.pc+2)
	c.pc = c.nnn() - 2
}

// SEVx skips next instruction if Vx == nn.
func (c *chip8) SEVx() {
	if c.vx() == c.nn() {
		c.pc += 2
	}
}

// SNEVx skips next instruction if Vx != nn.
func (c *chip8) SNEVx() {
	if c.vx() != c.nn() {
		c.pc += 2
	}
}

// SEVxVy skips next instruction if Vx == Vy.
func (c *chip8) SEVxVy() {
	if c.vx() == c.vy() {
		c.pc += 2
	}
}

// LD loads nn into Vx.
func (c *chip8) LD() {
	c.v[c.x()] = c.nn()
}

// ADD adds nn to Vx.
func (c *chip8) ADD() {
	c.v[c.x()] += c.nn()
}

// LDVxVy sets Vx = Vy.
func (c *chip8) LDVxVy() {
	c.v[c.y()] = c.vx()
}

// ORVxVy sets Vx |= Vy.
func (c *chip8) ORVxVy() {
	c.v[c.x()] |= c.vy()
}

// ANDVxVy sets Vx &= Vy.
func (c *chip8) ANDVxVy() {
	c.v[c.x()] &= c.vy()
}

// XORVxVy sets Vx ^= Vy.
func (c *chip8) XORVxVy() {
	c.v[c.x()] ^= c.vy()
}

// ADDVxVy sets Vx += Vy, sets VF on carry.
func (c *chip8) ADDVxVy() {
	result := int(c.vx()) + int(c.vy())
	c.v[c.x()] = byte(result)
	if result > 255 {
		c.v[0xF] = 1
	}
}

// SUBVxVy sets Vx -= Vy, sets VF if NOT borrow.
func (c *chip8) SUBVxVy() {
	c.v[0xF] = 0
	if c.vx() > c.vy() {
		c.v[0xF] = 1
	}

	c.v[c.x()] = max(0, c.vx()-c.vy())
}

// SHRVx sets Vx=Vx>>1, sets VF if least-significant bit is 1.
func (c *chip8) SHRVx() {
	c.v[0xF] = c.vx() & 0x01
	c.v[c.x()] = c.vx() >> 1
}

// SUBNVxVy sets Vx = Vy - Vx, sets VF if NOT carry.
func (c *chip8) SUBNVxVy() {
	c.v[0xF] = 0
	if c.vy() > c.vx() {
		c.v[0xF] = 1
	}

	c.v[c.x()] = c.vy() - c.vx()
}

// SHLVx sets Vx=Vx<<1, sets VF if most-significant bit is 1.
func (c *chip8) SHLVx() {
	c.v[0xF] = c.vx() & 0x80
	c.v[c.x()] = c.vx() << 1
}

// SNEVxVy skips next instruction if Vx != Vy.
func (c *chip8) SNEVxVy() {
	if c.vx() != c.vy() {
		c.pc += 2
	}
}

// LDI sets i to nnn.
func (c *chip8) LDI() {
	c.i = c.nnn()
}

// JPV jumps to location nnn + V0.
func (c *chip8) JPV() {
	c.pc = c.nnn() + uint16(c.v[0]) - 2
}

// RNDVx sets Vx = (random) & nn.
func (c *chip8) RNDVx() {
	c.v[c.x()] = uint8(c.rng.Uint32()) & c.nn()
}

// DRW draws n-byte sprite starting at memory location i at (Vx, Vy)
// and sets VF on collision.
func (c *chip8) DRW() {
	// Coordinates of the sprite on screen
	x, y := c.vx(), c.vy()

	// Clear collision flag
	c.v[0xF] = 0

	// Iterate over the bytes of the sprite
	for offset := range c.n() {
		row := c.mem[c.i+uint16(offset)]

		// Draw each bit in the row to the screen
		for col := range 8 {
			// Convert coordinate to linear index
			index := (int(y+offset) * DISPLAY_WIDTH) + int(x) + col
			spritePixel := row&(0x80>>col) != 0

			// If we try to draw past the right edge or out of bounds,
			// stop processing this row
			if index >= DISPLAY_SIZE {
				continue
			}

			displayPixel := c.display[index]
			if displayPixel && spritePixel {
				c.v[0xF] = 1
			}
			c.display[index] = displayPixel != spritePixel
		}
	}
}

// SKPVx skips next instruction if key with the value of Vx is pressed.
func (c *chip8) SKPVx() {
	if c.keypad.Pressed(c.vx()) {
		c.pc += 2
	}
}

// SKNPVx skips next instruction if key with the value of Vx is NOT pressed.
func (c *chip8) SKNPVx() {
	if !c.keypad.Pressed(c.vx()) {
		c.pc += 2
	}
}

// LDVxDT sets Vx = delay timer.
func (c *chip8) LDVxDT() {
	c.v[c.x()] = c.delayTimer
}

// LDVxK waits for a keypress and stores the value of the key in Vx.
func (c *chip8) LDVxK() {
	for key := range uint8(16) {
		if c.keypad.Pressed(key + 1) {
			c.v[c.x()] = key + 1
			return
		}
	}

	// Check again on the next cycle if a keypress was not found
	c.pc -= 2
}

// LDDTVx sets delay timer = Vx.
func (c *chip8) LDDTVx() {
	c.delayTimer = c.vx()
}

// LDSTVx sets sound timer = Vx.
func (c *chip8) LDSTVx() {
	c.soundTimer = c.vx()
}

// ADDIVx sets i += Vx.
func (c *chip8) ADDIVx() {
	c.i += uint16(c.vx())
}

// LDFVx sets i = address for sprite to digit Vx.
func (c *chip8) LDFVx() {
	c.i = uint16(FONT_START + c.vx()*5)
}

// LDBVx stores the decimal digits of Vx in memory locations [i:i+2] (BCD).
func (c *chip8) LDBVx() {
	copy(c.mem[c.i:], []byte{
		c.vx() / 100,
		(c.vx() % 100) / 10,
		c.vx() % 10,
	})
}

// LDIVx stores registers V0-Vx in memory starting at i.
func (c *chip8) LDIVx() {
	copy(c.mem[c.i:], c.v[:c.x()+1])
}

// LDVxI stores memory starting at i into register V0-Vx.
func (c *chip8) LDVxI() {
	copy(c.v[:], c.mem[c.i:c.i+uint16(c.x()+1)])
}
