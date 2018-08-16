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