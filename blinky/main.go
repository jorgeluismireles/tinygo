package main

import (
	"machine"
	"time"
)

func main() {
	led := machine.LED
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	for {
		led.Low()
		time.Sleep(time.Millisecond * 500)
		led.High()
		time.Sleep(time.Millisecond * 500)
	}
}

// Plug Rasberry pico into computer's USB port while holding down the board's reset button.
// Once plugged realese the reset button.
// Flash:
// $ tinygo flash -target=pico 