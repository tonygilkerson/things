package main

// This example code demonstrates Lora RX/TX With SX127x driver
// You need to connect SPI, RST, CS, DIO0 (aka IRQ) and DIO1 to use.

import (
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx127x"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

var (
	loraRadio *sx127x.Device

	SX127X_PIN_EN   = machine.GP15
	SX127X_PIN_RST  = machine.GP20
	SX127X_PIN_CS   = machine.GP17
	SX127X_PIN_DIO0 = machine.GP21 // (GP21--G0) Must be connected from pico to breakout for radio events IRQ to work
	SX127X_PIN_DIO1 = machine.GP22 // (GP22--G1)I don't now what this does
	SX127X_SPI      = machine.SPI0
)

func dioIrqHandler(machine.Pin) {
	loraRadio.HandleInterrupt()
}

func main() {
	time.Sleep(1 * time.Second)
	led := machine.LED //GP25
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	runLight(led,15)


	machine.SPI0.Configure(machine.SPIConfig{
		SCK: machine.SPI0_SCK_PIN, // GP18
		SDO: machine.SPI0_SDO_PIN, // GP19
		SDI: machine.SPI0_SDI_PIN, // GP16
	})

	println("\n# TinyGo Lora RX/TX test")
	println("# ----------------------")
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	SX127X_PIN_EN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SX127X_PIN_EN.High() //Enabled by default
	SX127X_PIN_RST.Configure(machine.PinConfig{Mode: machine.PinOutput})

	SX127X_SPI.Configure(machine.SPIConfig{Frequency: 500000, Mode: 0})

	println("main: create and start SX127x driver")
	loraRadio = sx127x.New(*SX127X_SPI, SX127X_PIN_RST)
	loraRadio.SetRadioController(sx127x.NewRadioControl(SX127X_PIN_CS, SX127X_PIN_DIO0, SX127X_PIN_DIO1))

	loraRadio.Reset()
	state := loraRadio.DetectDevice()
	if !state {
		panic("main: sx127x NOT FOUND !!!")
	} else {
		println("main: sx127x found")
	}

	// Prepare for Lora Operation
	loraConf := lora.Config{
		Freq:           lora.MHz_916_8,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPrivate,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)

	var count uint
	for {
		tStart := time.Now()

		//
		// RX
		//
		println("main: Receiving Lora for 10 seconds")
		for time.Since(tStart) < 10*time.Second {
			buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
			if err != nil {
				println("RX Error: ", err)
			} else if buf != nil {
				println("Packet Received: len=", len(buf), string(buf))
			}
		}
		println("main: End Lora RX")

		//
		// TX
		//
		strMsg := fmt.Sprintf("PICO-WITH-LORA-BREAKOUT-%v", count)
		msg := []byte(strMsg)

		println("LORA TX size=", len(msg), " -> ", string(msg))

		err := loraRadio.Tx(msg, LORA_DEFAULT_TXTIMEOUT_MS)

		if err != nil {
			println("TX Error:", err)
		}
		count++
	}

}

///////////////////////////////////////////////////////////////////////////////
//		functions
///////////////////////////////////////////////////////////////////////////////
func runLight(led machine.Pin, count int) {

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.High()
		time.Sleep(time.Millisecond * 100)
		led.Low()
		time.Sleep(time.Millisecond * 100)
		print(".")
	}

}