package environment

var Origin62ControlNameToNote = make(map[int64]string)
var Origin62ControlNoteToName = make(map[string]int64)

var InitializeOrigin62Mapping() {

	// Knobs
	Origin62ControlNoteToName[70] = "K1"
	Origin62ControlNoteToName[71] = "K2"
	Origin62ControlNoteToName[72] = "K3"
	Origin62ControlNoteToName[73] = "K4"
	Origin62ControlNoteToName[74] = "K5"
	Origin62ControlNoteToName[75] = "K6"
	Origin62ControlNoteToName[76] = "K7"
	Origin62ControlNoteToName[77] = "K8"

	Origin62ControlNameToNote["K1"] = 70
	Origin62ControlNameToNote["K2"] = 71
	Origin62ControlNameToNote["K3"] = 72
	Origin62ControlNameToNote["K4"] = 73
	Origin62ControlNameToNote["K5"] = 74
	Origin62ControlNameToNote["K6"] = 75
	Origin62ControlNameToNote["K7"] = 76
	Origin62ControlNameToNote["K8"] = 77

	// Sliders
	Origin62ControlNoteToName[20] = "S1"
	Origin62ControlNoteToName[21] = "S2"
	Origin62ControlNoteToName[22] = "S3"
	Origin62ControlNoteToName[23] = "S4"
	Origin62ControlNoteToName[24] = "S5"
	Origin62ControlNoteToName[25] = "S6"
	Origin62ControlNoteToName[26] = "S7"
	Origin62ControlNoteToName[27] = "S8"

	Origin62ControlNameToNote["S1"] = 20
	Origin62ControlNameToNote["S2"] = 21
	Origin62ControlNameToNote["S3"] = 22
	Origin62ControlNameToNote["S4"] = 23
	Origin62ControlNameToNote["S5"] = 24
	Origin62ControlNameToNote["S6"] = 25
	Origin62ControlNameToNote["S7"] = 26
	Origin62ControlNameToNote["S8"] = 27

	// System Controls
	Origin62ControlNoteToName[1] = "C4"
	Origin62ControlNoteToName[2] = "C5"
	Origin62ControlNoteToName[3] = "C1"
	Origin62ControlNoteToName[4] = "C3"
	Origin62ControlNoteToName[5] = "C2"
	Origin62ControlNoteToName[6] = "C6"

	Origin62ControlNameToNote["C4"] = 1
	Origin62ControlNameToNote["C5"] = 2
	Origin62ControlNameToNote["C1"] = 3
	Origin62ControlNameToNote["C3"] = 4
	Origin62ControlNameToNote["S2"] = 5
	Origin62ControlNameToNote["S6"] = 6

	Origin62ControlNoteToName[91] = "D1"
	Origin62ControlNameToNote["D1"] = 91
}

