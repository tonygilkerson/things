// Package lis2mdl implements a driver for the LIS2MDL,
// a magnetic sensor which is included on BBC micro:bit v1.5.
//
// Datasheet: https://www.st.com/resource/en/datasheet/lis2mdl.pdf
package lis2mdl // import "tinygo.org/x/drivers/lis2mdl"

import (
	"fmt"
	"math"
	"time"

	"tinygo.org/x/drivers"
)

// Device wraps an I2C connection to a LIS2MDL device.
type Device struct {
	bus        drivers.I2C
	Address    uint8
	PowerMode  uint8
	SystemMode uint8
	DataRate   uint8
}

// Configuration for LIS2MDL device.
type Configuration struct {
	PowerMode  uint8
	SystemMode uint8
	DataRate   uint8
}

// New creates a new LIS2MDL connection. The I2C bus must already be
// configured.
//
// This function only creates the Device object, it does not touch the device.
func New(bus drivers.I2C) Device {
	fmt.Printf("[DEBUG] ADDRESS: [%X]\n", ADDRESS)
	return Device{bus: bus, Address: ADDRESS}
}

// Connected returns whether LIS2MDL sensor has been found.
func (d *Device) Connected() bool {
	data := []byte{0}
	d.bus.ReadRegister(uint8(d.Address), WHO_AM_I, data)
	fmt.Printf("[DEBUG] WHO_AM_I: 0x%X, %d \n", data[0],data[0])
	// return data[0] == 0x40
	return data[0] == 0x4
}

// Configure sets up the LIS2MDL device for communication.
func (d *Device) Configure(cfg Configuration) {
	if cfg.PowerMode != 0 {
		d.PowerMode = cfg.PowerMode
	} else {
		d.PowerMode = POWER_NORMAL
	}

	if cfg.DataRate != 0 {
		d.DataRate = cfg.DataRate
	} else {
		d.DataRate = DATARATE_100HZ
	}

	if cfg.SystemMode != 0 {
		d.SystemMode = cfg.SystemMode
	} else {
		d.SystemMode = SYSTEM_CONTINUOUS
	}

	cmd := []byte{0}

	// reset
	cmd[0] = byte(1 << 5)
	err := d.bus.WriteRegister(uint8(d.Address), CFG_REG_A, cmd)
	if err != nil {
		fmt.Printf("[Configure] reset error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	// reboot
	cmd[0] = byte(1 << 6)
	err = d.bus.WriteRegister(uint8(d.Address), CFG_REG_A, cmd)
	if err != nil {
		fmt.Printf("[Configure] reboot error: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	// bdu
	cmd[0] = byte(1 << 4)
	err = d.bus.WriteRegister(uint8(d.Address), CFG_REG_C, cmd)
	if err != nil {
		fmt.Printf("[Configure] bdu error: %v", err)
	}

	// Temperature compensation is on for magnetic sensor (0x80)
	cmd[0] = byte(0x80)
	err = d.bus.WriteRegister(uint8(d.Address), CFG_REG_A, cmd)
	if err != nil {
		fmt.Printf("[Configure] reset error: %v", err)
	}

	// speed
	cmd[0] = byte(0x80 | d.DataRate)
	err = d.bus.WriteRegister(uint8(d.Address), CFG_REG_A, cmd)
	if err != nil {
		fmt.Printf("[Configure] speed error: %v", err)
	}

}

// ReadMagneticField reads the current magnetic field from the device and returns
// it in mG (milligauss). 1 mG = 0.1 µT (microtesla).
func (d *Device) ReadMagneticField() (x int32, y int32, z int32) {
	// turn back on read mode, even though it is supposed to be continuous?
	cmd := []byte{0}
	cmd[0] = byte(0x80 | d.PowerMode<<4 | d.DataRate<<2 | d.SystemMode)
	err := d.bus.WriteRegister(uint8(d.Address), CFG_REG_A, cmd)
	if err != nil {
		fmt.Printf("[ReadMagneticField] WriteRegister error: %v", err)
	}	
	time.Sleep(10 * time.Millisecond)

	data := make([]byte, 6)
	err = d.bus.ReadRegister(uint8(d.Address), OUTX_L_REG, data)
	if err != nil {
		fmt.Printf("[ReadMagneticField] ReadRegister error: %v", err)
	}		

	x = int32(int16((uint16(data[0]) << 8) | uint16(data[1])))
	y = int32(int16((uint16(data[2]) << 8) | uint16(data[3])))
	z = int32(int16((uint16(data[4]) << 8) | uint16(data[5])))

	return
}

// ReadCompass reads the current compass heading from the device and returns
// it in degrees. When the z axis is pointing straight to Earth and
// the y axis is pointing to North, the heading would be zero.
//
// However, the heading may be off due to electronic compasses would be effected
// by strong magnetic fields and require constant calibration.
func (d *Device) ReadCompass() (h int32) {
	x, y, _ := d.ReadMagneticField()
	xf, yf := float64(x)*0.15, float64(y)*0.15

	rh := (math.Atan2(yf, xf) * 180) / math.Pi
	if rh < 0 {
		rh = 360 + rh
	}

	return int32(rh)
}
