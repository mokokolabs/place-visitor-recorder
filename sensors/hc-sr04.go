package sensors

import (
	"context"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

const (
	soundSpeed = 331.3 + 0.606*21
	pulseDelay = 30000 * time.Nanosecond
)

// HCSR04 is a ultrasonic distance measurer
type HCSR04 struct {
	echo    rpio.Pin
	trigger rpio.Pin
}

// NewHCSR04 creates a new HCSR04 sensor
func NewHCSR04(echoPin int, triggerPin int) *HCSR04 {
	echo := rpio.Pin(echoPin)
	echo.Input()

	trigger := rpio.Pin(triggerPin)
	trigger.Output()
	trigger.Low()

	return &HCSR04{echo, trigger}
}

// Measure returns distance in cm
func (sensor HCSR04) Measure() (value float64) {
	result := make(chan float64)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			result <- -1
		default:
			sensor.trigger.High()
			time.Sleep(pulseDelay)
			sensor.trigger.Low()

			for {
				if sensor.echo.Read() == rpio.High {
					break
				}
			}

			startTime := time.Now()

			for {
				if sensor.echo.Read() == rpio.Low {
					break
				}
			}

			duration := time.Since(startTime)
			result <- float64(duration.Nanoseconds()) / 10000000 * (soundSpeed / 2)
		}
	}()

	return <-result
}
