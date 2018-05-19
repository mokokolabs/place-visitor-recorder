package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/looplab/fsm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var distanceThresholdStr = os.Getenv("DISTANCE_THRESHOLD")
var distanceThreshold, err = strconv.ParseFloat(distanceThresholdStr, 64)

var apiRoot = os.Getenv("API_ROOT")
var apiKey = os.Getenv("API_KEY")
var collectorURL = fmt.Sprintf("%v/api/v1/place_visit_events", apiRoot)

// Send event to API
func sendEvent(direction string) {
	payload := map[string]interface{}{
		"data": map[string]string{
			"key":         apiKey,
			"direction":   direction,
			"happened_at": time.Now().Format(time.RFC3339),
		},
	}

	bytePayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(collectorURL, "application/json", bytes.NewBuffer(bytePayload))

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(err)
	}
}

func main() {
	dm := InitDistanceMeasure()
	sm := fsm.NewFSM("both-idle", fsm.Events{
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
	}, fsm.Callbacks{
		"enter_entered": func(e *fsm.Event) {
			fmt.Println("Sending enter event")
		},
		"enter_exited": func(e *fsm.Event) {
			fmt.Println("Sending exit event")
		},
	})

	fmt.Println(sm.Current())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	oldOuter, oldInner := distanceThreshold, distanceThreshold

	for {
		select {
		default:
			fmt.Print("\033[H\033[2J")
			distanceOuter, distanceInner := dm.ReadValues()
			fmt.Printf("Distance inner is %.2f\n", distanceOuter)
			fmt.Printf("Distance outer is %.2f\n", distanceInner)

			if distanceOuter <= oldOuter {
				sm.Event("outer-activated")
			} else {
				sm.Event("outer-deactivated")
			}
			oldOuter = distanceOuter

			if distanceInner <= oldInner {
				sm.Event("inner-activated")
			} else {
				sm.Event("inner-deactivated")
			}
			oldInner = distanceInner

			time.Sleep(100 * time.Millisecond)
		case <-quit:
			return
		}
	}
}
