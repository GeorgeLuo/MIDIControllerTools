package main

import (
	"github.com/rakyll/portmidi"
	"github.com/GeorgeLuo/MIDIControllerTools/environment"
)

func main() {
	portmidi.Terminate()
	portmidi.Initialize()

	environment.InitializeDeviceLayout()

	environment.InitializePortToPortMap()
	environment.InitializePortNoteChannelMapping()

	environment.StartParallelize(10)
}