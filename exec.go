package keywatcher

import (
	"strconv"
	"os/exec"
	"fmt"
)

func (a *app) enable() {
	a.start("" +
		"setkeycodes 3a " + strconv.Itoa(KEY_0) + " &" + // capslock
		"wait")
}

func (a *app) disable() {
	a.start("" +
		"setkeycodes 3a 58 &" + // capslock
		"wait")
}

func (a *app) capsModeOn() {
	a.state.capsMode = true
	a.state.capsMode2 = false
	a.start("" +
		"setkeycodes 38 " + strconv.Itoa(KEY_1) + " &" + // alt
		"setkeycodes 0e 111 &" + // bksp -> del
		"setkeycodes 1e 105 &" + // a -> left
		"setkeycodes 11 103 &" + // w -> up
		"setkeycodes 20 106 &" + // d -> right
		"setkeycodes 1f 108 &" + // s -> down
		"wait")
}

// disable capsMode and capsMode2
func (a *app) capsModeOff() {
	a.state.capsMode = false
	a.state.capsMode2 = false
	a.start("" +
		"setkeycodes 38 56 &" + // restore alt
		"setkeycodes 0e 14 &" + // restore bksp
		"setkeycodes 1e 30 &" + // restore a
		"setkeycodes 11 17 &" + // restore w
		"setkeycodes 20 32 &" + // restore d
		"setkeycodes 1f 31 &" + // restore s
		"wait")
}

func (a *app) capsMode2On() {
	a.state.capsMode = false
	a.state.capsMode2 = true
	a.start("" +
		"setkeycodes 38 " + strconv.Itoa(KEY_1) + " &" + // alt
		"setkeycodes 0e 111 &" + // bksp -> del
		"setkeycodes 1e 102 &" + // a -> home
		"setkeycodes 11 104 &" + // w -> pgup
		"setkeycodes 20 107 &" + // d -> end
		"setkeycodes 1f 109 &" + // s -> pgdn
		"wait")
}

func (a *app) send(keys string) {
	a.start("xdotool key " + keys)
}

func (a *app) start(cmd string) {
	if _, err := exec.Command("/bin/sh", "-c", cmd).Output(); err != nil {
		fmt.Printf("Error during exec command: %s", err.Error())
	}
}
