// This package is for the Astro display
package astrodisplay

import (
	"image/color"
	"time"

	"tinygo.org/x/drivers/ssd1351"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

// Display
type AstroDisplay struct {
	status  string
	body    string
	display ssd1351.Device
}

// New returns a new display
// func New(spi *machine.SPI, rst, dc, cs, en, rw machine.Pin) AstroDisplay {
func New(display ssd1351.Device) AstroDisplay {

	// display := ssd1351.New(spi, rst, dc, cs, en, rw)
	// display.Configure(ssd1351.Config{
	// 	Width:        128,
	// 	Height:       128,
	// 	RowOffset:    0,
	// 	ColumnOffset: 0,
	// })

	display.Command(ssd1351.SET_REMAP_COLORDEPTH)
	display.Data(0x62)
	display.FillScreen(color.RGBA{1, 1, 1, 1})
	time.Sleep(time.Second * 1)
	display.FillScreen(color.RGBA{2, 2, 2, 2})
	time.Sleep(time.Second * 1)
	display.FillScreen(color.RGBA{3, 3, 3, 3})
	time.Sleep(time.Second * 1)
	display.FillScreen(color.RGBA{0, 0, 0, 0})
	time.Sleep(time.Second * 1)

	return AstroDisplay{
		status:  "Init",
		body:    "",
		display: display,
	}

}

// configure
func (ad *AstroDisplay) Configure() {

}

// set status
func (ad *AstroDisplay) SetStatus(status string) {
	ad.status = status
}

// set body
func (ad *AstroDisplay) SetBody(body string) {
	ad.body = body
}

// Write a status line that shows at the top of the display
func (ad *AstroDisplay) WriteStatus() {

	// green := color.RGBA{0, 255, 0, 255}
	red := color.RGBA{0, 0, 255, 255}
	black := color.RGBA{0, 0, 0, 0}
	// black := color.RGBA{55, 55, 55, 55}

	// clear status line
	ad.display.FillRectangle(0, 0, 128, 14, black)
	ad.display.DrawFastHLine(0, 127, 14, red)

	// tinyfont.WriteLine(&ad.Display, &freemono.Regular9pt7b, 5, 10, ad.Status, red)
	tinyfont.WriteLine(&ad.display, &freemono.Regular9pt7b, 5, 10, ad.status, red)
}

// Write the body content in middle of screen
func (ad *AstroDisplay) WriteBody() {
	// green := color.RGBA{0, 255, 0, 255}
	red := color.RGBA{0, 0, 255, 255}
	black := color.RGBA{0, 0, 0, 0}
	// black := color.RGBA{111, 111, 111, 111}

	// clear body
	ad.display.FillRectangle(0, 15, 128, 40, black)

	// tinyfont.WriteLine(&ad.Display, &freemono.Regular9pt7b, 5, 30, ad.Body, red)
	body := ad.body + "\nxxx"
	// tinyfont.WriteLine(&ad.Display, &tinyfont.Org01, 5, 30, ad.Body, red)
	tinyfont.WriteLine(&ad.display, &tinyfont.Org01, 5, 30, body, red)

	// tinyfont.WriteLine(&ad.Display, &tinyfont.LineFeed, 5, 50, ad.Body, red)

}
