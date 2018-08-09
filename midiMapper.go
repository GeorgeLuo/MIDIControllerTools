package main

import (
	"fmt"
	"github.com/rakyll/portmidi"
	"log"
	"github.com/GeorgeLuo/MIDIControllerTools/structures"
	// "github.com/spf13/viper"
)

import "sync"

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

	cutoffFilter = 102
)

type NoteMap struct {
	inControlNote  int64
	outControlNote int64
}

type Device struct {
	port   int
	name   string
	input  bool
	output bool
	open   bool
}

type StreamWrapper struct {
	ignoreStatus map[int64]bool
	underStream  portmidi.Stream
	portNum int
}

var OnChannels = [3]int64{channel1On, channel2On, channel3On}
var OffChannels = [3]int64{channel1Off, channel2Off, channel3Off}

// TODO: load in another class

var OnToOffSet = make(map[int]int)
var OnSet = make(map[int]bool)
var OffSet = make(map[int]bool)

var Filters = make(map[int64]bool)

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
	Filters[cutoffFilter] = true
}

func importDeviceConfig() {

}

func initializeDeviceLayout() {
	fmt.Println("initializing device layout")
	var numDevices = portmidi.CountDevices()
	deviceList := make([]Device, numDevices)
	for i := 0; i < numDevices; i++ {
		var portInfo = portmidi.Info(portmidi.DeviceID(i))
		fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name,
			portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)

		deviceList[i] = Device{port: i, name: portInfo.Name, input: portInfo.IsInputAvailable, output: portInfo.IsOutputAvailable, open: portInfo.IsOpened}
	}
	deviceLayout = DeviceLayout{devices: deviceList}

}

