package lps22x

const (

	// I2C address
	// lps22x_ADDRESS = 0x5C
	lps22x_ADDRESS = 0x5D

	// control/status registers
	lps22x_WHO_AM_I_REG  = 0x0F
	lps22x_CTRL1_REG     = 0x10
	lps22x_CTRL2_REG     = 0x11
	lps22x_STATUS_REG    = 0x27
	lps22x_PRESS_OUT_REG = 0x28
	lps22x_TEMP_OUT_REG  = 0x2B
)
