package environment

import (
	"github.com/rakyll/portmidi"
)

type Device struct {
	port   int
	name   string
	input  bool
	output bool
	open   bool
}

type DeviceLayout struct {
	devices []Device
}

type PortNoteChannel struct {
	port int
	note int64
	channel int64
}

type StreamWrapper struct {
	underStream  portmidi.Stream
	portNum int
}


type DeviceToDeviceMapping struct {
	DToDeMap []DToD `json:"DeviceToDeviceMap"`
}

type DToD struct {
	SourceDevice   string `json:"sourceDevice"`
	SourceNote   int64 `json:"sourceNote"`
	SourceChannel    int64    `json:"sourceChannel"`
	DestinationDevice   string `json:"destinationDevice"`
	DestinationNote   int64 `json:"destinationNote"`
	DestinationChannel    int64    `json:"destinationChannel"`
}