package buttons

import (
    "log"

    "periph.io/x/conn/v3/gpio"
    "periph.io/x/conn/v3/gpio/gpioreg"
    "periph.io/x/host/v3"
)

type Button struct {
	p gpio.PinIO
}

func NewButton(gpioPin string, handler func(chan bool)) (*Button, error) {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
        log.Fatal(err)
    }

    // Lookup a pin by its number:
    p := gpioreg.ByName(gpioPin)
    if p == nil {
        log.Println("Failed to find gpiopin")
    }

    // Set it as input, with an internal pull up resistor:
    if err := p.In(gpio.PullUp, gpio.BothEdges); err != nil {
        log.Println(err)
    }

    buttonChannel := make(chan bool)

    go handler(buttonChannel)

    // Wait for edges as detected by the hardware, and print the value read:
    for {
        p.WaitForEdge(-1)
		if !p.Read() {
            buttonChannel <- true
		}else {
            buttonChannel <- false
		}
    }
}