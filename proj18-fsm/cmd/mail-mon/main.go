package main

import (
	"aeg/pkg/fsm"
	"fmt"
	"log"
)

const (
	NothingDetected fsm.StateID = "NothingDetected"
	ReadingDetected fsm.StateID = "ReadingDetected"
	VehicleDetected fsm.StateID = "VehicleDetected"
	CarDetected     fsm.StateID = "CarDetected"
	MailDetected    fsm.StateID = "MailDetected"

	NSR    fsm.EventID = "NominalSensorReading"
	SSR    fsm.EventID = "SignificantSensorReading"
	IsCar  fsm.EventID = "IsCar"
	IsMail fsm.EventID = "IsMail"

	CarSensorReadingThreshold  = 100
	MailSensorReadingThreshold = 200
	CarReadingCountThreshold   = 2
	MailReadingCountThreshold  = 2
)

func main() {
	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("no impl here, run this instead:\ngo test -v ./...\n")

}

type MailMonContext struct {
	NothingCount     int
	ReadingCount     int
	CarReadingCount  int
	MailReadingCount int
	SensorReading    int
	CarDetectedCount int
	MailDetectedCount int

}

func (mmc *MailMonContext) String() string {
	mmcCopy := *mmc
	return fmt.Sprintf("MailMonContext: %v", mmcCopy)
}

// NothingDetectedAction represents the action executed on entering the NothingDetected state.
type NothingDetectedAction struct{}

func (a *NothingDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MailMonContext)
	ctx.NothingCount += 1

	log.Printf("nothing detected")
	return fsm.NoOp
}

// ReadingDetectedAction represents the action executed on entering the ReadingDetected state.
type ReadingDetectedAction struct{}

func (a *ReadingDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MailMonContext)
	ctx.NothingCount = 0 // reset
	ctx.ReadingCount += 1

	if ctx.SensorReading >= MailSensorReadingThreshold {
		ctx.MailReadingCount += 1
		log.Printf("increment possible mail count")
		return fsm.NoOp
	}
	if ctx.SensorReading >= CarSensorReadingThreshold {
		ctx.CarReadingCount += 1
		log.Printf("increment possible car count")
		return fsm.NoOp
	}

	log.Printf("we got an SSR but it was not a car or mail, that does not seem right")
	return fsm.NoOp

}

// VehicleDetectedAction represents the action executed on entering the VehicleDetected state.
type VehicleDetectedAction struct{}

func (a *VehicleDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*MailMonContext)
	ctx.NothingCount = 0 // reset
	ctx.ReadingCount = 0 // reset

	if ctx.MailReadingCount >= MailReadingCountThreshold {
		ctx.CarReadingCount = 0  // reset
		ctx.MailReadingCount = 0 //reset
		log.Printf("trigger IsMail event")
		return IsMail

	} else {
		ctx.CarReadingCount = 0  // reset
		ctx.MailReadingCount = 0 //reset
		log.Printf("trigger IsCar event")
		return IsCar
	}

}

// CarDetectedAction represents the action executed on entering the CarDetected state.
type CarDetectedAction struct{}

func (a *CarDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	ctx := eventCtx.(*MailMonContext)
	ctx.CarDetectedCount += 1

	log.Printf("a car has been detected")
	return fsm.NoOp
}

// MailDetectedAction represents the action executed on entering the MailDelivery state.
type MailDetectedAction struct{}

func (a *MailDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	ctx := eventCtx.(*MailMonContext)
	ctx.MailDetectedCount += 1
	
	log.Printf("a mail truck has been detected")
	return fsm.NoOp
}

func newMailMonFSM() *fsm.StateMachine {

	fsm := fsm.StateMachine{
		Current:  fsm.Default,
		Previous: fsm.Default,
		States: fsm.States{

			fsm.Default: fsm.State{
				Events: fsm.Events{NSR: NothingDetected},
			},

			NothingDetected: fsm.State{
				Action: &NothingDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: ReadingDetected,
				},
			},

			ReadingDetected: fsm.State{
				Action: &ReadingDetectedAction{},
				Events: fsm.Events{
					NSR: VehicleDetected,
					SSR: ReadingDetected,
				},
			},

			VehicleDetected: fsm.State{
				Action: &VehicleDetectedAction{},
				Events: fsm.Events{
					IsCar:  CarDetected,
					IsMail: MailDetected,
				},
			},

			CarDetected: fsm.State{
				Action: &CarDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: ReadingDetected,
				},
			},

			MailDetected: fsm.State{
				Action: &MailDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: ReadingDetected,
				},
			},
		},
	}

	// log.Printf("\n------------------------\n\n%+v\n\n----------------\n",fsm)

	return &fsm
}
