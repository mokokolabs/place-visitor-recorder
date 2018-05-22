package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var apiRoot = os.Getenv("API_ROOT")
var apiKey = os.Getenv("API_KEY")
var collectorURL = fmt.Sprintf("%v/api/v1/place_visit_events", apiRoot)

// SendEvent send event to API
func SendEvent(direction string) {
	happenedAt := time.Now().Format(time.RFC3339)

	payload := map[string]interface{}{
		"data": map[string]string{
			"key":         apiKey,
			"direction":   direction,
			"happened_at": happenedAt,
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
	} else {
		fmt.Printf("Sent %v event at %v\n", direction, happenedAt)
	}
}
