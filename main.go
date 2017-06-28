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
	"time"
)

const (
	KEY_0 = iota + 500
	KEY_1
)

var debug bool
var capsMode bool
var capsMode2 bool
var key1Pressed bool
var pressedButtonsCount int
var capsPressed bool
var lastCapsPressed time.Time

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
				fmt.Printf("key: %s, code: %d, event: %d\n", i.KeyString(), i.Code, i.Value)
			}

			// @todo remove this workaround
			if pressedButtonsCount < 0 {
				fmt.Print("keywatcher: workaround [pressedButtonsCount < 0]\n")
				pressedButtonsCount = 0
			}

			switch uint(i.Code) {
			case KEY_0: // CAPS LOCK
				if int(i.Value) == 1 {

					// @todo remove this workaround
					if time.Since(lastCapsPressed) < 300 * time.Millisecond {
						fmt.Print("keywatcher: workaround [double caps pressed]\n")
						pressedButtonsCount = 0
					}

					lastCapsPressed = time.Now()
					capsPressed = true
					if pressedButtonsCount == 0 {
						capsModeOn()
					}
				}
				if int(i.Value) == 0 {
					capsPressed = false
					if pressedButtonsCount == 0 {
						capsModeOff()
					}
				}
			case KEY_1:
				if int(i.Value) == 1 {
					key1Pressed = true
					if pressedButtonsCount == 0 {
						capsMode2On()
					}
				}
				if int(i.Value) == 0 {
					key1Pressed = false
					if pressedButtonsCount == 0 {
						capsModeOn()
					}
				}
			case
				// arrows
				103, 105, 106, 108,
				//home, end, pgup, pgdn
				102, 104, 107, 109,
				// delete, bksp
				111, 14,
				// a, s, d, w
				30,	17,	32,	31,
				// alt
				56:
				if int(i.Value) == 1 {
					pressedButtonsCount++
				}
				if int(i.Value) == 0 {
					pressedButtonsCount--
					if pressedButtonsCount == 0 {
						if capsMode2 && !key1Pressed && !capsPressed ||
							capsMode && !capsPressed {
							capsModeOff()
						}
						if capsPressed {
							if key1Pressed {
								capsMode2On()
							} else {
								capsModeOn()
							}
						}
					}
				}
			}
			if debug {
				fmt.Printf("capsPressed: %t, key1Pressed: %t, capsMode: %t, capsMode2: %t, pressedButtonsCount: %d\n",
					capsPressed, key1Pressed, capsMode, capsMode2, pressedButtonsCount)
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

// disable capsMode and capsMode2
func capsModeOff() {
	capsMode = false
	capsMode2 = false
	start("" +
		"setkeycodes 38 56 &" + // restore alt
		"setkeycodes 0e 14 &" + // restore bksp
		"setkeycodes 1e 30 &" + // restore a
		"setkeycodes 11 17 &" + // restore w
		"setkeycodes 20 32 &" + // restore d
		"setkeycodes 1f 31 &" + // restore s
		"")
}

func capsMode2On() {
	capsMode = false
	capsMode2 = true
	start("" +
		"setkeycodes 38 " + strconv.Itoa(KEY_1) + " &" + // alt
		"setkeycodes 0e 111 &" + // bksp -> del
		"setkeycodes 1e 102 &" + // a -> home
		"setkeycodes 11 104 &" + // w -> pgup
		"setkeycodes 20 107 &" + // d -> end
		"setkeycodes 1f 109 &" + // s -> pgdn
		"")
}

func send(keys string) {
	start("xdotool key " + keys)
}

func start(cmd string) {
	if _, err := exec.Command("/bin/sh", "-c", cmd).Output(); err != nil {
		log.Printf("Error during exec command: %s", err.Error())
	}
}
