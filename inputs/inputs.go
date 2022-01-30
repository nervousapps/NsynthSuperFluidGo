package inputs

import (
	"log"
	"time"
	// "os/user"
	"reflect"

	"github.com/go-daq/smbus"
)

func find(slice []byte, val byte) (bool) {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

type Inputs struct {
	touchChannel chan []byte
	encsChannel chan []byte
	potsChannel chan []byte
}

func (inputs *Inputs) readMcu() {
	buf := make([]byte, 16)
	touch := make([]byte, 2)
	encs := make([]byte, 4)
	pots := make([]byte, 6)
	previousBuf := make([]byte, len(buf))
	previousTouch := make([]byte, len(touch))
	previousEncs := make([]byte, len(encs))
	previousPots := make([]byte, len(pots))

	// usr, err := user.Current()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if usr.Name != "root" {
	// 	log.Fatal("need root access")
	// }

	c, err := smbus.Open(1, 0x47)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	for {
		err = c.ReadBlockData(0x47, 0, buf)
		if err != nil {
			log.Print(err)
			continue
		}
		// joined := bytes.Join(buf, []byte(" "))

		// Compute data
		copy(touch, buf[:2])
		copy(encs, buf[2:6])
		copy(pots, buf[6:12])
		// chk := buf[12]

		if reflect.DeepEqual(previousBuf, buf) {
			continue
		}
		copy(previousBuf, buf)

		if !reflect.DeepEqual(previousTouch, touch) && !find(touch, 255) {
			// go touch_handler(touch)
			inputs.touchChannel <- touch
			copy(previousTouch, touch)
		}

		if !reflect.DeepEqual(previousEncs, encs) {
			// go encoders_handler(encs)
			inputs.encsChannel <- encs
			copy(previousEncs, encs)
		}

		if !reflect.DeepEqual(previousPots, pots) {
			// go potentiometers_handler(pots)
			inputs.potsChannel <- pots
			copy(previousPots, pots)
		}

		time.Sleep(1*time.Millisecond)
	}
}

func StartInputs(touchHandler, encodersHandler, potentiometersHandler func(chan []byte)) (Inputs){
	log.Printf("Starting inputs !\n")
	inputs := &Inputs {
		touchChannel : make(chan []byte),
		encsChannel : make(chan []byte),
		potsChannel : make(chan []byte),
	}
	
	go touchHandler(inputs.touchChannel)
	go encodersHandler(inputs.encsChannel)
	go potentiometersHandler(inputs.potsChannel)

	go inputs.readMcu()

	return *inputs
}

func (inputs *Inputs) CloseInputs() {
	log.Printf("Stopping !\n")
	close(inputs.encsChannel)
	close(inputs.touchChannel)
	close(inputs.potsChannel)
}
