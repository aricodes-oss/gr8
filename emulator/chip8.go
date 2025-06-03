package emulator

import (
	"errors"
	"io"
	"os"
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
	v []byte

	// Clock signal, typically set at 60hz
	clock *time.Ticker

	// Shutdown channel, for asynchronous operation
	done chan bool
}

// LoadBuffer takes ROM data and puts it into memory
func (c *chip8) LoadBuffer(buf io.Reader) error {
	rom, err := io.ReadAll(buf)
	if err != nil {
		return err
	}

	// Check that the ROM size does not exceed available memory
	if len(rom) > int(rom_start+len(c.mem)) {
		return errors.New("ROM is too large to fit in memory")
	}

	// Copy the ROM data into memory starting at rom_start
	copy(c.mem[rom_start:], rom)

	return nil
}

// LoadFile loads a ROM file into memory.
func (c *chip8) LoadFile(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return c.LoadBuffer(fd)
}

// Cycle runs one emulation cycle.
func (c *chip8) Cycle() error {
	return nil
}

// Run runs the emulator.
func (c *chip8) Run() {
	// Create a new signal channel, in case the old one was closed
	c.done = make(chan bool)

	// Decrement the timers separately from the system clock
	go func() {
		for {
			select {
			case <-c.timerClock.C:
				c.timerTick()
			case <-c.done:
				return
			}
		}
	}()

	// Run Cycle on system clock tick, or exit if stopped
	for {
		select {
		case <-c.clock.C:
			c.Cycle()
		case <-c.done:
			return
		}
	}
}

// Stop stops the background emulation process
func (c *chip8) Stop() {
	close(c.done)
}
