package environment

import (
	"github.com/rakyll/portmidi"
	"sync"
	"fmt"
	"log"
)

// start listening for jobs, streams contains all input streams detected
func StartParallelize(maxEvents int) {
	loops := make([]func(), 0)

	dispatchJobs := func() {
		ExecuteDispatchJobs()
	}

	loops = append(loops, dispatchJobs)
	readChannels := func() {
		ReadFromIncomingDevices(10)
	}
	loops = append(loops, readChannels)

	Parallelize(loops)
}

func ExecuteDispatchJobs() {

	for {
		if GlobalJobQueue.Len() > 0 {
			deviceJob := GlobalJobQueue.Poll()
			if GlobalMode == 0 {
				TunnelData (deviceJob.events(), deviceJob.commandDestination())
			}
		}
	}
}

func TunnelData(events []portmidi.Event, port int) {
	fmt.Println("tunneling", events, "to port", port)
	var destinationStream = DestinationPortToStreams[port]
	destinationStream.Write(events)
}

func ReadFromIncomingDevices(maxEvents int) {
	var i = 0
	readFunctions := make([]func(), len(SourcePortToStreams))

	for port, sourceStream := range SourcePortToStreams { 
		// fmt.Println("loop", i)
		var funcNum = i // need to initialize value otherwise Parallel uses the final value of i
		var portNum = port
		var tempSourceStream = sourceStream
		readFunctions[funcNum] = func() {
			// fmt.Println("setting func", funcNum)
			ReadFromIncomingDevice(StreamWrapper{underStream: tempSourceStream, portNum: portNum}, maxEvents)
		}

		i++
	}
	fmt.Println(readFunctions)
	Parallelize(readFunctions)
} 

// read from device, translate inbound events to outbound events, add outbound events to jobQueue
func ReadFromIncomingDevice(stream StreamWrapper, maxEvents int) {
	fmt.Println(stream.portNum, stream.underStream, "sees", SourcePortToStreams)
	var outboundPort int
	for {
		events, err := stream.underStream.Read(maxEvents)
		if err != nil {
			log.Fatal(err)
		} else {
			var outboundEvents = make([]portmidi.Event, 0)
			for j := 0; j < len(events); j++ {
				if stream.portNum != 2 {
					fmt.Println(stream.portNum, events)
				} else {
					if !IgnoreChannel[events[j].Status] {
						fmt.Println(stream.portNum, events)
					}
				}
				
				destinationPnc, exists := CheckForReceivedEventMapping(stream.portNum, events[j].Data1, events[j].Status)
				if exists {
					fmt.Println("mapping found in ReceivedEventMapping", destinationPnc)
					outboundEvents = append(outboundEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: destinationPnc.channel, Data1: destinationPnc.note, Data2: events[j].Data2})
					outboundPort = destinationPnc.port
				} else {
					p, exists := GetPortMapping(stream.portNum, events[j].Status)
					if exists {
						fmt.Println("mapping found in PortMapping for port", stream.portNum, "sending to port", p, "event", events[j])
						outboundEvents = append(outboundEvents, portmidi.Event{Timestamp: events[j].Timestamp, Status: events[j].Status, Data1: events[j].Data1, Data2: events[j].Data2})
						outboundPort = p
					}

					// if control is not mapped and the inbound port is not mapped to another port, don't do anything
				} 
			}
			if len(outboundEvents) > 0 {
				toQueue := DeviceJob{destination: outboundPort, outboundEvents: outboundEvents}
				GlobalJobQueue.Push(toQueue)
			}
		}
	}
}

// given note, channel, and port, return new event, if no mapping exists = false
func CheckForReceivedEventMapping(port int, note int64, channel int64) (PortNoteChannel, bool) {
	pnc, exists := PortNoteChannelMap[PortNoteChannel{port: port, note: note, channel: channel}]
	return pnc, exists
}

func GetPortMapping(port int, channel int64) (int, bool) {
	p, exists := SourcePortToDestinationPortMap[port]
	exists = exists && !IgnoreChannel[channel]
	return p, exists
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
