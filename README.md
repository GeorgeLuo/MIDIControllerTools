# MIDIControllerTools

A tool to tunnel MIDI instructions between devices. The primary application is using a dumb MIDI controller as a master to a smarter device (multi-channel synth here), solving the problem of only being able to access the interface of one channel at a time.

Testing implementation optimized for JD-XI as output device.

TODO: initialize using config files.
TODO: implement a mapping mode to change control mappings dynamically.
TODO: port application to micro-controller.

Reference for MIDI mapping standards:
http://www.logosfoundation.org/kursus/1075.html
