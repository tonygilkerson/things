package sprinkler

import (
	"log"
	"time"
	"tinygo.org/x/drivers"
)

// Device wraps an I2C connection to a Soil device.
type Device struct {
	bus            drivers.I2C
	address        uint16
	moistureBuf    [2]uint8
	temperatureBuf [4]uint8
}

/*
  RELAY control Reg 0x11

  Bit	Desc	                R/W
  --- --------------------- ----
  7	  LED1 / 1: ON 0:OFF	  R/W
  6	  LED2 / 1: ON 0:OFF	  R/W
  5	  LED3 / 1: ON 0:OFF	  R/W
  4	  LED4 / 1: ON 0:OFF	  R/W
  3	  RELAY1 / 1: ON 0:OFF	R/W
  2	  RELAY2 / 1: ON 0:OFF	R/W
  1	  RELAY3 / 1: ON 0:OFF	R/W
  0	  RELAY4 / 1: ON 0:OFF	R/W

	H-bridge circuit

	   +---------+-------------+
     |         |             |
     |         |             |
		 ^        \  relay1      \  relay3
		 |         |             |
	 [Vin]       +-----[M]-----+
     |         |             |
		 |        \  relay2      \  relay4
     |         |             |
     |         |             |
	   +---------+-------------+

		 * stand by - relay 1,2,3,4  off
		 * turn on water supply - relay 1 & 4 on, 2 & 3 off - 00000000
		 * turn off water supply - relay 1 & 4 off, 2 & 3 on

		 01234567
		 --------
		 10011001 - led1, led4, relay1, relay4 = on, 153 dec, 0x99
		 01100110 - led2, led3, relay2, relay3 = off, 102 dec, 0x66

*/

const relayControlRegister = 0x11
const SPRINKLER_ON_CMD = 0x99
const SPRINKLER_OFF_CMD = 0x66
const SPRINKLER_STANDBY_CMD = 0x00

// var onCommand = []uint8{0x11, 0x99}
// var offCommand = []uint8{0x11, 0x66}
// var standbyCommand = []uint8{0x11, 0x00}

// New creates a new sprinkler device. The I2C bus must already be configured.
func New(bus drivers.I2C, address uint16) *Device {
	return &Device{
		bus:     bus,
		address: address,
	}
}

func (d *Device) TurnOn(){

	log.Println("sprinkler.TurnOn: Turn on")
	onCmd := []uint8{relayControlRegister,SPRINKLER_ON_CMD}
	err := d.bus.Tx(d.address, onCmd, nil)
	doOrDie(err)
	time.Sleep(time.Millisecond * 500)
	
	log.Println("sprinkler.TurnOn: Standby")
	standbyCmd := []uint8{relayControlRegister,SPRINKLER_STANDBY_CMD}
	err = d.bus.Tx(d.address, standbyCmd, nil)
	doOrDie(err)
	time.Sleep(time.Millisecond * 500)
}

func (d *Device) TurnOff(){

	log.Println("sprinkler.TurnOn: Turn off")
	onCmd := []uint8{relayControlRegister,SPRINKLER_OFF_CMD}
	err := d.bus.Tx(d.address, onCmd, nil)
	doOrDie(err)
	time.Sleep(time.Millisecond * 500)
	
	log.Println("sprinkler.TurnOn: Standby")
	standbyCmd := []uint8{relayControlRegister,SPRINKLER_STANDBY_CMD}
	err = d.bus.Tx(d.address, standbyCmd, nil)
	doOrDie(err)
	time.Sleep(time.Millisecond * 500)
}

func doOrDie(err error) {
	if err != nil {
		log.Panicf("Oops %v", err)
	}
}