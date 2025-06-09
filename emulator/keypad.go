package emulator

type keypad uint16

// Pressed returns whether or not the specified bit is set.
func (k keypad) Pressed(key uint8) bool {
	return (k>>(key-1))&0x01 == 1
}

// Press returns a new bitmap with the selected bit set.
func (k keypad) Press(key uint8) keypad {
	return k | 1<<key - 1
}

// Release returns a new bitmap with the selected bit unset.
func (k keypad) Release(key uint8) keypad {
	return k & ^(1<<key - 1)
}
