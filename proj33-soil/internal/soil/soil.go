package soil
// Note there is now a tinygo driver for this but it is on the dev branch
//      I should check out this driver and if it works for me switch ASAP
//      see: https://github.com/tinygo-org/drivers/blob/dev/examples/seesaw/main.go
//           https://gophers.slack.com/archives/CDJD3SUP6/p1708793230660769?thread_ts=1708486573.676089&cid=CDJD3SUP6
//
// This package taken from:
//
//  @ysoldak Yurii Soldak in the #tinygo channel
//  https://github.com/ysoldak/plantbot/blob/seesaw/src/seesaw/example/example.go
//
// See also:
// https://learn.adafruit.com/adafruit-stemma-soil-sensor-i2c-capacitive-moisture-sensor
// https://www.electrokit.com/en/product/jordfuktighetssensor-kapacitiv-i2c/
//
// https://github.com/adafruit/Adafruit-STEMMA-Soil-Sensor-PCB
//
// https://github.com/adafruit/Adafruit_Seesaw/blob/master/Adafruit_seesaw.cpp
// https://github.com/adafruit/Adafruit_BusIO/blob/master/Adafruit_I2CDevice.cpp
//
// The reading for soil moisture can range from about 300 to 1,000 at the extremes
// In soil, you'll see this range from about 0 to 1023
// Note that it does depend on how packed/loose the soil is!

import (
	"time"

	"tinygo.org/x/drivers"
)

const Address = 0x36

// BASE = 0x0F and OFFSET = 0x10
var readMoistureCommand = []uint8{0x0F, 0x10}

// BASE = 0x00 and OFFSET = 0x10
var readTemperatureCommand = []uint8{0x00, 0x04}

// Device wraps an I2C connection to a Soil device.
type Device struct {
	bus            drivers.I2C
	Address        uint16
	moistureBuf    [2]uint8
	temperatureBuf [4]uint8
}

// New creates a new SeeSaw Soil Sensor connection. The I2C bus must already be configured.
//
// This function only creates the Device object, it does not touch the device.
func New(bus drivers.I2C) *Device {
	return &Device{
		bus:     bus,
		Address: Address,
	}
}

// Read returns the moisture reading in range 0 to 1023
func (d *Device) ReadMoisture() (moisture uint16, err error) {

	for retry := 0; retry < 5; retry++ {
		err = d.bus.Tx(d.Address, readMoistureCommand, nil)
		if err != nil {
			continue
		}
		time.Sleep(time.Duration(3000+retry*1000) * time.Microsecond)
		err = d.bus.Tx(d.Address, nil, d.moistureBuf[:])
		if err == nil {
			return (uint16(d.moistureBuf[0]) << 8) | uint16(d.moistureBuf[1]), nil
		}
	}
	return // returns last error

}


// Read Temperature in degrees Celsius
func (d *Device) ReadTemperature() (temperature float64, err error) {
	for retry := 0; retry < 5; retry++ {
		err = d.bus.Tx(d.Address, readTemperatureCommand, nil)
		if err != nil {
			continue
		}
		time.Sleep(time.Duration(3000+retry*1000) * time.Microsecond)
		err = d.bus.Tx(d.Address, nil, d.temperatureBuf[:])
		if err == nil {
			tempRaw := ( uint32(d.temperatureBuf[0]) << 24) | uint32(d.temperatureBuf[1]) << 16 | uint32(d.temperatureBuf[2]) << 8 | uint32(d.temperatureBuf[3])
			// tempC := (float64(tempRaw) / 100_000)
			tempF := float64(tempRaw)/float64(100_000)*1.8 + float64(32)
			return tempF, nil
		}
	}
	return // returns last error
}