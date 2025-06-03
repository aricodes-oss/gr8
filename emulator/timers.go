package emulator

import "time"

// The speed our delay and sound timers count down at, independent from
// the clock cycle.
const TIMER_SPEED = 16 * time.Millisecond

func (c *chip8) timerTick() {
	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}

	if c.soundTimer > 0 {
		c.soundTimer -= 1
	}
}
