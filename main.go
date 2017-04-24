package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/MarinX/keylogger"
)

func main() {
	var device = flag.Int("device", -1, "Listen device with id")
	var code = flag.Uint("keycode", 0, "Listen key code")
	var event = flag.Int("event", 1, "Listen key event")
	var cmd = flag.String("cmd", "", "Run command")

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

	for i := range in {
		if i.Type == keylogger.EV_KEY {
			if *cmd == "" {
				fmt.Printf("key:\t%s,\tcode:\t%d,\tevent:\t%d\n", i.KeyString(), i.Code, i.Value)
			} else if uint(i.Code) == *code && int(i.Value) == *event {
				if err := exec.Command(*cmd).Start(); err != nil {
					log.Printf("Error during exec command: %s", err.Error())
				}
			}
		}
	}
}
