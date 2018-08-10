package main

import (
	"fmt"
	"github.com/rakyll/portmidi"
	"log"
	"github.com/GeorgeLuo/MIDIControllerTools/structures"
	"github.com/GeorgeLuo/MIDIControllerTools/loops"
	"github.com/GeorgeLuo/MIDIControllerTools/environment"
	"github.com/spf13/viper"
)

// import "sync"

/*
	Design:
	EnvironmentSession: map of devices

	SDK functions:
	ListenAll: output array of event -> instrument mappings
	map device toggles to any device toggle

	TODO: should close stream at end of session, otherwise
	switching instruments can lead to irregular behavior
	TODO: read yaml config to initialize mapping
	TODO: make loops async to switch channels outside of
	tunneling loop
*/

func main() {
	// fmt.Printf("Reading from midi channels\n")

	// initializeChannelConstants()

	// var numDevices = portmidi.CountDevices() // returns the number of MIDI devices
	// fmt.Printf("Num devices: %d\n\n", numDevices)
	// // portmidi.Info(deviceID) // returns info about a MIDI device

	// for i := 0; i < numDevices; i++ {
	// 	var portInfo = portmidi.Info(portmidi.DeviceID(i))
	// 	fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name,
	// 		portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)
	// }

	// var defaultInputDeviceID = portmidi.DefaultInputDeviceID() // returns the ID of the system default input
	// fmt.Printf("Default input device ID: %d\n", defaultInputDeviceID)

	// var defaultOutDeviceID = portmidi.DefaultOutputDeviceID() // returns the ID of the system default output
	// fmt.Printf("Default output device ID: %d\n\n", defaultOutDeviceID)

	// InStreamOrigin62, err := portmidi.NewInputStream(portmidi.DeviceID(1), 1024)

	// InStreamJDXI, err := portmidi.NewInputStream(portmidi.DeviceID(2), 1024)
	// var jdxiIgnoreMap = make(map[int64]bool)
	// jdxiIgnoreMap[248] = true

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var originToJDXINoteMap = make(map[int64]int64)
	// var originToJDXIChannelMap = make(map[int64]int64)

	// originToJDXINoteMap[70] = 102
	// originToJDXIChannelMap[176] = 177

	// var jobQueue = NewJobQueue()

	portmidi.Initialize() // need to run here or else the streams are out of scope
	origin62mapping.InitializeOrigin62Mapping()
	hardWareMap := initializeHardwareMap()
	tunnelMap := initializeTunnelingMap(hardWareMap)
	startParallelize(10, tunnelMap)

	// OutStreamJDXI, err := portmidi.NewOutputStream(portmidi.DeviceID(6), 1024, 0)
	// readFromDeviceWriteToDevice(*InStreamOrigin62, *OutStreamJDXI, originToJDXINoteMap, originToJDXIChannelMap, 10, 2)

	// var inStreams = []StreamWrapper{StreamWrapper{underStream: *InStreamJDXI, ignoreStatus: jdxiIgnoreMap}, StreamWrapper{underStream: *InStreamOrigin62, ignoreStatus: make(map[int64]bool)}}
	// readFromIncomingDevices(inStreams, 10)

}
