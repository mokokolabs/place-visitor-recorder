package main

import (
	"github.com/looplab/fsm"
	"github.com/stianeikeland/go-rpio"
	"os"
	"os/signal"
	"strconv"
	"time"

	"place-visitor-recorder/api"
	"place-visitor-recorder/sensors"
)

var stateEvents = fsm.Events{
	// Visitor activated outer first
	{Name: "outer-activated", Src: []string{"both-idle"}, Dst: "outer-active"},
	// Visitor activated inner first
	{Name: "inner-activated", Src: []string{"both-idle"}, Dst: "inner-active"},
	// Visitor activated also inner
	{Name: "inner-activated", Src: []string{"outer-active"}, Dst: "entered"},
	// Visitor activated also outer
	{Name: "outer-activated", Src: []string{"inner-active"}, Dst: "exited"},
	// Outer got deactivated on enter, inner still active
	{Name: "outer-deactivated", Src: []string{"entered"}, Dst: "inner-active"},
	// Inner got deactivated on exit, outer still active
	{Name: "inner-deactivated", Src: []string{"exited"}, Dst: "outer-active"},
	// Single one deactivates
	{Name: "outer-deactivated", Src: []string{"outer-active"}, Dst: "both-idle"},
	{Name: "inner-deactivated", Src: []string{"inner-active"}, Dst: "both-idle"},
}

var distanceThresholdStr = os.Getenv("DISTANCE_THRESHOLD")
var distanceThreshold, err = strconv.ParseFloat(distanceThresholdStr, 64)

func main() {
	sm := fsm.NewFSM("both-idle", stateEvents, fsm.Callbacks{
		"enter_entered": func(e *fsm.Event) {
			api.SendEvent("enter")
		},
		"enter_exited": func(e *fsm.Event) {
			api.SendEvent("exit")
		},
	})

	if err := rpio.Open(); err != nil {
		panic(err)
	}
	defer rpio.Close()

	outerSensor := sensors.NewHCSR04(21, 22)
	innerSensor := sensors.NewHCSR04(0, 1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	for {
		select {
		default:
			distanceOuter := outerSensor.Measure()
			distanceInner := innerSensor.Measure()

			if distanceOuter <= distanceThreshold {
				sm.Event("outer-activated")
			} else {
				sm.Event("outer-deactivated")
			}

			if distanceInner <= distanceThreshold {
				sm.Event("inner-activated")
			} else {
				sm.Event("inner-deactivated")
			}

			time.Sleep(50 * time.Millisecond)
		case <-quit:
			return
		}
	}
}
