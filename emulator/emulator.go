package emulator

import (
	"errors"
	"image"
	"image/color"
	"io"
	"math/rand/v2"
	"os"
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

	// Draw draws the display buffer to an *image.RGBA.
	Draw(image *image.RGBA)
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

// -- *chip8 public implementation

// LoadFile loads a ROM file into memory.
func (c *chip8) LoadFile(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	return c.LoadBuffer(fd)
}

// LoadBuffer takes ROM data and puts it into memory
func (c *chip8) LoadBuffer(buf io.Reader) error {
	rom, err := io.ReadAll(buf)
	if err != nil {
		return err
	}

	// Check that the ROM size does not exceed available memory
	if len(rom) > int(ROM_START+uint16(len(c.mem))) {
		return errors.New("ROM is too large to fit in memory")
	}

	// Copy the ROM data into memory starting at rom_start
	copy(c.mem[ROM_START:], rom)

	return nil
}

// Cycle runs one emulation cycle.
func (c *chip8) Cycle() error {
	err := c.dispatch(c.opcode())
	if err != nil {
		panic(err)
	}
	c.pc += 2
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

func (c *chip8) Draw(image *image.RGBA) {
	for idx, pixel := range c.display {
		image.Set(idx%DISPLAY_WIDTH, idx/DISPLAY_WIDTH, colorFor(pixel))
	}
}

func colorFor(pixel bool) color.Color {
	if pixel {
		return color.White
	}

	return color.Black
}

func baseChip8(clockSpeed time.Duration) *chip8 {
	c := &chip8{}

	// Bring font data into memory
	c.loadFont()

	// Clock speeds varied over the years and different games
	// expect different system clocks
	c.clock = time.NewTicker(clockSpeed)

	// The timers decrement indepently of the system clock
	c.timerClock = time.NewTicker(TIMER_SPEED)

	c.rng = rand.New(rand.NewPCG(uint64(time.Now().Unix()), 0))

	// Traditionally there would be a bootloader here that sets this
	c.pc = ROM_START

	// Blank out all keypad bits
	c.keypad = keypad(0)

	return c
}
