package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	for {
		select {
		default:
			distanceA, distanceB := dm.ReadValues()
			fmt.Print("\033[H\033[2J")
			fmt.Printf("Distance A is %.2f\n", distanceA)
			fmt.Printf("Distance B is %.2f\n", distanceB)

			time.Sleep(100 * time.Millisecond)
		case <-quit:
			return
		}
	}
}
