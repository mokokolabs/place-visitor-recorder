package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/us020"
)

const echoAPin = "P1_13"
const triggerAPin = "P1_15"
const echoBPin = "P1_16"
const triggerBPin = "P1_18"

// DistanceMeasure Measure distance with two sensors
type DistanceMeasure struct {
	readerA *us020.US020
	readerB *us020.US020
}

// InitDistanceMeasure Setup GPIO pins
func InitDistanceMeasure() DistanceMeasure {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}

	echoA, err := embd.NewDigitalPin(echoAPin)
	if err != nil {
		panic(err)
	}

	triggerA, err := embd.NewDigitalPin(triggerAPin)
	if err != nil {
		panic(err)
	}

	echoB, err := embd.NewDigitalPin(echoBPin)
	if err != nil {
		panic(err)
	}

	triggerB, err := embd.NewDigitalPin(triggerBPin)
	if err != nil {
		panic(err)
	}

	readerA := us020.New(echoA, triggerA, nil)
	readerB := us020.New(echoB, triggerB, nil)

	return DistanceMeasure{
		readerA: readerA,
		readerB: readerB,
	}
}

// ReadValues reads distance values
func (dm DistanceMeasure) ReadValues() (distanceA, distanceB float64) {
	a, err := dm.readerA.Distance()
	if err != nil {
		panic(err)
	}

	// b, err := dm.readerB.Distance()
	// if err != nil {
	// 	panic(err)
	// }

	return a, 0
}

// Cleanup everything
func (dm DistanceMeasure) Cleanup() {
	defer dm.readerA.Close()
	defer dm.readerB.Close()
	defer embd.CloseGPIO()
}
