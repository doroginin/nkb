package nkb

import (
	"strconv"
	"os/exec"
	"fmt"
	"strings"
)

func (a *app) enable() {
	a.start("" +
		"setkeycodes 3a " + strconv.Itoa(KEY_0) + " &" + // capslock
		"setkeycodes e049 97 &" + // pgup -> ctrl
		"setkeycodes e01d 127 &" + // ctrl -> menu
		"wait")
}

func (a *app) disable() {
	a.start( "" +
		"setkeycodes 3a 58 &" + // capslock
		"sudo setkeycodes e049 104 &" + // ctrl -> pgup
		"setkeycodes e01d 97 &" + // ctrl -> menu
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

var user string // user who is logged in

func init() {
	var err error
	if user, err = sh(`who | grep '/dev/tty' | grep -oP '^.*?(?=\s)'`); err != nil {
		fmt.Printf("Error during get logged in user: %s\n", err.Error())
	} else {
		user = strings.Trim(user, "\n")
		fmt.Printf("Logged in user: %s\n", user)
	}
}

func (a *app) send(keys string) {
	a.start(`su ` + user + ` -c "export DISPLAY=':0.0'; xdotool key ` + keys + `"`)
}

func (a *app) start(cmd string) {
	if out, err := sh(cmd); err != nil {
		fmt.Printf("Error during exec command: %s, out: %s\n", err.Error(), out)
	}
}

func sh(cmd string) (string, error) {
	if out, err := exec.Command("/bin/sh", "-c", cmd).Output(); err == nil {
		return string(out), nil
	} else {
		return "", nil
	}
}
