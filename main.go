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
	KEY_2
	KEY_3
	KEY_4
	KEY_5
)

func main() {
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


	var capsMode bool
	var altPressed bool
	var pressedButtonsCount int
	var capsPressed bool

	for {
		select {
		case i := <-in:
			if i.Type != keylogger.EV_KEY {
				continue
			}
			//fmt.Printf("key:\t%s,\tcode:\t%d,\tevent:\t%d\n", i.KeyString(), i.Code, i.Value)
			switch uint(i.Code) {
			case KEY_0:// CAPS LOCK
				if int(i.Value) == 1 {
					capsPressed = true
					if !capsMode { // Press
						capsMode = true
						capsModeOn()
						if altPressed {
							//capsMode2On()
						}
					}
				}
				if int(i.Value) == 0 {
					capsPressed = false
					if pressedButtonsCount == 0 {
						capsMode = false
						capsModeOff()
					}
				}
			case KEY_3:
				if int(i.Value) == 1 {
					//altPressed = true
					capsMode2On()
				}
				if int(i.Value) == 0 {
					//altPressed = true
					if capsMode {
						capsMode2Off()
					}
				}
				//if uint(i.Code) == KEY_4 && int(i.Value) ==
				//if uint(i.Code) == 555 && int(i.Value) == 1 {
				//	send("ctrl")
				//}
				//if uint(i.Code) == KEY_3 && int(i.Value) == 1 {
				//	send("ctrl+b")
				//}
			case 56:
				if int(i.Value) == 1 {
					altPressed = true
				}
				if int(i.Value) == 0 {
					altPressed = false
				}
			case 103, 105, 106, 108:
				if int(i.Value) == 1 {
					pressedButtonsCount++
				}
				if int(i.Value) == 0 {
					pressedButtonsCount--
					if pressedButtonsCount == 0 && capsMode && !capsPressed {
						capsMode = false
						capsModeOff()
					}
				}
			}

		case <-quit:
			return
		}
	}
	//_=capsPressed
}

func disableKeys() {
	start("" +
		"setkeycodes 3a " + strconv.Itoa(KEY_0) + " &", // capslock
	//"setkeycodes e02ae037 "+strconv.Itoa(KEY_4)+" &"+ // prtsc

	//"setkeycodes e047 0 &"+ // home
	//"setkeycodes e04f 0 &"+ // end
	//"setkeycodes e049 0 &"+ // pgup
	//"setkeycodes e051 0 &"+ // pgdn
	//"setkeycodes e050 0 &"+ // down
	//"setkeycodes e048 0 &"+ // up
	//"setkeycodes e04b 0 &"+ // left
	//"setkeycodes e04d 0 &"+ // right

	//"setkeycodes e053 0 &", // delete
	)
}

func restoreKeys() {
	start("" +
		"setkeycodes 3a 58 &", // capslock
	//"setkeycodes e02ae037 "+strconv.Itoa(KEY_4)+" &"+ // prtsc

	//"setkeycodes e047 0 &"+ // home
	//"setkeycodes e04f 0 &"+ // end
	//"setkeycodes e049 0 &"+ // pgup
	//"setkeycodes e051 0 &"+ // pgdn
	//"setkeycodes e050 0 &"+ // down
	//"setkeycodes e048 0 &"+ // up
	//"setkeycodes e04b 0 &"+ // left
	//"setkeycodes e04d 0 &"+ // right

	//"setkeycodes e053 111 &", // delete
	)
}

func capsModeOn() {
	start("" +
		//"setkeycodes 15 " + strconv.Itoa(KEY_1) + " &" + // y
		//"setkeycodes 19 " + strconv.Itoa(KEY_2) + " &" + // p
		//"setkeycodes 33 99 &" + // , -> menu

		"setkeycodes 38 " + strconv.Itoa(KEY_3) + " &" + // alt
		"setkeycodes 0e 111 &" + // bksp -> del

		//"setkeycodes 39 "+strconv.Itoa(KEY_3)+" &"+ /KEY_1/ ,

		//"setkeycodes 16 104 &"+ // u -> pgup
		//"setkeycodes 17 103 &"+ // i -> up
		//"setkeycodes 18 109 &"+ // o -> pgdn
		//"setkeycodes 23 102 &"+ // h -> home
		//"setkeycodes 24 105 &"+ // j -> left
		//"setkeycodes 25 108 &"+ // k -> down
		//"setkeycodes 26 106 &"+ // l -> right
		//"setkeycodes 27 107 &"+ // ; -> end

		"setkeycodes 1e 105 &" + // a -> left
		"setkeycodes 11 103 &" + // w -> up
		"setkeycodes 20 106 &" + // d -> right
		"setkeycodes 1f 108 &" + // s -> down

		//"setkeycodes 20  29 &"+ // d -> ctrl
		//"setkeycodes 21  42 &"+ // f -> shift

		"")
}

func capsModeOff() {
	start("" +
		//"setkeycodes 15 21 &" + // y
		//"setkeycodes 19 25 &" + // p
		//"setkeycodes 33 51 &" + // ,

		"setkeycodes 38 56 &" + // alt
		"setkeycodes 0e 14 &" + // bksp

		//"setkeycodes 39 57 &"+ // ,

		//"setkeycodes 16 22 &"+ // u
		//"setkeycodes 17 23 &"+ // i
		//"setkeycodes 18 24 &"+ // o
		//"setkeycodes 23 35 &"+ // h
		//"setkeycodes 24 36 &"+ // j
		//"setkeycodes 25 37 &"+ // k
		//"setkeycodes 26 38 &"+ // l
		//"setkeycodes 27 39 &"+ // ;

		"setkeycodes 1e 30 &" + // a -> left
		"setkeycodes 11 17 &" + // w -> up
		"setkeycodes 20 32 &" + // d -> right
		"setkeycodes 1f 31 &" + // s -> down

		//"setkeycodes 20 32 &"+ // d
		//"setkeycodes 21 33 &"+ // f

		"")
}

func capsMode2On() {
	start("" +
		"setkeycodes 1e 102 &" + // a -> home
		"setkeycodes 11 104 &" + // w -> pgup
		"setkeycodes 20 107 &" + // d -> end
		"setkeycodes 1f 109 &" + // s -> pgdn
		"",
	)
}

func capsMode2Off() {
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
