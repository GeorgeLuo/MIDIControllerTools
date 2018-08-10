package structures

type KnobToKnobMapping struct {
	FromDevice int64
	ToDevice   int64
	FromKnob   int64
	ToKnob     int64
}

type NoteMap struct {
	inControlNote  int64
	outControlNote int64
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