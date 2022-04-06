package main

// See:
// github.com/tinygo-org/tinygo/blob/release/src/machine/board_pico.go
// github.com/tinygo-org/tinygo/blob/release/src/machine/machine_rp2040.go

import (
	"machine"
	"runtime/volatile"
	"time"
)

// This program turns on the pico led when button is pressed
// and turns off led when button is realased
// Plug Rasberry pico into computer's USB port while holding down the board's reset button.
// Once plugged realese the reset button.
func bruteForce(led machine.Pin) {
	btn := machine.GPIO15
	btn.Configure(machine.PinConfig{
		Mode: machine.PinInput,
	})
	for {
		time.Sleep(100 * time.Millisecond)
		if btn.Get() {
			led.Set(true)
		} else {
			led.Set(false)
		}
	}
}


// From tinygo/src/examples/pininterrupt
// Without debouncing!
//
//                     | mode machine.*    | change machine.*       | pico?
//                     |-------------------|------------------------|-----
// circuitplay_express | PinInputPulldown  | PinFalling             | fails flip/flop
// pca10040            | PinInputPullup    | PinRising              | does not work for btn to v+
// stm32               | PinInputPulldown  | PinRising | PinFalling | fails flip/flop
// wioterminal         | PinInput          | PinFalling             | ?
// other               | ?                 | ?                      |

func interrupt(led machine.Pin) {
	var state volatile.Register8
	state.Set(0)

	btn := machine.GPIO15
	mode := machine.PinInput
	change := machine.PinFalling
	btn.Configure(machine.PinConfig{Mode: mode})
	err := btn.SetInterrupt(change, func(pin machine.Pin) {
		if state.Get() != 0 {
			state.Set(0)
			led.Low()
		} else {
			state.Set(1)
			led.High()
		}
	})
	if err != nil {
		println("pin interrupt error:", err.Error())
	}
	for {
		time.Sleep(time.Hour)
	}
}

// Flash:
// $ cd ~/go/src/gitlab.com/jmireles/tinygo
// $ sudo tinygo flash -target=pico inputs1/main.go
func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	bruteForce(led) // worked OK
	//interrupt(led)
}

