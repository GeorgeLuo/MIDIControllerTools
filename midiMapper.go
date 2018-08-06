package main

import (
	"fmt"
	"log"
	"github.com/rakyll/portmidi"
)

/*
	Design:
	EnvironmentSession: map of devices

	SDK functions:
	ListenAll: output array of event -> instrument mappings
	map device toggles to any device toggle
*/

var deviceLayout DeviceLayout

type DeviceLayout struct {
	devices []Device
}

type Device struct {
	name string
	input bool
	output bool
	open bool
}

func initializeDeviceLayout() {
	portmidi.Initialize()
	var numDevices = portmidi.CountDevices()
	deviceList := make([]Device, numDevices)
	for i := 0; i < numDevices; i++ {
		var portInfo = portmidi.Info(portmidi.DeviceID(i))
		fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name, 
			portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)

		deviceList[i] = Device{name: portInfo.Name, input: portInfo.IsInputAvailable, output: portInfo.IsOutputAvailable, open: portInfo.IsOpened}
	}
	deviceLayout = DeviceLayout{devices: deviceList}

}

type KnobToKnobMapping struct {
	FromDevice int64
	ToDevice int64
	FromKnob int64
	ToKnob int64
}

func pollDevice(stream portmidi.Stream) {
	
	result, err := stream.Poll()
	if err != nil {
		log.Fatal(err)
	}

	if result {
		fmt.Println("Something happened!")
	} else {
		fmt.Println("Waiting ...")
	}
}

func readFromStreamForSamplesNEveryTimeT(stream portmidi.Stream, numSamples int, duration int, maxEvents int) {
	for i := 0; i < numSamples; i++ {
		events, err := stream.Read(maxEvents)	
		if err != nil {
			log.Fatal(err)
		} else {
			for j := 0; j < len(events); j++ {
				fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)
			}
		}
	}
}

func readFromIncomingDevices(streams []portmidi.Stream) {
	for i := 0; i < len(streams); i++ {

	}
}

func readFromDeviceWriteToDevice(OutStream portmidi.Stream, InStream portmidi.Stream, maxEvents int) {
	for {
		events, err := OutStream.Read(maxEvents)
		if err != nil {
			log.Fatal(err)
		} else {
			if len(events) > 0 {
				InStream.Write(events)
				for j := 0; j < len(events); j++ {
					fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)
				}
			}
		}
	}
}

func main(){
	fmt.Printf("Reading from midi channels\n")
	
	portmidi.Initialize()

	var numDevices = portmidi.CountDevices() // returns the number of MIDI devices
	fmt.Printf("Num devices: %d\n\n", numDevices)
	// portmidi.Info(deviceID) // returns info about a MIDI device

	for i := 0; i < numDevices; i++ {
		var portInfo = portmidi.Info(portmidi.DeviceID(i))
		fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name, 
			portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)
	}

	var defaultInputDeviceID = portmidi.DefaultInputDeviceID() // returns the ID of the system default input
	fmt.Printf("Default input device ID: %d\n", defaultInputDeviceID)

	var defaultOutDeviceID = portmidi.DefaultOutputDeviceID() // returns the ID of the system default output
	fmt.Printf("Default output device ID: %d\n\n", defaultOutDeviceID)

	InStream, err := portmidi.NewInputStream(portmidi.DeviceID(1), 1024)
	OutStream, err := portmidi.NewOutputStream(portmidi.DeviceID(6), 1024, 0)

	if err != nil {
		log.Fatal(err)
	}

	readFromDeviceWriteToDevice(*InStream, *OutStream, 10)

}