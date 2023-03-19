package main

import (
	"aeg/pkg/fsm"
	"log"
)

const (
	Off fsm.StateID = "Off"
	On  fsm.StateID = "On"
	
	SwitchOff fsm.EventID = "SwitchOff"
	SwitchOn  fsm.EventID = "SwitchOn"
)


func main() {
	// Log to the console with date, time and filename prepended
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	log.Printf("No impl here, run this instead:\ngo test -v ./...\n")

}

// OffAction represents the action executed on entering the Off state.
type OffAction struct{}

func (a *OffAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	log.Println("The light has been switched off")
	return fsm.NoOp
}

// OnAction represents the action executed on entering the On state.
type OnAction struct{}

func (a *OnAction) Execute(eventCtx fsm.EventContext) fsm.EventID {
	log.Println("The light has been switched on")
	return fsm.NoOp
}

func newLightSwitchFSM() *fsm.StateMachine {
	
	lightSwitchFSM := fsm.StateMachine{
		Current: fsm.Default,
		Previous: fsm.Default,
		States: fsm.States{
			fsm.Default: fsm.State{
				Events: fsm.Events{ SwitchOff: Off,},
			},
			Off: fsm.State{
				Action: &OffAction{},
				Events: fsm.Events{SwitchOn: On,},
			},
			On: fsm.State{
				Action: &OnAction{},
				Events: fsm.Events{SwitchOff: Off,},
			},
		},
	}

	log.Printf("\n------------------------\n\n%+v\n\n----------------\n",lightSwitchFSM)

	return &lightSwitchFSM
}

//
// newLightSwitchFSM2 does the same as newLightSwitchFSM but using a different method
// I am just trying to figure out which I like best
//
func newLightSwitchFSM2()  *fsm.StateMachine {
	
	var stateMachine fsm.StateMachine
	var states fsm.States
	var state fsm.State
	var events fsm.Events

	states = make(fsm.States)
	
	stateMachine.Current = fsm.Default
	stateMachine.Previous = fsm.Default
	
	//
	// Default state
	//
	var defaultState fsm.State
	defaultState.Action = nil
	events = make(fsm.Events)
	events[SwitchOff] = Off
	defaultState.Events = events
	states[fsm.Default] = defaultState
	
	//
	// Off state
	//
	var offState fsm.State
	offState.Action = &OffAction{}
	events = make(fsm.Events)
	events[SwitchOn] = On
	offState.Events = events
	states[Off] = offState
	
	
	//
	// On state
	//
	var onState fsm.State
	onState.Action = &OnAction{}
	events = make(fsm.Events)
	state.Events = events
	events[SwitchOff] = Off
	onState.Events = events
	states[On] = onState


	stateMachine.States = states

	log.Printf("\n------------------------\n\n%+v\n\n----------------\n",stateMachine)
	return &stateMachine
}
