package structures


type DeviceLayout struct {
	devices []Device
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

func ReadFromIncomingDevice(stream StreamWrapper, maxEvents int) {
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

func TunnelData(events []portmidi.Event, InStream portmidi.Stream, outChannel int) {
	var sendEvents = make([]portmidi.Event, 0)
	for j := 0; j < len(events); j++ {
		// fmt.Println(events[j].Timestamp, events[j].Status, events[j].Data1, events[j].Data2)
		sendEvents = append(sendEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: events[j].Status, Data1: events[j].Data1, Data2: events[j].Data2})
	}
	// fmt.Println(sendEvents, InStream)
	InStream.Write(sendEvents)
}

