package environment

import (
	"github.com/rakyll/portmidi"
	"fmt"
)

var GlobalMode = 0

func InitializeDeviceLayout() {
	fmt.Println("initializing device layout")
	var numDevices = portmidi.CountDevices()
	deviceList := make([]Device, numDevices)
	for i := 0; i < numDevices; i++ {
		var portInfo = portmidi.Info(portmidi.DeviceID(i))

		deviceList[i] = Device{port: i, name: portInfo.Name, input: portInfo.IsInputAvailable, output: portInfo.IsOutputAvailable, open: portInfo.IsOpened}

		if portInfo.IsInputAvailable && WatchPortByString[portInfo.Name] {
			tempStream, err := portmidi.NewInputStream(portmidi.DeviceID(i), 1024)
			if(err != nil) {
				fmt.Errorf("error initializing input stream %g", err)
			}
			var portInfo = portmidi.Info(portmidi.DeviceID(i))
			fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name,
				portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)

			SourcePortToStreams[i] = *tempStream
			SourceDeviceToPorts[portInfo.Name] = i
			fmt.Println(portInfo.Name, "set to port", i, "at", tempStream)
		}

		if(portInfo.IsOutputAvailable) {
			tempStream, err := portmidi.NewOutputStream(portmidi.DeviceID(i), 1024, 0)
			if(err != nil) {
				fmt.Errorf("error initializing input stream %g", err)
			}
			var portInfo = portmidi.Info(portmidi.DeviceID(i))
			fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name,
				portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)

			DestinationPortToStreams[i] = *tempStream
			DestinationDeviceToPort[portInfo.Name] = i
		}
	}

	// now initialize mapping by name of device
	for sourceDevice, port := range SourceDeviceToPorts {
		outputDevice, exists := DeviceToDeviceMap[sourceDevice]
		if exists {
			SourcePortToDestinationPortMap[port] = DestinationDeviceToPort[outputDevice]
		}
	}


	fmt.Println("SourcePortToStreams", SourcePortToStreams)
	EnvironmentDeviceLayout = DeviceLayout{devices: deviceList}
}

func InitializePortToPortMap() {
	// SourcePortToDestinationPortMap[1] = 6
}

func ReadMappingConfig() {

}

func InitializePortNoteChannelMapping() {
	// 176 means channel 1, 177 = channel 2, 178 = channel 3. Note mappings indicate knob to knob.
	// JD-XI note 102 is cutoff filter. note 98

	PortNoteChannelMap[PortNoteChannel{port: 1, note: 70, channel: 176}] = PortNoteChannel{port: 7, note: 102, channel: 176}
	PortNoteChannelMap[PortNoteChannel{port: 1, note: 71, channel: 176}] = PortNoteChannel{port: 7, note: 102, channel: 177}
	PortNoteChannelMap[PortNoteChannel{port: 1, note: 72, channel: 176}] = PortNoteChannel{port: 7, note: 102, channel: 178}
	PortNoteChannelMap[PortNoteChannel{port: 1, note: 73, channel: 176}] = PortNoteChannel{port: 7, note: 102, channel: 185}

	PortNoteChannelMap[PortNoteChannel{port: 1, note: 74, channel: 176}] = PortNoteChannel{port: 7, note: 98, channel: 176}
	PortNoteChannelMap[PortNoteChannel{port: 1, note: 75, channel: 176}] = PortNoteChannel{port: 7, note: 98, channel: 177}
	PortNoteChannelMap[PortNoteChannel{port: 1, note: 76, channel: 176}] = PortNoteChannel{port: 7, note: 98, channel: 178}
	PortNoteChannelMap[PortNoteChannel{port: 1, note: 77, channel: 176}] = PortNoteChannel{port: 7, note: 98, channel: 185}
}