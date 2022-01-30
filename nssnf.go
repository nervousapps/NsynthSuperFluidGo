package main

import (
	"fmt"
	"io/ioutil"
	"math"

	// "os"

	"math/rand"

	"github.com/coral/fluidsynth2"

	"log"
	"time"

	"buttons"
	"inputs"
	"screen"
)

var main_screen *screen.Screen
var err error
var displayChan chan bool
var ui []screen.Msg

const audio_driver string = "jack"
const midi_driver string = "jack"
const soundfont string = "Super_Italo_DiscoFont_Director_s_Cut.sf2"
const RAND_MAX int = 126

var synth fluidsynth2.Synth
var s fluidsynth2.Settings

const CENTER int = 0
const POT1 int = 1
const POT2 int = 2
const POT3 int = 3
const POT4 int = 4
const POT5 int = 5
const POT6 int = 6
const TOUCH int = 7
const BUTTON1 int = 8
const BUTTON2 int = 9
const BUTTON3 int = 10
const BUTTON4 int = 11


func main() {
	// os.Clearenv()

	main_screen, err = screen.NewScreen()
	if err!= nil {
		log.Fatal(err)
	}

	initFluidSynth()

	for i := 0; i<=BUTTON4 ; i++{
		message := screen.Msg{}
		message.SetMsg(0,0,"")
		ui = append(ui, message)
	}

	forgeUi(CENTER, 10, 30, "NsynthSuperHard")
	
	time.Sleep(1000* time.Millisecond)

	displayChan = make(chan bool)
	// go main_screen.DisplayGif("images/golang.gif", displayChan)

	go buttons.NewButton("GPIO5", handlerButton1)
	go buttons.NewButton("GPIO6", handlerButton2)
	go buttons.NewButton("GPIO26", handlerButton3)
	go buttons.NewButton("GPIO13", handlerButton4)

	go inputs.StartInputs(touchHandler, encodersHandler, potentiometersHandler)

	// for i := 0; i<10 ; i++ {
	// 	for j := 0; j<10 ; j++ {
	// 		forgeUi(TOUCH, int(i*4)+80, int(j*4), ".")
	// 	}
	// }

	// go main_screen.DisplayGif("images/sexy-weiner.gif", displayChan)

	// time.Sleep(1000* time.Millisecond)
	// displayChan <- true

	// go playMidi()

	for {
		// time.Sleep(1000* time.Millisecond)
		playMidi()
	}
}

func forgeUi(index int, x int, y int, msg string, v ...interface{}){
	message := screen.Msg{}
	msg = fmt.Sprintf(msg, v...)
	message.SetMsg(x, y, msg)
	ui[index] = message
	main_screen.DisplayMessage(ui)
}

func handlerButton1(buttonChannel chan bool) {
	for {
		select{
			case pushed := <- buttonChannel:
			if pushed{
				log.Printf("Pushed 1")
				forgeUi(BUTTON4, 0, 10, "B1")
			}else{
				log.Printf("Released 1")
				forgeUi(BUTTON4, 0, 10, "  ")
			}
		}
	}
}

func handlerButton2(buttonChannel chan bool) {
	for {
		select{
			case pushed := <- buttonChannel:
			if pushed{
				log.Printf("Pushed 3")
				forgeUi(BUTTON4, 10, 10, "B2")
			}else{
				log.Printf("Released 2")
				forgeUi(BUTTON4, 10, 10, "  ")
			}
		}
	}
}

func handlerButton3(buttonChannel chan bool) {
	for {
		select{
			case pushed := <- buttonChannel:
			if pushed{
				log.Printf("Pushed 3")
				forgeUi(BUTTON4, 20, 10, "B3")
			}else{
				log.Printf("Released 3")
				forgeUi(BUTTON4, 20, 10, "  ")
			}
		}
	}
}

func handlerButton4(buttonChannel chan bool) {
	for {
		select{
			case pushed := <- buttonChannel:
			if pushed{
				log.Printf("Pushed 4")
				forgeUi(BUTTON4, 30, 10, "B4")
			}else{
				log.Printf("Released 4")
				forgeUi(BUTTON4, 30, 10, "  ")
			}
		}
	}
}

func touchHandler(touchChannel chan []byte) {
	for data := range touchChannel {
		log.Printf("Handler touch: %v\n", data)
		// for i := 0; i<10 ; i++ {
		// 	for j := 0; j<10 ; j++ {
		// 		forgeUi(TOUCH, int(i*4)+80, int(j*4), ".")
		// 	}
		// }
		forgeUi(TOUCH, int(data[0]*4)+80, int(data[1]*4), "*")
	}
}

