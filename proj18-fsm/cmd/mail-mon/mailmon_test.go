package main

import (
	"testing"
)

func TestMailMonStateMachine(t *testing.T) {
	var err error

	// Create a new instance of the light switch state machine.
	mailMonFSM := newMailMonFSM()
	mailMonContext := &MailMonContext{
		NothingCount:      0,
		ReadingCount:      0,
		CarReadingCount:   0,
		MailReadingCount:  0,
		SensorReading:     0,
		CarDetectedCount:  0,
		MailDetectedCount: 0,
	}

	mailMonContext.SensorReading = 0
	err = mailMonFSM.SendEvent(NSR, mailMonContext)
	if err != nil {
		t.Errorf("%v", err)
	}

	mailMonContext.SensorReading = 0
	err = mailMonFSM.SendEvent(NSR, mailMonContext)

	//
	// A car
	//
	mailMonContext.CarDetectedCount = 0
	mailMonContext.MailDetectedCount = 0
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)

	mailMonContext.SensorReading = 0
	err = mailMonFSM.SendEvent(NSR, mailMonContext)
	if mailMonContext.CarDetectedCount != 1 && mailMonContext.MailDetectedCount != 0 {
		t.Errorf("Expected CarDetectedCount=1 and MailDetectedCount=0 got: %v, %v", mailMonContext.CarDetectedCount, mailMonContext.MailDetectedCount)
	}

	//
	// A mail
	//
	mailMonContext.CarDetectedCount = 0
	mailMonContext.MailDetectedCount = 0
	mailMonContext.SensorReading = MailSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = MailSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = MailSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)

	mailMonContext.SensorReading = 0
	err = mailMonFSM.SendEvent(NSR, mailMonContext)
	if mailMonContext.CarDetectedCount != 0 && mailMonContext.MailDetectedCount != 1 {
		t.Errorf("Expected CarDetectedCount=0 and MailDetectedCount=1 got: %v, %v", mailMonContext.CarDetectedCount, mailMonContext.MailDetectedCount)
	}

	//
	// A slow approaching (or leaving) mail truck
	// some of the mail truck reading will look like a car
	//
	mailMonContext.CarDetectedCount = 0
	mailMonContext.MailDetectedCount = 0

	// car
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)

	// mail
	mailMonContext.SensorReading = MailSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = MailSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = MailSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)

	// car
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)
	mailMonContext.SensorReading = CarSensorReadingThreshold + 1
	err = mailMonFSM.SendEvent(SSR, mailMonContext)

	mailMonContext.SensorReading = 0
	err = mailMonFSM.SendEvent(NSR, mailMonContext)
	if mailMonContext.CarDetectedCount != 0 && mailMonContext.MailDetectedCount != 1 {
		t.Errorf("Expected CarDetectedCount=0 and MailDetectedCount=1 got: %v, %v", mailMonContext.CarDetectedCount, mailMonContext.MailDetectedCount)
	}
}
