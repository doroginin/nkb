package nkb

import (
	"fmt"
	"time"
	"github.com/MarinX/keylogger"
	"flag"
	"os"
)

type state struct {
	enabled bool
	capsMode bool
	capsMode2 bool
	caps2Pressed bool
	pressedButtons map[uint]struct{}
	capsPressed bool
	lastCapsPressed time.Time
}

type app struct {
	devices []*keylogger.InputDevice
	device int
	verbose bool
	ch chan keylogger.InputEvent
	state state
}

func New() (*app, error) {
	app := &app{state:state{
		pressedButtons:make(map[uint]struct{}),
	}}
	devices, err := keylogger.NewDevices()
	if err != nil {
		return nil, err
	}
	app.devices = devices

	if !app.flags() {
		return nil, fmt.Errorf("bad usage")
	}
	if err := app.prepare(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *app) flags() bool {
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	verbose := cmd.Bool("v", false, "Verbose mode")
	device := cmd.Int("d", 3, "Listen device with id")
	if err := cmd.Parse(os.Args[1:]); err == flag.ErrHelp || *device < 0 || len(a.devices) <= *device {
		if err != flag.ErrHelp {
			cmd.Usage()
		}
		fmt.Println("\nAvailable devices: ")
		for _, val := range a.devices {
			fmt.Println("Id: ", val.Id, "Device: ", val.Name)
		}
		return false
	}
	a.device = *device
	a.verbose = *verbose
	return true
}

func (a *app) prepare() error {
	rd := keylogger.NewKeyLogger(a.devices[a.device])
	in, err := rd.Read()
	if err != nil {
		return err
	}
	a.ch = in
	return nil
}
