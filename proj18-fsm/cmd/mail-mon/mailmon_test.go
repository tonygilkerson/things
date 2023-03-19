package main

import (
	"testing"
	// "aeg/pkg/fsm"
)


func TestMailMonStateMachine(t *testing.T) {
	var err error

	// Create a new instance of the light switch state machine.
	mailMonFSM := newMailMonFSM()
	mailMonContext := &MailMonContext {
		NothingCount: 0,
		CandidateCount: 0,
		CarCount: 0,
		MailCount: 0,
	}


	// Set the initial "NSR" state in the state machine.
	err = mailMonFSM.SendEvent(NSR, mailMonContext)
	if err != nil {
		t.Errorf("Couldn't set the initial state of the state machine, err: %v", err)
	}

	// Sen NSR again to see the count increase
	err = mailMonFSM.SendEvent(NSR, mailMonContext)
	if err != nil {
		t.Errorf("Couldn't set NSR state for the second time, err: %v", err)
	}



}