package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/MarinX/keylogger"
)

func main() {
	var device = flag.Int("d", -1, "Listen device with id")
	flag.Parse()

	devices, err := keylogger.NewDevices()
	if err != nil {
		fmt.Println(err)
		return
	}

	if *device < 0 {
		flag.Usage()
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

	var capsMode bool
	var capsPressed bool
	for i := range in {
		if i.Type == keylogger.EV_KEY {
			//fmt.Printf("key:\t%s,\tcode:\t%d,\tevent:\t%d\n", i.KeyString(), i.Code, i.Value)
			if uint(i.Code) == 58 { // CAPS LOCK
				if int(i.Value) == 1 { // Press
					capsPressed = true
					if !capsMode { // Turn on caps mode if it was disabled
						capsMode = true
						capsModeOn()
					} else {
						capsMode = false
					}
				}
				if int(i.Value) == 0 { // Release
					capsPressed = false
					if !capsMode { // Turn off caps mode only if no key was pressed with caps
						capsModeOff()
					}
				}
			} else {
				// if any key press during caps lock pressing then it's caps modifier mode
				// and we need to switch caps mode off after caps will be released
				if capsMode && capsPressed {
					capsMode = false
				}
			}
		}
	}
}

func disableKeys() {
	if _, err := exec.Command("/bin/sh", "-c", ""+
		"setkeycodes e047 0 &"+ // home
		"setkeycodes e04f 0 &"+ // end
		"setkeycodes e049 0 &"+ // pgup
		"setkeycodes e051 0 &"+ // pgdn
		"setkeycodes e050 0 &"+ // down
		"setkeycodes e048 0 &"+ // up
		"setkeycodes e04b 0 &"+ // left
		"setkeycodes e04d 0 &"+ // right
		"setkeycodes e053 0 &", // delete
	).Output(); err != nil {
		log.Printf("Error during exec command: %s", err.Error())
	}
}

func xkbHook() {
	//disable caps here
	//enable menu on prtsc
}

func capsModeOn() {
	if _, err := exec.Command("/bin/sh", "-c", ""+
		"setkeycodes 16 104 &"+ // u -> pgup
		"setkeycodes 17 103 &"+ // i -> up
		"setkeycodes 18 109 &"+ // o -> pgdn
		"setkeycodes 23 102 &"+ // h -> home
		"setkeycodes 24 105 &"+ // j -> left
		"setkeycodes 25 108 &"+ // k -> down
		"setkeycodes 26 106 &"+ // l -> right
		"setkeycodes 27 107 &"+ // ; -> end
		"setkeycodes 21  42 &"+ // f -> shift
		"setkeycodes 20  29 &"+ // d -> ctrl
		"setkeycodes 0e 111 &", // bksp -> del
	).Output(); err != nil {
		log.Printf("Error during exec command: %s", err.Error())
	}
}

func capsModeOff() {
	if _, err := exec.Command("/bin/sh", "-c", ""+
		"setkeycodes 16 22 &"+ // u
		"setkeycodes 17 23 &"+ // i
		"setkeycodes 18 24 &"+ // o
		"setkeycodes 23 35 &"+ // h
		"setkeycodes 24 36 &"+ // j
		"setkeycodes 25 37 &"+ // k
		"setkeycodes 26 38 &"+ // l
		"setkeycodes 27 39 &"+ // ;
		"setkeycodes 21 33 &"+ // f
		"setkeycodes 20 32 &"+ // d
		"setkeycodes 0e 14 &", // bksp
	).Output(); err != nil {
		log.Printf("Error during exec command: %s", err.Error())
	}
}
