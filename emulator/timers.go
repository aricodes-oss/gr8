package emulator

func (c *chip8) timerTick() {
	if c.delayTimer > 0 {
		c.delayTimer -= 1
	}

	if c.soundTimer > 0 {
		c.soundTimer -= 1
	}
}
