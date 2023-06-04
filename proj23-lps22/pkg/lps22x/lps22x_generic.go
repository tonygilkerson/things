//go:build !nano_33_ble

package lps22x

// Configure sets up the lps22x device for communication.
func (d *Device) Configure() {
	// set to block update mode
	d.bus.WriteRegister(d.Address, lps22x_CTRL1_REG, []byte{0x02})
}
