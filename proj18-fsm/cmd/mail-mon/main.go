package main

import (
	"aeg/pkg/fsm"
	"fmt"
	"log"
)

const (
	NothingDetected fsm.StateID = "NothingDetected"
	CandidateDetected fsm.StateID = "CandidateDetected"
	CarDetected fsm.StateID = "CarDetected"
	MailDetected fsm.StateID = "MailDetected"

	
	NSR fsm.EventID = "NominalSensorReading"
	SSR fsm.EventID = "SignificantSensorReading"
	IsCar fsm.EventID = "IsCar"
	IsMail fsm.EventID = "IsMail"
	
)



func main() {
	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	log.Printf("no impl here, run this instead:\ngo test -v ./...\n")

}

type MailMonContext struct {
	NothingCount int
	CandidateCount int
	CarCount int
	MailCount int

}
func (mmc *MailMonContext) String() string {
	return fmt.Sprintf("MailMonContext: NothingCount: %v",mmc.NothingCount)
}

// NothingDetectedAction represents the action executed on entering the NothingDetected state.
type NothingDetectedAction struct{}

func (a *NothingDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
  
	ctx := eventCtx.(*MailMonContext)
	
	ctx.NothingCount += 1
	
	log.Printf("nothing detected, NothingCount: %v",ctx.NothingCount)
	return fsm.NoOp
}

// CandidateDetectedAction represents the action executed on entering the CandidateDetected state.
type CandidateDetectedAction struct{}

func (a *CandidateDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	log.Println("a candidate has been detected")
	return fsm.NoOp
}

// CarDetectedAction represents the action executed on entering the CarDetected state.
type CarDetectedAction struct{}

func (a *CarDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	log.Println("a car has been detected")
	return fsm.NoOp
}

// MailDetectedAction represents the action executed on entering the MailDelivery state.
type MailDetectedAction struct{}

func (a *MailDetectedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	log.Println("a mail delivery has been detected")
	return fsm.NoOp
}

func newMailMonFSM() *fsm.StateMachine {
	
	fsm := fsm.StateMachine{
		Current: fsm.Default,
		Previous: fsm.Default,
		States: fsm.States{

			fsm.Default: fsm.State{
				Events: fsm.Events{ NSR: NothingDetected,},
			},

			NothingDetected: fsm.State{
				Action: &NothingDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: CandidateDetected,
				},
			},

			CandidateDetected: fsm.State{
				Action: &CandidateDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: CandidateDetected,
					IsCar: CarDetected,
					IsMail: MailDetected,
				},
			},

			CarDetected: fsm.State{
				Action: &CarDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: CandidateDetected,
				},
			},

			MailDetected: fsm.State{
				Action: &MailDetectedAction{},
				Events: fsm.Events{
					NSR: NothingDetected,
					SSR: CandidateDetected,
				},
			},
		},
	}

	log.Printf("\n------------------------\n\n%+v\n\n----------------\n",fsm)

	return &fsm
}

