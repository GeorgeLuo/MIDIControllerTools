package environment

import (
	"github.com/rakyll/portmidi"
)

var EnvironmentDeviceLayout DeviceLayout

var WatchPort = map[int] bool {
	1 : true,
	2 : true,
}

// ie. watch "JD-Xi" and "Origin62"
var WatchPortByString = map[string] bool {
	"JD-Xi" : true,
	"Origin62" : true,
}

// ie. output from Origin62 goes to input of JD-Xi
var DeviceToDeviceMap = map[string] string {
	"Origin62" : "JD-Xi",
}

var SourceDeviceToPorts = map[string] int {

}

var DestinationDeviceToPort = map[string] int {

}

var IgnoreChannel = map[int64] bool {
	248 : true,
}

var SourcePortToStreams = make(map[int]portmidi.Stream)
var DestinationPortToStreams = make(map[int]portmidi.Stream)

var PortNoteChannelMap = make(map[PortNoteChannel]PortNoteChannel)

var GlobalJobQueue = NewJobQueue()

var SourcePortToDestinationPortMap = make(map[int]int)