// This package is for the Astro display
package astrodisplay

import (
	"image/color"
	"machine"

	// "time"

	"tinygo.org/x/drivers/ssd1351"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// DEVTODO - add sync.Mutex to ensure only one client channel is active at a time

// Display
type AstroDisplay struct {
	status     string
	prevStatus string
	body       string
	prevBody   string
	display    *ssd1351.Device
}

// New returns a new display
// func New(spi *machine.SPI, rst, dc, cs, en, rw machine.Pin) AstroDisplay {
// func New(spi SPI, display ssd1351.Device, rst, dc, en, rw, cs machine.Pin) AstroDisplay {
func New(spi *machine.SPI, display ssd1351.Device, rst, dc, en, rw, cs machine.Pin) AstroDisplay {

	display = ssd1351.New(spi, rst, dc, cs, en, rw)
	display.Configure(ssd1351.Config{
		Width:        128,
		Height:       128,
		RowOffset:    0,
		ColumnOffset: 0,
	})

	// display.Command(ssd1351.SET_REMAP_COLORDEPTH)
	// display.Data(0x62)
	display.FillScreen(color.RGBA{0, 0, 0, 0})

	return AstroDisplay{
		status:     "Init",
		prevStatus: "",
		body:       "",
		prevBody:   "",
		display:    &display,
	}

}

// configure
func (ad *AstroDisplay) Configure() {

}

// set status
func (ad *AstroDisplay) SetStatus(status string) {
	ad.prevStatus = ad.status
	ad.status = status
}

// set body
func (ad *AstroDisplay) SetBody(body string) {
	ad.prevBody = ad.body
	ad.body = body
}

// Write a status line that shows at the top of the display
func (ad *AstroDisplay) WriteStatus() {

	// green := color.RGBA{0, 255, 0, 255}
	red := color.RGBA{0, 0, 255, 255}
	//white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 0}

	// clear status line
	//ad.display.FillRectangle(0, 0, 128, 20, white)
	ad.display.DrawFastHLine(0, 127, 14, red)

	tinyfont.WriteLine(ad.display, &freemono.Regular9pt7b, 5, 10, ad.prevStatus, black)
	tinyfont.WriteLine(ad.display, &freemono.Regular9pt7b, 5, 10, ad.status, red)

}

// Write the body content in middle of screen
func (ad *AstroDisplay) WriteBody() {

	// green := color.RGBA{0, 255, 0, 255}
	red := color.RGBA{0, 0, 255, 255}
	//white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 0}

	// clear body
	//ad.display.FillRectangle(0, 15, 128, 40, black)

	// tinyfont.WriteLine(&ad.Display, &freemono.Regular9pt7b, 5, 30, ad.Body, red)
	// tinyfont.WriteLine(&ad.Display, &tinyfont.Org01, 5, 30, ad.Body, red)
	tinyfont.WriteLine(ad.display, &tinyfont.Org01, 5, 30, ad.prevBody, black)
	tinyfont.WriteLine(ad.display, &tinyfont.Org01, 5, 30, ad.body, red)

	// tinyfont.WriteLine(&ad.Display, &tinyfont.LineFeed, 5, 50, ad.Body, red)

}
