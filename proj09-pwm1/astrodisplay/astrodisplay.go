// This package is for the Astro display
package astrodisplay

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1351"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// Display
type AstroDisplay struct {
	Status  string
	Body    string
	Display ssd1351.Device
}

// New returns a new display
func New(spi machine.SPI, rst, dc, cs, en, rw machine.Pin) AstroDisplay {

	display := ssd1351.New(spi, rst, dc, cs, en, rw)
	display.Configure(ssd1351.Config{
		Width:        128,
		Height:       128,
		RowOffset:    0,
		ColumnOffset: 0,
	})

	display.Command(ssd1351.SET_REMAP_COLORDEPTH)
	display.Data(0x62)
	display.FillScreen(color.RGBA{0, 0, 0, 0})

	return AstroDisplay{
		Status:  "Init",
		Body:    "",
		Display: display,
	}

}

// Write a status line that shows at the top of the display
func (ad *AstroDisplay) WriteStatus() {

	// green := color.RGBA{0, 255, 0, 255}
	red := color.RGBA{0, 0, 255, 255}
	black := color.RGBA{0, 0, 0, 0}
	// black := color.RGBA{55, 55, 55, 55}

	// clear status line
	ad.Display.FillRectangle(0, 0, 128, 14, black)
	ad.Display.DrawFastHLine(0, 127, 14, red)

	tinyfont.WriteLine(&ad.Display, &freemono.Regular9pt7b, 5, 10, ad.Status, red)

}

// Write the body content in middle of screen
func (ad *AstroDisplay) WriteBody() {
	// green := color.RGBA{0, 255, 0, 255}
	red := color.RGBA{0, 0, 255, 255}
	black := color.RGBA{0, 0, 0, 0}
	// black := color.RGBA{111, 111, 111, 111}

	// clear body
	ad.Display.FillRectangle(0, 15, 128, 40, black)

	tinyfont.WriteLine(&ad.Display, &freemono.Regular9pt7b, 5, 30, ad.Body, red)

}
