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

const (

/*
	declare constants for channel 1-3
	10010000= 90= 144| Chan 1      Note on        |      "        |      "
	10010001= 91= 145| Chan 2         "           |      "        |      "
	10010010= 92= 146| Chan 3         "           |      "        |      "

	10000000= 80= 128| Chan 1      Note off       |  Note Number  | Note Velocity
	10000001= 81= 129| Chan 2         "           |   (0-127)     |   (0-127)
	10000010= 82= 130| Chan 3         "           |     see       |      "
*/

	channel1On = 144
	channel2On = 145
	channel3On = 146

	channel1Off = 128
	channel2Off = 129
	channel3Off = 130
)

type Device struct {
	name string
	input bool
	output bool
	open bool
}

var OnChannels = [3]int64{channel1On, channel2On, channel3On}
var OffChannels = [3]int64{channel1Off, channel2Off, channel3Off}

// TODO: load in another class

var OnToOffSet = make(map[int]int)
var OnSet = make(map[int]bool)
var OffSet = make(map[int]bool)

func initializeChannelConstants() {

	OnToOffSet[channel1On] = 128
	OnToOffSet[channel2On] = 129
	OnToOffSet[channel2On] = 130

	OnSet[channel1On] = true
	OnSet[channel2On] = true
	OnSet[channel3On] = true

	OffSet[channel1Off] = true
	OffSet[channel2Off] = true
	OffSet[channel3Off] = true
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

func isOnCommand(command int64) bool {
	if OnSet[int(command)] {
		return true
	} else {
		return false
	}
}

func isOffCommand(command int64) bool {
	if OffSet[int(command)] {
		return true
	} else {
		return false
	}
}

func onCommand(channel int) int64 {
	return OnChannels[channel]
}

func offCommand(channel int) int64 {
	return OffChannels[channel]
}

func readFromDeviceWriteToDevice(OutStream portmidi.Stream, InStream portmidi.Stream, maxEvents int, outChannel int) {
	var status int64

	for {
		events, err := OutStream.Read(maxEvents)
		if err != nil {
			log.Fatal(err)
		} else {
			var sendEvents = make([]portmidi.Event, len(events))
			if len(events) > 0 {
				for j := 0; j < len(events); j++ {
					fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)
					if isOnCommand(events[j].Status) {
						status = onCommand(outChannel)
					} else if isOffCommand(events[j].Status) {
						status = offCommand(outChannel)
					}
					sendEvents[j] = portmidi.Event{Timestamp: events[j].Timestamp, Status: status, Data1: events[j].Data1, Data2: events[j].Data2}
				}
				InStream.Write(sendEvents)
			}
		}
	}
}

func main(){
	fmt.Printf("Reading from midi channels\n")
	
	portmidi.Initialize()
	initializeChannelConstants()

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

	readFromDeviceWriteToDevice(*InStream, *OutStream, 10, 2)

}