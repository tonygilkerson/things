package main

// In this example, a Lora packet will be sent every 10s
// module will be in RX mode between two transmissions

import (
	"machine"
	"time"

	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/sx126x"
)

const (
	LORA_DEFAULT_RXTIMEOUT_MS = 1000
	LORA_DEFAULT_TXTIMEOUT_MS = 5000
)

var (
	loraRadio *sx126x.Device
	txmsg     = []byte("Hi from Gateway")
)

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// run light
  runLight(10)
	time.Sleep(time.Second * 10)

	println("High")
	machine.LED.High()
	time.Sleep(time.Second * 10)
	
	println("Low")
	machine.LED.Low()
	
	time.Sleep(time.Second * 10)
	println("Continue")

	//
	// Input Pin
	//
	// pinPB10 := machine.PB10
	// pinPB10.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	// pinPA9 := machine.PA9
	// pinPA9.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	// pinPA0 := machine.PA0
	// pinPA0.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})


	// Create the driver
	loraRadio = sx126x.New(spi)
	loraRadio.SetDeviceType(sx126x.DEVICE_TYPE_SX1262)

	// Create radio controller for target
	loraRadio.SetRadioController(newRadioControl())

	// Detect the device
	state := loraRadio.DetectDevice()
	if !state {
		panic("sx126x not detected.")
	}

	loraConf := lora.Config{
		Freq:           lora.MHz_916_8,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Ldr:            lora.LowDataRateOptimizeOff,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPrivate,
		LoraTxPowerDBm: 20,
	}

	loraRadio.LoraConfig(loraConf)

	var count uint
	for {
		start := time.Now()

		// println("pinPB10 ", pinPB10.Get())
		// println("pinPA9 ", pinPA9.Get())
		// println("pinPA0 ", pinPA0.Get())

		println("Receiving for 5 seconds")
		for time.Since(start) < 5*time.Second {
			buf, err := loraRadio.Rx(LORA_DEFAULT_RXTIMEOUT_MS)
			if err != nil {
				println("RX Error: ", err)
			} else if buf != nil {
				println("Packet Received: len=", len(buf), string(buf))
				// runLight(5)
			}
		}

		println("Send TX size=", len(txmsg), " -> ", string(txmsg))
		err := loraRadio.Tx(txmsg, LORA_DEFAULT_TXTIMEOUT_MS)
		if err != nil {
			println("TX Error:", err)
		}
		count++

		runLight(3)
	}

}

func runLight(count int) {

	// run light
	led := machine.LED

	// blink run light for a bit seconds so I can tell it is starting
	for i := 0; i < count; i++ {
		led.Low()
		time.Sleep(time.Millisecond * 200)
		// Do high last because we want it to be off and for some reason
		// high is off on lore E5 board, strange
		led.High()
		time.Sleep(time.Millisecond * 200)
	}
}
