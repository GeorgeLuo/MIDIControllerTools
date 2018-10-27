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
	// environment.InitializePortNoteChannelMapping()

	environment.ReadMappingConfig("config/channel_1_master.json")

	environment.StartParallelize(10)
}