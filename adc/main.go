package main

import (
	"machine"
	"time"
)

// Connections:
//                      100kohms
// [38] GND     --------^v^v^v^-,
//                         |    |
// [36] 3V3_OUT ----------------'
//                         |
// [34] GPIO28  -----------'

// Flash:
// $ cd ~/go/src/gitlab.com/jmireles/tinygo
// $ sudo tinygo flash -target=pico adc/main.go
func main() {
	machine.InitADC()

	sensor := machine.ADC{machine.ADC2} // GPIO28 pin-34
	sensor.Configure(machine.ADCConfig{})

	//simple(&sensor) // OK worked
	simple1(&sensor)
	//blink(&sensor)
}

func NewLED() *machine.Pin {
	led := machine.LED
	led.Configure(machine.PinConfig{
		Mode: machine.PinOutput,
	})
	return &led	
}

func ledChanger(ch<- chan bool) {
	led := NewLED()
	for {
		select {
		case state := <-ch:
			if state {
				led.High()
			} else {
				led.Low()
			}
		}
	}
}


// When the potentiometer is turned past the midway point,
// the build-in LED will light up.
func simple(sensor *machine.ADC) {
	led := NewLED()
	for {
		val := sensor.Get()
		if val < 0x8000 {
			led.Low()
		} else {
			led.High()
		}
		time.Sleep(time.Millisecond * 100)
	}
}



// Same as simple except led changes are channel's
func simple1(sensor *machine.ADC) {
	ch := make(chan bool, 1)
	go ledChanger(ch)
	for {
		val := sensor.Get()
		if val < 0x8000 {
			ch<- false
		} else {
			ch<- true
		}
		time.Sleep(time.Millisecond * 100)
	}
}



// time.NewTicker/time.NewTimer are unimplemeted yet 0.22!
// github.com/tinygo-org/tinygo issue #1037
// tinygo:ld.lld: error: undefined symbol: time.modTimer
func blink(sensor *machine.ADC) {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		high := false
		led := NewLED()
		for {
			select {
			case <-ticker.C:
				if high {
					led.High()
				} else {
					led.Low()
				}
				high = !high
			}
		}
	}()
	for {
		val := sensor.Get()
		if val < 0x8000 {
			ticker.Reset(time.Millisecond * 250)
		} else {
			ticker.Reset(time.Millisecond * 500)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

