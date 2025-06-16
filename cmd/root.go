/*
Copyright © 2025 Aria Taylor <ari@aricodes.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"

	"gr8/emulator"

	"github.com/spf13/cobra"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

var Scale int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gr8",
	Short: "A simple CHIP-8 emulator in Go",
	Long:  `A simple CHIP-8 emulator in Go.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		chip8, err := emulator.NewEmulator(file, emulator.DEFAULT_CLOCK_SPEED)
		if err != nil {
			return err
		}

		cfg := opengl.WindowConfig{
			Title: file,
			Bounds: pixel.R(
				0,
				0,
				float64(emulator.DISPLAY_WIDTH*Scale),
				float64(emulator.DISPLAY_HEIGHT*Scale),
			),
			VSync: true,
		}

		win, err := opengl.NewWindow(cfg)
		if err != nil {
			return err
		}

		go chip8.Run()
		defer chip8.Stop()

		for !win.Closed() {
			frame := chip8.Frame()
			if frame == nil {
				win.Update()
				continue
			}

			pictureData := pixel.PictureDataFromImage(frame)
			texture := pixel.NewSprite(pictureData, pictureData.Rect)
			texture.Draw(win, pixel.IM.Scaled(pixel.ZV, float64(Scale)).Moved(win.Bounds().Center()))

			win.Update()
		}

		return nil
	},
}

func init() {
	rootCmd.Flags().IntVarP(&Scale, "scale", "s", 16, "screen scaling factor")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
