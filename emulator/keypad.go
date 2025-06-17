package emulator

import (
	"github.com/gopxl/pixel/v2"
)

// TODO: make this configurable
var Keybinds = []pixel.Button{
	pixel.KeyX,

	pixel.Key1,
	pixel.Key2,
	pixel.Key3,

	pixel.KeyQ,
	pixel.KeyW,
	pixel.KeyE,

	pixel.KeyA,
	pixel.KeyS,
	pixel.KeyD,

	pixel.KeyZ,
	pixel.KeyC,

	pixel.Key4,
	pixel.KeyR,
	pixel.KeyF,
	pixel.KeyV,
}

type Keypad interface {
	Pressed(uint8) bool
	Press(uint8) Keypad
	Release(uint8) Keypad
}

type keypad uint16

// Pressed returns whether or not the specified bit is set.
func (k keypad) Pressed(key uint8) bool {
	return (k>>(key))&0x01 == 1
}

// Press returns a new bitmap with the selected bit set.
func (k keypad) Press(key uint8) Keypad {
	return k | (1 << (key))
}

// Release returns a new bitmap with the selected bit unset.
func (k keypad) Release(key uint8) Keypad {
	return k & ^(1 << (key))
}

// Passthrough stubs for the emulator
func (c *chip8) Pressed(key uint8) bool {
	return c.keypad.Pressed(key)
}

func (c *chip8) Press(key uint8) Keypad {
	c.keypad = c.keypad.Press(key).(keypad)
	return c.keypad
}

func (c *chip8) Release(key uint8) Keypad {
	c.keypad = c.keypad.Release(key).(keypad)
	return c.keypad
}