type KnobToKnobMapping struct {
	FromDevice int64
	ToDevice   int64
	FromKnob   int64
	ToKnob     int64
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

func readFromIncomingDevices(streams []StreamWrapper, maxEvents int) {
	fmt.Println(len(streams))
	readFunctions := make([]func(), len(streams))
	for i := 0; i < len(streams); i++ {
		// fmt.Println("loop", i)
		var funcNum = i // need to initialize value otherwise Parallel uses the final value of i
		readFunctions[funcNum] = func() {
			// fmt.Println("setting func", funcNum)
			readFromIncomingDevice(streams[funcNum], maxEvents)
		}
	}
	fmt.Println(readFunctions)

	Parallelize(readFunctions)
} 

// get all input streams from device layout
func getAllInputStreams() []StreamWrapper {
	var inputStreams = make([]StreamWrapper, 0)
	initializeDeviceLayout()


	fmt.Println("listening from Devices:")
	for i := 0; i < len(deviceLayout.devices); i++ {
		tempDevice := deviceLayout.devices[i]
		if tempDevice.input == true {
			fmt.Println(tempDevice.name, "stream", tempDevice.port, portmidi.DeviceID(tempDevice.port))
			tempStream, _ := portmidi.NewInputStream(portmidi.DeviceID(tempDevice.port), 1024)
			fmt.Println(tempStream)
			inputStreams = append(inputStreams, StreamWrapper{underStream: *tempStream, portNum: tempDevice.port, ignoreStatus: make(map[int64]bool)})
		}
	}

	return inputStreams
}

var jobQueue = NewJobQueue()


// func Rule struct {
// 	// when mode
// 	mode int
// 	S
// 	action func()
// }

func executeDispatchJobs() {
	OutStreamJDXI, _ := portmidi.NewOutputStream(portmidi.DeviceID(6), 1024, 0)
	// fmt.Println("processing events", midiEvents.commandSource(), midiEvents.events(), OutStreamJDXI)

	for {
		if jobQueue.Len() > 0 {
			midiEvents := jobQueue.Poll()
			if globalMode == 0 {
				tunnelData (midiEvents.events(), *OutStreamJDXI, 2)
			}
		}
	}
}

// start listening for jobs, streams contains all input streams detected
func startParallelize(maxEvents int) {
	loops := make([]func(), 0)

	dispatchJobs := func() {
		executeDispatchJobs()
	}

	loops = append(loops, dispatchJobs)
	readChannels := func() {
		readFromIncomingDevices(getAllInputStreams(), 10)
	}
	loops = append(loops, readChannels)

	Parallelize(loops)
}

/* change mode between
	0 = play mode using mapping
	1 = map control mode, set by combination of midi events (maybe push two control buttons)
*/
var globalMode = 0
func changeMode(mode int) {
	// var modeOn = true
	// for modeOn {

	// 	modeOn
	// }
	globalMode = mode
}

func (j deviceJob) commandSource() int {
    return j.source
}

func (j deviceJob) events() []portmidi.Event {
    return j.inputEvents
}

func readFromIncomingDevice(stream StreamWrapper, maxEvents int) {
	for {
		events, err := stream.underStream.Read(maxEvents)
		if err != nil {
			log.Fatal(err)
		} else {
			for j := 0; j < len(events); j++ {
				if !stream.ignoreStatus[events[j].Status] {
					// fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)
					toQueue := deviceJob{source: stream.portNum, inputEvents: events}
					// next = jobInterface{commandSource: j, events: events}
					jobQueue.Push(toQueue)
				}
			}
		}
	}
}

func isOnCommand(command int64) bool {
	if OnSet[int(command)] {
		return true
	} else {
		return false
	}
}

// TODO: support both protocols for on-off, command and 0 velocity
func isOffCommand(command int64) bool {
	if OffSet[int(command)] {
		return true
	} else {
		return false
	}
}

func onCommand(channel int) int64 {
	return OnChannels[channel-1]
}

func offCommand(channel int) int64 {
	return OffChannels[channel-1]
}

// TODO: account for note = 0
func noteMapContains(noteMap map[int64]int64, check int64) bool {
	if noteMap[check] == 0 {
		return false
	} else {
		return true
	}
}

func isFilter(note int64) bool {
	return Filters[note]
}

/*
	TODO: change streams to channels, will reduce number of parameters by determining readable channel to MIDI channel.
*/
func readFromDeviceWriteToDevice(OutStream portmidi.Stream, InStream portmidi.Stream, noteMap map[int64]int64, channelMap map[int64]int64, maxEvents int, outChannel int) {
	var status int64
	fmt.Println(outChannel)
	for {
		events, err := OutStream.Read(maxEvents)
		if err != nil {
			log.Fatal(err)
		} else {
			var sendEvents = make([]portmidi.Event, 0)
			if len(events) > 0 {
				fmt.Println("num events", len(events))
				for j := 0; j < len(events); j++ {
					fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)

					if noteMapContains(noteMap, events[j].Data1) {
						// rewrite note (knob or slide)
						status = channelMap[events[j].Status]
						note := noteMap[events[j].Data1]
						sendEvents = append(sendEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: status, Data1: note, Data2: events[j].Data2})
						fmt.Println("sending: ", events[j].Timestamp, status, note, events[j].Data2)
						if isFilter(noteMap[events[j].Data1]) {
							sendEvents = append(sendEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: status, Data1: note + int64(1), Data2: events[j].Data2})
							fmt.Println("sending: ", events[j].Timestamp, status, note+int64(1), events[j].Data2)
						}
					} else {
						if isOnCommand(events[j].Status) {
							status = onCommand(outChannel)
						} else if isOffCommand(events[j].Status) {
							status = offCommand(outChannel)
						}
						sendEvents = append(sendEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: status, Data1: events[j].Data1, Data2: events[j].Data2})
					}
				}
				InStream.Write(sendEvents)
			}
		}
	}
}

func tunnelData (events []portmidi.Event, InStream portmidi.Stream, outChannel int) {
	var sendEvents = make([]portmidi.Event, 0)
	for j := 0; j < len(events); j++ {
		// fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)
		sendEvents = append(sendEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: events[j].Status, Data1: events[j].Data1, Data2: events[j].Data2})
	}
	// fmt.Println(sendEvents, InStream)
	InStream.Write(sendEvents)
}

func Parallelize(functions []func()) {
	var waitGroup sync.WaitGroup
	for i := 0; i < len(functions); i++ {
		waitGroup.Add(len(functions))
	}

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

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

	portmidi.Initialize() // need to run here or else the streams are out of scope
	startParallelize(10)

	// OutStreamJDXI, err := portmidi.NewOutputStream(portmidi.DeviceID(6), 1024, 0)
	// readFromDeviceWriteToDevice(*InStreamOrigin62, *OutStreamJDXI, originToJDXINoteMap, originToJDXIChannelMap, 10, 2)

	// var inStreams = []StreamWrapper{StreamWrapper{underStream: *InStreamJDXI, ignoreStatus: jdxiIgnoreMap}, StreamWrapper{underStream: *InStreamOrigin62, ignoreStatus: make(map[int64]bool)}}
	// readFromIncomingDevices(inStreams, 10)

}
