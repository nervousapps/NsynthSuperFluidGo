module nsynthsuperfluid

go 1.17

require (
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	periph.io/x/devices/v3 v3.6.11
)

require (
	github.com/rainu/launchpad v1.2.0 // indirect
	github.com/rakyll/portmidi v0.0.0-20191102002215-74e95e8bc9b1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/stianeikeland/go-rpio/v4 v4.4.0 // indirect
	gitlab.com/gomidi/midi v1.21.0 // indirect
	gitlab.com/gomidi/rtmididrv v0.14.0 // indirect
	golang.org/x/sys v0.0.0-20210525143221-35b2ab0089ea // indirect
	periph.io/x/conn/v3 v3.6.8 // indirect
	periph.io/x/d2xx v0.0.1 // indirect
	periph.io/x/host/v3 v3.7.0 // indirect
)

require (
	buttons v1.0.0
	fluidsynth2 v0.1.0
	github.com/go-daq/smbus v0.0.0-20201216173259-5725b4593606
	github.com/rakyll/launchpad v0.1.0
	github.com/raspberrypi-go-drivers/button v0.0.0-20201218194336-d7e7fbd8a9e3
	github.com/xthexder/go-jack v0.0.0-20201026211055-5b07fb071116
	inputs v1.0.0
	screen v1.0.0
)

replace (
	buttons v1.0.0 => ./buttons
	fluidsynth2 v0.1.0 => ./fluidsynth2@v0.1.0
	inputs v1.0.0 => ./inputs
	screen v1.0.0 => ./screen
)