func encodersHandler(encsChannel chan []byte) {
	previous_encs := make([]byte, 4)
	for encs := range encsChannel {
		// Change program with encoder 0
		if previous_encs[0] != encs[0] && encs[0]/2 < 126 {
			previous_encs[0] = encs[0]
			log.Printf("Handler enc 0: %v\n", encs[0])
			forgeUi(CENTER, 0, 30, "%d", int(encs[0]))
		}
		// Change program with encoder 1
		if previous_encs[1] != encs[1] && encs[1]/2 < 126 {
			previous_encs[1] = encs[1]
			go synth.ProgramChange(0, uint8(encs[1]/2))
			name := synth.SFGetPresetName(0)
			forgeUi(CENTER, 0, 30, name)
			log.Printf("Program change to ")
			log.Println(name)
		}
		// Change program with encoder 2
		if previous_encs[2] != encs[2] && encs[2]/2 < 126 {
			previous_encs[2] = encs[2]
			log.Printf("Handler enc 2: %v\n", encs[2])
			forgeUi(CENTER, 0, 30, "%d", int(encs[2]))
		}
		// Change program with encoder 3
		if previous_encs[3] != encs[3] && encs[3]/2 < 126 {
			previous_encs[3] = encs[3]
			log.Printf("Handler enc 3: %v\n", encs[3])
			forgeUi(CENTER, 0, 30, "%d", int(encs[3]))
		}
	}
}

func potentiometersHandler(potsChannel chan []byte) {
	previous_pots := make([]byte, 6)
	for pots := range potsChannel {
		if previous_pots[0] != pots[0] && (previous_pots[0]+2 < pots[0] || previous_pots[0]+2 > pots[0]) {
			go s.SetNum("synth.gain", math.Round((float64(pots[0])/126.0)*100)/100)
			forgeUi(POT1, 0, 64, "%0.2f", math.Round((float64(pots[0])/126.0)*100)/100)
			previous_pots[0] = pots[0]
		}
		if previous_pots[1] != pots[1] && (previous_pots[1]+2 < pots[1] || previous_pots[1]+2 > pots[1]) {
			forgeUi(POT2, 20, 54, "%d", int(pots[1])/2)
			previous_pots[1] = pots[1]
		}
		if previous_pots[2] != pots[2] && (previous_pots[2]+2 < pots[2] || previous_pots[2]+2 > pots[2]) {
			forgeUi(POT3, 40, 64, "%d", int(pots[2])/2)
			previous_pots[2] = pots[2]
		}
		if previous_pots[3] != pots[3] && (previous_pots[3]+2 < pots[3] || previous_pots[3]+2 > pots[3]) {
			forgeUi(POT4, 60, 54, "%d", int(pots[3])/2)
			previous_pots[3] = pots[3]
		}
		if previous_pots[4] != pots[4] && (previous_pots[4]+2 < pots[4] || previous_pots[4]+2 > pots[4]) {
			forgeUi(POT5, 80, 64, "%d", int(pots[4])/2)
			previous_pots[4] = pots[4]
		}
		if previous_pots[5] != pots[5] && (previous_pots[5]+2 < pots[5] || previous_pots[5]+2 > pots[5]) {
			forgeUi(POT6, 100, 54, "%d", int(pots[5])/2)
			previous_pots[5] = pots[5]
		}
	}
}

func initFluidSynth() {

	s = fluidsynth2.NewSettings()
	fmt.Println("\nAvaliable audio drivers:")
	for _, value := range s.GetOptions("audio.driver") {
		fmt.Println(value)
	}

	fmt.Println("\nAvaliable midi drivers:")
	for _, value := range s.GetOptions("midi.driver") {
		fmt.Println(value)
	}

	// Easy way to set audio backend
	s.SetString("audio.driver", audio_driver)
	s.SetString("midi.driver", midi_driver)

	s.SetInt("audio.jack.autoconnect", 1)
	s.SetInt("midi.autoconnect", 1)

	s.SetNum("synth.gain", 0.20)

	synth = fluidsynth2.NewSynth(s)

	sf_id := synth.SFLoad(soundfont, true)
	fmt.Printf("Soundfont id : %d\n", sf_id)

	player := fluidsynth2.NewPlayer(synth)

	// presets_names := synth.SFGetPresetsName()
	
	// player.Add("Super Mario 64 - Medley.mid")

	// Example of how to play from memory
	dat, err := ioutil.ReadFile("Super Mario 64 - Medley.mid")
	if err != nil {
		panic(err)
	}

	player.AddMem(dat)

	player.SetBPM(300)
	player.SetTempo(300)

	fluidsynth2.NewAudioDriver(s, synth)

	// player.Play()
	// player.Join()
}


func playMidi() {
	// for {
		for j := 0; j < 10; j++ {
			/* Generate a random key */
			rand.Seed(time.Now().UnixNano())
			min := 30
			max := 100
			// fmt.Println(rand.Intn(max - min + 1) + min)
			note := rand.Intn(max - min + 1) + min
			/* Play a note */
			synth.NoteOn(0, uint8(note), 80)
			/* Sleep for 1 second */
			time.Sleep(175 * time.Millisecond);
			/* Stop the note */
			synth.NoteOff(0, uint8(note))
		}
	// }
}