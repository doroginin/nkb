package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/MarinX/keylogger"
)

const (
	KEY_0 = iota + 500
	KEY_1
)

var debug bool
var capsMode bool
var capsMode2 bool
var key3Pressed bool
var pressedButtonsCount int
var capsPressed bool

func main() {
	debug = false
	cmd := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var device = cmd.Int("d", 3, "Listen device with id")

	devices, err := keylogger.NewDevices()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := cmd.Parse(os.Args[1:]); err == flag.ErrHelp || *device < 0 || len(devices) <= *device {
		if err != flag.ErrHelp {
			cmd.Usage()
		}
		fmt.Println("\nAvailable devices: ")
		for _, val := range devices {
			fmt.Println("Id: ", val.Id, "Device: ", val.Name)
		}
		return
	}

	rd := keylogger.NewKeyLogger(devices[*device])
	in, err := rd.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	disableKeys()
	defer restoreKeys()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case i := <-in:
			if i.Type != keylogger.EV_KEY {
				continue
			}
			if debug {
				fmt.Printf("key:\t%s,\tcode:\t%d,\tevent:\t%d\n", i.KeyString(), i.Code, i.Value)
			}
			switch uint(i.Code) {
			case KEY_0: // CAPS LOCK
				if int(i.Value) == 1 {
					capsPressed = true
					if !capsMode {
						capsModeOn()
					}
				}
				if int(i.Value) == 0 {
					capsPressed = false
					if pressedButtonsCount <= 0 {
						pressedButtonsCount = 0
						capsModeOff()
					}
				}
			case KEY_1:
				if int(i.Value) == 1 {
					key3Pressed = true
					capsMode2On()
				}
				if int(i.Value) == 0 {
					key3Pressed = false
					if pressedButtonsCount <= 0 {
						pressedButtonsCount = 0
						capsMode2Off()
					}
				}
			case 103, 105, 106, 108, 102, 104, 107, 109, 111:
				if int(i.Value) == 1 {
					pressedButtonsCount++
				}
				if int(i.Value) == 0 {
					pressedButtonsCount--
					if pressedButtonsCount <= 0 {
						pressedButtonsCount = 0
						if !key3Pressed {
							capsMode2Off()
						}
						if !capsPressed {
							capsModeOff()
						}
					}
				}
			}
			if debug {
				fmt.Printf("capsMode: %t, capsMode2: %t, key3Pressed: %t, pressedButtonsCount: %d," +
					" capsPressed: %t'\n", capsMode, capsMode2, key3Pressed, pressedButtonsCount, capsPressed)
			}

		case <-quit:
			return
		}
	}
}

func disableKeys() {
	start("" +
		"setkeycodes 3a " + strconv.Itoa(KEY_0) + " &", // capslock
	)
}

func restoreKeys() {
	start("" +
		"setkeycodes 3a 58 &", // capslock
	)
}

func capsModeOn() {
	capsMode = true
	capsMode2 = false
	start("" +
		"setkeycodes 38 " + strconv.Itoa(KEY_1) + " &" + // alt
		"setkeycodes 0e 111 &" + // bksp -> del
		"setkeycodes 1e 105 &" + // a -> left
		"setkeycodes 11 103 &" + // w -> up
		"setkeycodes 20 106 &" + // d -> right
		"setkeycodes 1f 108 &" + // s -> down
		"")
}

func capsModeOff() {
	capsMode = false
	capsMode2 = false
	start("" +
		"setkeycodes 38 56 &" + // alt
		"setkeycodes 0e 14 &" + // bksp
		"setkeycodes 1e 30 &" + // a -> left
		"setkeycodes 11 17 &" + // w -> up
		"setkeycodes 20 32 &" + // d -> right
		"setkeycodes 1f 31 &" + // s -> down
		"")
}

func capsMode2On() {
	capsMode = false
	capsMode2 = true
	start("" +
		"setkeycodes 1e 102 &" + // a -> home
		"setkeycodes 11 104 &" + // w -> pgup
		"setkeycodes 20 107 &" + // d -> end
		"setkeycodes 1f 109 &" + // s -> pgdn
		"",
	)
}

func capsMode2Off() {
	capsMode = true
	capsMode2 = false
	start("" +
		"setkeycodes 1e 105 &" + // a -> left
		"setkeycodes 11 103 &" + // w -> up
		"setkeycodes 20 106 &" + // d -> right
		"setkeycodes 1f 108 &" + // s -> down
		"",
	)
}

func send(keys string) {
	start("xdotool key " + keys)
}

func start(cmd string) {
	if _, err := exec.Command("/bin/sh", "-c", cmd).Output(); err != nil {
		log.Printf("Error during exec command: %s", err.Error())
	}
}
