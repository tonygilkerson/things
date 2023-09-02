package marty

// To run tests
// $ go test -v ./...
//

import (
	"testing"
)

func TestMartyStateMachine(t *testing.T) {

	var m *Marty
	

	//
	// A car arriving
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(FarRising,m.Ctx)
	m.StateMachine.SendEvent(NearRising,m.Ctx)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 1 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("A car arriving\nexpected: {DefaultCount:1 ArrivedCount:1 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// A car departing
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(NearRising,m.Ctx)
	m.StateMachine.SendEvent(FarRising,m.Ctx)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 1 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("A car departing\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:0 DepartedCount:1 DepartingCount:1 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// FalseAlarm from the Arriving direction
	// A car approaching but stops short, turns around, backups up or something
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(FarRising,m.Ctx)
	m.StateMachine.SendEvent(FarFalling,m.Ctx)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 1 {
		// all good
	} else {
		t.Errorf("FalseAlarm from the Arriving direction\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:1}\ngot:      %+v", m.Ctx)
	}

	//
	// FalseAlarm from the Departing direction
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(NearRising,m.Ctx)
	m.StateMachine.SendEvent(NearFalling,m.Ctx)

	if m.Ctx.DefaultCount == 1 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 1 {
		// all good
	} else {
		t.Errorf("FalseAlarm from the Departing direction\nexpected: {DefaultCount:1 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:0 FalseAlarmCount:1}\ngot:      %+v", m.Ctx)
	}

	//
	// Error from the Departing direction
	// Error - should never get two Rising events in a row from the same direction
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(NearRising,m.Ctx)
	m.StateMachine.SendEvent(NearRising,m.Ctx)

	if m.Ctx.DefaultCount == 0 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 1 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error from the Departing direction\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:0 DepartedCount:0 DepartingCount:1 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// Error from the Arriving direction
	// Error - should never get two Rising events in a row from the same direction
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(FarRising,m.Ctx)
	m.StateMachine.SendEvent(FarRising,m.Ctx)

	if m.Ctx.DefaultCount == 0 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 1 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error from the Arriving direction\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// Default goes to Default if LD or RD
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(NearRising,m.Ctx)
	m.StateMachine.SendEvent(FarRising,m.Ctx)
	m.StateMachine.SendEvent(FarFalling,m.Ctx)
	m.StateMachine.SendEvent(NearFalling,m.Ctx)

	if m.Ctx.DefaultCount == 3 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 1 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 0 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Error default goes to default\nexpected: {DefaultCount:3 ArrivedCount:0 ArrivingCount:0 DepartedCount:1 DepartingCount:1 ErrorCount:0 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}

	//
	// I have see this but I am not sure how it happens.  
	// I think the PIRs are timing out at different rates 
	// or I have a hardware issue, or maybe I am running the PIRs with the wrong voltage
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(NearRising,m.Ctx)
	m.StateMachine.SendEvent(FarFalling,m.Ctx)


	if m.Ctx.DefaultCount == 0 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 0 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 1 &&
		m.Ctx.ErrorCount == 1 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Falling out of order 1\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:0 DepartedCount:0 DepartingCount:1 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}


	//
	// I have see this but I am not sure how it happens.  I think the PIRs are timing out at different rates
	//
	t.Logf("----------------------------------\n")
	m = New()
	m.ResetContext()
	m.StateMachine.SendEvent(FarRising,m.Ctx)
	m.StateMachine.SendEvent(NearFalling,m.Ctx)


	if m.Ctx.DefaultCount == 0 &&
		m.Ctx.ArrivedCount == 0 &&
		m.Ctx.ArrivingCount == 1 &&
		m.Ctx.DepartedCount == 0 &&
		m.Ctx.DepartingCount == 0 &&
		m.Ctx.ErrorCount == 1 &&
		m.Ctx.FalseAlarmCount == 0 {
		// all good
	} else {
		t.Errorf("Falling out of order 2\nexpected: {DefaultCount:0 ArrivedCount:0 ArrivingCount:1 DepartedCount:0 DepartingCount:0 ErrorCount:1 FalseAlarmCount:0}\ngot:      %+v", m.Ctx)
	}


}
