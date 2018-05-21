package main

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	echoOuterPin    = 21
	triggerOuterPin = 22
	echoInnerPin    = 0
	triggerInnerPin = 1
)

const (
	soundSpeed = 331.3 + 0.606*21
	pulseDelay = 30000 * time.Nanosecond
)

// DistanceMeasure Measure distance with two sensors
type DistanceMeasure struct {
	echoOuter    rpio.Pin
	triggerOuter rpio.Pin
	echoInner    rpio.Pin
	triggerInner rpio.Pin
}

// InitDistanceMeasure Setup GPIO pins
func InitDistanceMeasure() DistanceMeasure {
	if err := rpio.Open(); err != nil {
		panic(err)
	}

	echoOuter := rpio.Pin(echoOuterPin)
	echoOuter.Input()
	triggerOuter := rpio.Pin(triggerOuterPin)
	triggerOuter.Output()
	echoInner := rpio.Pin(echoInnerPin)
	echoInner.Input()
	triggerInner := rpio.Pin(triggerInnerPin)
	triggerInner.Output()

	return DistanceMeasure{
		echoOuter:    echoOuter,
		triggerOuter: triggerOuter,
		echoInner:    echoInner,
		triggerInner: triggerInner,
	}
}

func (dm DistanceMeasure) read(trigger rpio.Pin, echo rpio.Pin) (value float64) {
	trigger.High()
	time.Sleep(pulseDelay)
	trigger.Low()

	var start time.Time
	for s := echo.Read(); s == rpio.Low; {
		start = time.Now()
	}

	for s := echo.Read(); s == rpio.High; {
		return float64(time.Since(start).Nanoseconds()) / 10000000 * (soundSpeed / 2)
	}

	return 0
}

// ReadValues reads distance values
func (dm DistanceMeasure) ReadValues() (distanceOuter, distanceInner float64) {

	outer := dm.read(dm.triggerOuter, dm.echoOuter)
	inner := dm.read(dm.triggerInner, dm.echoInner)

	return outer, inner
}

// Cleanup pins
func (dm DistanceMeasure) Cleanup() {
	rpio.Close()
}
