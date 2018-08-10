package loops

func readFromIncomingDevices(streams []StreamWrapper, maxEvents int) {
	fmt.Println(len(streams))
	readFunctions := make([]func(), len(streams))
	for i := 0; i < len(streams); i++ {
		// fmt.Println("loop", i)
		var funcNum = i // need to initialize value otherwise Parallel uses the final value of i
		readFunctions[funcNum] = func() {
			// fmt.Println("setting func", funcNum)
			readFromIncomingDevice(streams[funcNum], maxEvents)
		}
	}
	fmt.Println(readFunctions)

	Parallelize(readFunctions)
} 

// get all input streams from device layout
func getAllInputStreams() []StreamWrapper {
	var inputStreams = make([]StreamWrapper, 0)
	initializeDeviceLayout()


	fmt.Println("listening from Devices:")
	for i := 0; i < len(deviceLayout.devices); i++ {
		tempDevice := deviceLayout.devices[i]
		if tempDevice.input == true {
			fmt.Println(tempDevice.name, "stream", tempDevice.port, portmidi.DeviceID(tempDevice.port))
			tempStream, _ := portmidi.NewInputStream(portmidi.DeviceID(tempDevice.port), 1024)
			fmt.Println(tempStream)
			inputStreams = append(inputStreams, StreamWrapper{underStream: *tempStream, portNum: tempDevice.port, ignoreStatus: make(map[int64]bool)})
		}
	}

	return inputStreams
}

func executeDispatchJobs() {
	OutStreamJDXI, _ := portmidi.NewOutputStream(portmidi.DeviceID(6), 1024, 0)
	// fmt.Println("processing events", midiEvents.commandSource(), midiEvents.events(), OutStreamJDXI)

	for {
		if jobQueue.Len() > 0 {
			midiEvents := jobQueue.Poll()
			if globalMode == 0 {
				tunnelData (midiEvents.events(), *OutStreamJDXI, 2)
			}
		}
	}
}

// start listening for jobs, streams contains all input streams detected
func startParallelize(maxEvents int) {
	loops := make([]func(), 0)

	dispatchJobs := func() {
		executeDispatchJobs()
	}

	loops = append(loops, dispatchJobs)
	readChannels := func() {
		readFromIncomingDevices(getAllInputStreams(), 10)
	}
	loops = append(loops, readChannels)

	Parallelize(loops)
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

