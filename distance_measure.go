package main

import (
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/sensor/us020"
)

const echoOuterPin = "P1_13"
const triggerOuterPin = "P1_15"
const echoInnerPin = "P1_16"
const triggerInnerPin = "P1_18"

// DistanceMeasure Measure distance with two sensors
type DistanceMeasure struct {
	readerOuter *us020.US020
	readerInner *us020.US020
}

// InitDistanceMeasure Setup GPIO pins
func InitDistanceMeasure() DistanceMeasure {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}

	echoOuter, err := embd.NewDigitalPin(echoOuterPin)
	if err != nil {
		panic(err)
	}

	triggerOuter, err := embd.NewDigitalPin(triggerOuterPin)
	if err != nil {
		panic(err)
	}

	echoInner, err := embd.NewDigitalPin(echoInnerPin)
	if err != nil {
		panic(err)
	}

	triggerInner, err := embd.NewDigitalPin(triggerInnerPin)
	if err != nil {
		panic(err)
	}

	readerOuter := us020.New(echoOuter, triggerOuter, nil)
	readerInner := us020.New(echoInner, triggerInner, nil)

	return DistanceMeasure{
		readerOuter: readerOuter,
		readerInner: readerInner,
	}
}

// ReadValues reads distance values
func (dm DistanceMeasure) ReadValues() (distanceOuter, distanceInner float64) {
	outer, err := dm.readerOuter.Distance()
	if err != nil {
		panic(err)
	}

	inner, err := dm.readerInner.Distance()
	if err != nil {
		panic(err)
	}

	return outer, inner
}

// Cleanup everything
func (dm DistanceMeasure) Cleanup() {
	defer dm.readerOuter.Close()
	defer dm.readerInner.Close()
	defer embd.CloseGPIO()
}
