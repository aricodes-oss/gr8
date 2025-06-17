package emulator

type Keypad interface {
  Pressed(uint8) bool
  Press(uint8) Keypad
  Release(uint8) Keypad
}

type keypad uint16

// Pressed returns whether or not the specified bit is set.
func (k keypad) Pressed(key uint8) bool {
	return (k>>(key-1))&0x01 == 1
}

// Press returns a new bitmap with the selected bit set.
func (k keypad) Press(key uint8) Keypad {
	return k | 1<<key - 1
}

// Release returns a new bitmap with the selected bit unset.
func (k keypad) Release(key uint8) Keypad {
	return k & ^(1<<key - 1)
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
