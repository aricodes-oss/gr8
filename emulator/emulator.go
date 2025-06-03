package emulator

import (
	"io"
	"time"
)

const DEFAULT_CLOCK_SPEED = 1428 * time.Microsecond // Roughly 700 instructions/second

type Emulator interface {
	// LoadFile loads a ROM file from disk into emulator memory.
	LoadFile(path string) error

	// LoadBuffer loads a ROM file from a buffer into emulator memory.
	LoadBuffer(buf io.Reader) error

	// Cycle runs one emulation cycle.
	Cycle() error

	// Run runs the emulator in the background. Call Stop() to end it.
	Run()

	// Stop stops the background emulation process.
	Stop()
}

// NewEmulator takes a path to a ROM file and returns an Emulator with that ROM loaded.
func NewEmulator(rom_path string, clockSpeed time.Duration) (Emulator, error) {
	c := baseChip8(clockSpeed)
	err := c.LoadFile(rom_path)

	return c, err
}

// NewEmulatorFromBuf takes a ROM buffer and returns an Emulator with that ROM loaded.
func NewEmulatorFromBuf(buf io.Reader, clockSpeed time.Duration) (Emulator, error) {
	c := baseChip8(clockSpeed)
	err := c.LoadBuffer(buf)

	return c, err
}

func baseChip8(clockSpeed time.Duration) *chip8 {
	c := &chip8{}
	c.clock = time.NewTicker(clockSpeed)
	// Separately from the system clock, the timers decrement at 60hz
	c.timerClock = time.NewTicker(TIMER_SPEED)
	c.loadFont()

	return c
}
