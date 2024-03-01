package tm1637

const (
	TM1637_CMD1   = 0x40
	TM1637_CMD2   = 0xC0
	TM1637_CMD3   = 0x80
	TM1637_DSP_ON = 0x08

	// The black 4-digit and 6-digit TM1637 modules on Amazon and eBay (from
	// robotdyn.com and diymore.cc) contain a 10nF capacitor with a 10kOhm pullup
	// registor on each of the DIO and CLK lines. The RC time constant is 100
	// microseconds, which forces the bit transition delay to be somewhere in
	// the 100 microsecond range.
	//
	// The blue 4-digit TM1637 modules contain a much lower capacitor (220 or 470
	// picofareds?), so the bit transition delay can be as low as 3-5
	// microseconds.
	//
	// TODO: This should be an adjustable parameter.
	TM1637_DELAY  = uint8(100)
)

// 7-segment characters encoding for 0-9, A-Z, a-z, blank, dash, star
var segments []byte = []byte{
	0x3F, 0x06, 0x5B, 0x4F, 0x66, 0x6D, 0x7D, 0x07, 0x7F, 0x6F,
	0x77, 0x7C, 0x39, 0x5E, 0x79, 0x71, 0x3D, 0x76, 0x06, 0x1E,
	0x76, 0x38, 0x55, 0x54, 0x3F, 0x73, 0x67, 0x50, 0x6D, 0x78,
	0x3E, 0x1C, 0x2A, 0x76, 0x6E, 0x5B, 0x00, 0x40, 0x63}