package main


/* 

Raspberry pico power supply
 
            5.1V                    5.0V                       3.2V
           VBUS                    VSYS                        3V3
    J1     [40]                    [39]                        [36]
  +---+      |                       |       +------------+     |
  |  [1]-----+----------+----[D1]----+-------|5,8     1,10|-----+
  |  [2]--USB-DM (TP2)  |            |       |            |
  |  [3]--USB-DP (TP3)  R10          R2      |     U2     |
  |  [4]                |            |       |            |
  |  [5]--GND           +-- GPIO24   +-------+6     3,9,11|--+
  +---+                 |            |       +------------+  |
                        R1          [37]                    GND
                        |          3V3_EN
                        GND         4.9V


Raspberry pico J0 connector (40 pins)
                                          +-----| USB |-----+ 
UART-0 TX | I2C-0 SDA | SPI-0 RX  | GP0  [ 1]      J1     [40] VBUS
UART-0 RX | I2C-0 SCL | SPI-0 CSn | GP1  [ 2]             [39] VSYS
                                    GND  [ 3]             [38] GND
            I2C-1 SDA | SPI-0 SCK | GP2  [ 4]             [37] 3V3_EN
            I2C-2 SCL | SPI-0 TX  | GP3  [ 5]             [36] 3V3_OUT
UART-1 TX | I2C-0 SDA | SPI-0 RX  | GP4  [ 6]             [35]         | ADC-VREF
UART-2 RX | I2C-0 SCL | SPI-0 CSn | GP5  [ 7]             [34] GP28    | ADC-2
                                    GND  [ 8]             [33] GND     | AGND
            I2C-1 SDA | SPI-0 SCK | GP6  [ 9]             [32] GP27    | ADC-1     | I2C-1 SCL
            I2C-1 SCL | SPI-0 TX  | GP7  [10]             [31] GP26    | ADC-0     | I2C-1 SDA
UART-1 TX | I2C-0 SDA | SPI-1 RX  | GP8  [11]             [30] RUN
UART-2 RX | I2C-0 SCL | SPI-1 CSn | GP9  [12]             [29] GP22
                                  | GND  [13]             [28] GND
            I2C-1 SDA | SPI-1 SCK | GP10 [14]             [27] GP21    |           | I2C-0 SCL
            I2C-1 SCL | SPI-1 TX  | GP11 [15]             [26] GP20    |           | I2C-0 SDA
UART-0 TX | I2C-0 SDA | SPI-1 RX  | GP12 [16]             [25] GP19    | SPI-0 TX  | I2C-1 SCL
UART-0 RX | I2C-0 SCL | SPI-1 CSn | GP13 [17]             [24] GP18    | SPI-0 SCK | I2C-1 SDA
                                    GND  [18]             [23] GND       
            I2C-1 SDA | SPI-1 SCK | GP14 [19]    debug    [22] GP17    | SPI-0 CSn | I2C-0 SCL | UART-1 RX
            I2C-1 SCL | SPI-1 TX  | GP15 [20]  [-][-][-]  [21] GP16    | SPI-0 RX  | I2C-0 SDA | UART-0 TX
                                                S  G  S
                                                W  N  W
                                                C  D  D
                                                L     I
                                                K     O


Optoisolated input circuit:
     
        +----------------------------------[20] GP15
        |       +-------------resistor-----[36] 
        |       |
emitter |       | collector
    +--[D]-----[+]--+
    |       ^       |
    |    window     | ITR8102 top view
    |       ^       |
    +--[+]-----[E]--+
 anode  |       | cathode 
        |       +----------|>---resistor--- GND
      | o                 LED
    --|   button 
      | o
        |
        +---------------------------------- 5.1V
*/

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

