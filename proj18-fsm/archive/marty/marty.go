package marty

import (
	"fmt"
	"log"

	"github.com/tonygilkerson/marty/pkg/fsm"
)

const (
	// States
	Arriving   fsm.StateID = "Arriving"
	Arrived    fsm.StateID = "Arrived"
	Departing  fsm.StateID = "Departing"
	Departed   fsm.StateID = "Departed"
	FalseAlarm fsm.StateID = "FalseAlarm"
	Error      fsm.StateID = "Error"

	//Events
	FarRising   fsm.EventID = "FarRising"
	FarFalling  fsm.EventID = "FarFalling"
	NearRising  fsm.EventID = "NearRising"
	NearFalling fsm.EventID = "NearFalling"
	Reset       fsm.EventID = "Reset"
)

type Context struct {
	DefaultCount    int
	ArrivedCount    int
	ArrivingCount   int
	DepartedCount   int
	DepartingCount  int
	ErrorCount      int
	FalseAlarmCount int
}

func (c *Context) String() string {
	cCopy := *c
	return fmt.Sprintf("Context: %+v\n", cCopy)
}

type Marty struct {
	StateMachine fsm.StateMachine
	Ctx          Context
}


func (m *Marty) ResetContext() {
	m.Ctx = Context{
		DefaultCount:    0,
		ArrivedCount:    0,
		ArrivingCount:   0,
		DepartedCount:   0,
		DepartingCount:  0,
		ErrorCount:      0,
		FalseAlarmCount: 0,
	}
}

// // MarshallMetrics will format the Context into a message that can be sent
// func (m *Marty) MarshallMetrics() string {
// 	msg := fmt.Sprintf("mbx|%d|%d|%d|%d",
// 		m.Ctx.ArrivedCount,
// 		m.Ctx.DepartedCount,
// 		m.Ctx.ErrorCount,
// 		m.Ctx.FalseAlarmCount,
// 	)

// 	return msg
// }

// // UnmarshallMetrics will unmarshall a message that was produced by MarshallMetrics
// func UnmarshallMetrics(msg string) Context {

// 	var ctx Context

// 	msgParts := strings.Split(msg, "|")

// 	if msgParts[0] != "mbx" {
// 		log.Printf("expected metrics message got: %v\n", msg)
// 		return ctx
// 	}

// 	ctx.ArrivedCount, _ = strconv.Atoi(msgParts[1])
// 	ctx.DepartedCount, _ = strconv.Atoi(msgParts[2])
// 	ctx.ErrorCount, _ = strconv.Atoi(msgParts[3])
// 	ctx.FalseAlarmCount, _ = strconv.Atoi(msgParts[4])

// 	return ctx
// }

// DefaultAction
type DefaultAction struct{}

func (a *DefaultAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.DefaultCount += 1

	log.Printf("DefaultAction\n")
	return fsm.NoOp
}

// ArrivedAction
type ArrivedAction struct{}

func (a *ArrivedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.ArrivedCount += 1

	log.Printf("ArrivedAction\n")
	return Reset
}

type DepartedAction struct{}

func (a *DepartedAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.DepartedCount += 1

	log.Printf("DepartedAction\n")
	return Reset
}

// ArrivingAction
type ArrivingAction struct{}

func (a *ArrivingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.ArrivingCount += 1

	log.Printf("ArrivingAction\n")
	return fsm.NoOp
}

// DepartingAction
type DepartingAction struct{}

func (a *DepartingAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.DepartingCount += 1

	log.Printf("DepartingAction\n")
	return fsm.NoOp
}

// DepartedAction

// ErrorAction
type ErrorAction struct{}

func (a *ErrorAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.ErrorCount += 1

	log.Printf("ErrorAction\n")
	return fsm.NoOp
}

// FalseAlarmAction
type FalseAlarmAction struct{}

func (a *FalseAlarmAction) Execute(eventCtx fsm.EventContext) fsm.EventID {

	ctx := eventCtx.(*Context)
	ctx.FalseAlarmCount += 1

	log.Printf("FalseAlarmAction\n")
	return Reset
}

func New() *Marty {

	var marty Marty
	marty.StateMachine = fsm.StateMachine{
		Current:  fsm.Default,
		Previous: fsm.Default,
		States: fsm.States{

			fsm.Default: fsm.State{
				Action: &DefaultAction{},
				Events: fsm.Events{
					FarRising:   Arriving,
					NearRising:  Departing,
					FarFalling:  fsm.Default,
					NearFalling: fsm.Default,
				},
			},

			Arriving: fsm.State{
				Action: &ArrivingAction{},
				Events: fsm.Events{
					FarFalling: FalseAlarm,
					NearRising: Arrived,
				},
			},

			Arrived: fsm.State{
				Action: &ArrivedAction{},
				Events: fsm.Events{
					Reset: fsm.Default,
				},
			},

			Departing: fsm.State{
				Action: &DepartingAction{},
				Events: fsm.Events{
					NearFalling: FalseAlarm,
					FarRising:   Departed,
				},
			},

			Departed: fsm.State{
				Action: &DepartedAction{},
				Events: fsm.Events{
					Reset: fsm.Default,
				},
			},

			FalseAlarm: fsm.State{
				Action: &FalseAlarmAction{},
				Events: fsm.Events{
					Reset: fsm.Default,
				},
			},
		},
	}

	return &marty
}
