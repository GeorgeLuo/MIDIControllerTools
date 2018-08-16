package environment

import (
	"github.com/rakyll/portmidi"
)

var EnvironmentDeviceLayout DeviceLayout

var WatchPort = map[int] bool {
	1 : true,
	2 : true,
}

var IgnoreChannel = map[int64] bool {
	248 : true,
}

var SourcePortToStreams = make(map[int]portmidi.Stream)
var DestinationPortToStreams = make(map[int]portmidi.Stream)

var PortNoteChannelMap = make(map[PortNoteChannel]PortNoteChannel)

var GlobalJobQueue = NewJobQueue()

var SourcePortToDestinationPortMap = make(map[int]int)