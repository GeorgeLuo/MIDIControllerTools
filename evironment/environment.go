package environment

var deviceLayout DeviceLayout

const (

	/*
		declare constants for channel 1-3
		10010000= 90= 144| Chan 1      Note on        |      "        |      "
		10010001= 91= 145| Chan 2         "           |      "        |      "
		10010010= 92= 146| Chan 3         "           |      "        |      "

		10000000= 80= 128| Chan 1      Note off       |  Note Number  | Note Velocity
		10000001= 81= 129| Chan 2         "           |   (0-127)     |   (0-127)
		10000010= 82= 130| Chan 3         "           |     see       |      "
	*/

	channel1On = 144
	channel2On = 145
	channel3On = 146

	channel1Off = 128
	channel2Off = 129
	channel3Off = 130

	cutoffFilter = 102
)

var OnToOffSet = make(map[int]int)
var OnSet = make(map[int]bool)
var OffSet = make(map[int]bool)

var Filters = make(map[int64]bool)

var OnChannels = [3]int64{channel1On, channel2On, channel3On}
var OffChannels = [3]int64{channel1Off, channel2Off, channel3Off}

func initializeChannelConstants() {

	OnToOffSet[channel1On] = 128
	OnToOffSet[channel2On] = 129
	OnToOffSet[channel2On] = 130

	OnSet[channel1On] = true
	OnSet[channel2On] = true
	OnSet[channel3On] = true

	OffSet[channel1Off] = true
	OffSet[channel2Off] = true
	OffSet[channel3Off] = true
	Filters[cutoffFilter] = true
}

func importDeviceConfig() {

}

func initializeDeviceLayout() {
	fmt.Println("initializing device layout")
	var numDevices = portmidi.CountDevices()
	deviceList := make([]Device, numDevices)
	for i := 0; i < numDevices; i++ {
		var portInfo = portmidi.Info(portmidi.DeviceID(i))
		fmt.Printf("Port %d:\nName: %s\nInput Available: %t\nOutput Available: %t\nIs Open: %t\n\n", i, portInfo.Name,
			portInfo.IsInputAvailable, portInfo.IsOutputAvailable, portInfo.IsOpened)

		deviceList[i] = Device{port: i, name: portInfo.Name, input: portInfo.IsInputAvailable, output: portInfo.IsOutputAvailable, open: portInfo.IsOpened}
	}
	deviceLayout = DeviceLayout{devices: deviceList}

}

/* change mode between
	0 = play mode using mapping
	1 = map control mode, set by combination of midi events (maybe push two control buttons)
*/
var globalMode = 0
func changeMode(mode int) {
	// var modeOn = true
	// for modeOn {

	// 	modeOn
	// }
	globalMode = mode
}