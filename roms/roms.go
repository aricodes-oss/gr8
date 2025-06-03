// Package roms provides test roms for use in validating the emulator.
// Thanks to https://github.com/Timendus/chip8-test-suite for the ROM files
package roms

import _ "embed"

//go:embed 1-chip8-logo.ch8
var Chip8Logo []byte

//go:embed 2-ibm-logo.ch8
var IBMLogo []byte

//go:embed 3-corax+.ch8
var Corax []byte

//go:embed 4-flags.ch8
var Flags []byte

//go:embed 5-quirks.ch8
var Quirks []byte

//go:embed 6-keypad.ch8
var Keypad []byte

//go:embed 7-beep.ch8
var Beep []byte

//go:embed 8-scrolling.ch8
var Scrolling []byte
