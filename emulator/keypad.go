package emulator

type keypad uint16

// Pressed returns whether or not the specified bit is set.
func (k keypad) Pressed(key uint8) bool {
	return (k>>key)&0x01 == 1
}
