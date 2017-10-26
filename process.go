package keywatcher

import (
	"os/signal"
	"syscall"
	"github.com/MarinX/keylogger"
	"fmt"
	"time"
	"os"
)

const (
	KEY_0 = iota + 500
	KEY_1
)

const (
	EVENT_RELEASE = iota
	EVENT_PRESS
	EVENT_HOLD
)

func (a *app) Run() {
	//a.enable()
	//defer a.disable()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case i := <-a.ch:
			if i.Type != keylogger.EV_KEY {
				continue
			}
			if a.verbose {
				fmt.Printf("key: %s, code: %d, event: %d\n", i.KeyString(), i.Code, i.Value)
			}

			event := int(i.Value)
			code := uint(i.Code)

			switch code {
			case 58:
				a.processCaps0(event)
			case KEY_0: // CAPS LOCK
				a.processCaps(event)
			case KEY_1:
				a.processCaps2(event)
			case
				// capsmode and capsmode2 off
				14, // bksp
				56, // alt
				30, 17, 32, 31, // a, s, d, w

				// capsmode on
				103, 105, 106, 108, // arrows

				// capsmode or capsmode 2 on
				111, // delete

				// capsmode2 on
				102, 104, 107, 109: //home, end, pgup, pgdn

				a.processButtons(code, event)
			}
			if a.verbose {
				fmt.Printf("capsPressed: %t, caps2Pressed: %t, capsMode: %t, capsMode2: %t, pressedButtons: %v\n",
					a.state.capsPressed, a.state.caps2Pressed, a.state.capsMode, a.state.capsMode2, a.state.pressedButtons)
			}

		case <-quit:
			return
		}
	}
}

func (a *app) processCaps0(event int) {
	if event == EVENT_PRESS {
		if time.Since(a.state.lastCapsPressed) < 300*time.Millisecond {
			fmt.Println("double caps pressed => enabling")
			a.state.pressedButtons = make(map[uint]struct{})
			a.state.enabled = true
			a.enable()
		}
		a.state.lastCapsPressed = time.Now()
	}
}

func (a *app) processCaps(event int) {
	if event == EVENT_PRESS {
		if time.Since(a.state.lastCapsPressed) < 300*time.Millisecond {
			fmt.Println("double caps pressed => disabling")
			a.state.pressedButtons = make(map[uint]struct{})
			a.state.enabled = false
			a.disable()
		}
		a.state.lastCapsPressed = time.Now()

		a.state.capsPressed = true
		if len(a.state.pressedButtons) == 0 {
			a.capsModeOn()
		}
	}
	if event == EVENT_RELEASE {
		a.state.capsPressed = false
		if len(a.state.pressedButtons) == 0 {
			a.capsModeOff()
		}
	}
}

func (a *app) processCaps2(event int) {
	if event == EVENT_PRESS {
		a.state.caps2Pressed = true
		if len(a.state.pressedButtons) == 0 {
			a.capsMode2On()
		}
	}
	if event == EVENT_RELEASE {
		a.state.caps2Pressed = false
		if len(a.state.pressedButtons) == 0 {
			if a.state.capsPressed {
				a.capsModeOn()
			} else {
				a.capsModeOff()
			}
		}
	}
}

func (a *app) processButtons(code uint, event int) {
	if a.shouldIgnoreEvent(code, event) {
		return
	}
	if event == EVENT_PRESS {
		a.state.pressedButtons[code] = struct{}{}
	}
	if event == EVENT_RELEASE {
		delete(a.state.pressedButtons, code)
		if len(a.state.pressedButtons) == 0 {
			if a.state.capsMode2 && !a.state.caps2Pressed && !a.state.capsPressed ||
				a.state.capsMode && !a.state.capsPressed {
				a.capsModeOff()
			}
			if a.state.capsPressed {
				if a.state.caps2Pressed {
					a.capsMode2On()
				} else {
					a.capsModeOn()
				}
			}
		}
	}
}

func (a *app) shouldIgnoreEvent(code uint, event int) bool {
	switch code {
	// capsmode and capsmode2 off
	case 14, // bksp
		30, 17, 32, 31, // a, s, d, w
		56: // alt
		if (a.state.capsMode || a.state.capsMode2) && event == EVENT_PRESS {
			return true
		}
		// capsmode on
	case 103, 105, 106, 108: // arrows
		if !a.state.capsMode && event == EVENT_PRESS {
			return true
		}
		// capsmode or capsmode 2 on
	case 111: // delete
		if !a.state.capsMode && !a.state.capsMode2 && event == EVENT_PRESS {
			return true
		}
		// capsmode2 on
	case 102, 104, 107, 109: //home, end, pgup, pgdn
		if !a.state.capsMode2 && event == EVENT_PRESS {
			return true
		}
	}
	return false
}