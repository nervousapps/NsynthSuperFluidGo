package main

import (
	"fmt"
	"io/ioutil"
	"math"

	// "os"

	"math/rand"

	// "github.com/coral/fluidsynth2"
	"fluidsynth2"

	// "github.com/xthexder/go-jack"
	"github.com/rakyll/launchpad"

	"log"
	"os"
	"time"

	"os/signal"

	"buttons"
	"inputs"
	"screen"
)

var main_screen *screen.Screen
var err error
var displayChan chan bool
var ui []screen.Msg
var ui_loading []screen.Msg
var ui_bar []screen.Msg
var available_presets []fluidsynth2.PresetName
var available_soundfonts []string
var currentVoices []fluidsynth2.Voice
var sf_id int
var programIndex int
var menuIndex int
var sfIndex int
var genTypeIndex fluidsynth2.FluidGenType

const audio_driver string = "jack"
const midi_driver string = "alsa_seq"
const default_soundfont string = "Super_Italo_DiscoFont_Director_s_Cut.sf2"
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
const MENU_LINE1 int = 12
const MENU_LINE2 int = 13
const MENU_LINE3 int = 14

func resetUi() []screen.Msg {
	var reseted_ui []screen.Msg
	for i := 0; i <= MENU_LINE3; i++ {
		message := screen.Msg{}
		message.SetMsg(0, 0, false, "")
		reseted_ui = append(reseted_ui, message)
	}
	return reseted_ui
}

func main() {
	// os.Clearenv()

	main_screen, err = screen.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	ui = resetUi()
	ui_loading = resetUi()
	ui_bar = resetUi()

	go forgeUiLoading(CENTER, 10, 30, "NsynthSuperHard")

	initFluidSynth()

	// go initLaunchpad()
	
	menuIndex = 1
	genTypeIndex = 0

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

func forgeUi(index int, x int, y int, invert bool, msg string, v ...interface{}) {
	message := screen.Msg{}
	msg = fmt.Sprintf(msg, v...)
	message.SetMsg(x, y, invert, msg)
	ui[index] = message
	main_screen.DisplayMessage(ui)
}

func forgeUiLoading(index int, x int, y int, msg string, v ...interface{}) {
	message := screen.Msg{}
	msg = fmt.Sprintf(msg, v...)
	message.SetMsg(x, y, true, msg)
	ui_loading[index] = message
	main_screen.DisplayLoadingMessage(ui_loading, false)
}

func forgeUiBarX(index int, x int, y int, percent int, msg string, v ...interface{}) {
	message := screen.Msg{}
	msg = fmt.Sprintf(msg, v...)
	message.SetMsg(x, y, true, msg)
	ui_bar[index] = message
	main_screen.DisplayBarXMessage(ui_bar, percent)
}

func handlerButton1(buttonChannel chan bool) {
	for {
		select {
		case pushed := <-buttonChannel:
			if pushed {
				log.Printf("Pushed 1")
				// forgeUi(BUTTON1, 0, 10, true, "B1")
			} else {
				log.Printf("Released 1")
				if menuIndex == 0 {
					sfIndex = sfIndex -1
					if sfIndex < 0 {
						sfIndex = len(available_soundfonts)-1
					}
				}else if menuIndex == 1 {
					sfIndex = sfIndex
				}else if menuIndex == 2 {
					sfIndex = sfIndex +1
					if sfIndex > len(available_soundfonts)-1 {
						sfIndex = 0
					}
				}
				selectSF()
				// forgeUi(BUTTON1, 0, 10, true, "  ")
			}
		}
	}
}

func handlerButton2(buttonChannel chan bool) {
	for {
		select {
		case pushed := <-buttonChannel:
			if pushed {
				log.Printf("Pushed 2")
				// forgeUi(BUTTON2, 10, 10, false, "B2")
			} else {
				log.Printf("Released 2")
				index := programIndex
				if menuIndex == 0 {
					index = index -1
					if index < 0 {
						index = len(available_presets)-1
					}
				}else if menuIndex == 1 {
					index = index
				}else if menuIndex == 2 {
					index = index +1
					if index > len(available_presets)-1 {
						index = 0
					}
				}
				synth.ProgramSelect(0,
					sf_id,
					available_presets[index].Bank,
					available_presets[index].Num)
				currentVoices = synth.GetVoicelist()
				log.Printf("Soundfont ID : %v", sf_id)
				log.Printf("Program change to %v", index)
				log.Println(available_presets[index].Name)
				log.Printf("Voices %v", currentVoices)
				// forgeUi(BUTTON2, 10, 10, false, "  ")
			}
		}
	}
}

func handlerButton3(buttonChannel chan bool) {
	for {
		select {
		case pushed := <-buttonChannel:
			if pushed {
				log.Printf("Pushed 3")
				forgeUi(BUTTON3, 20, 10, false, "B3")
			} else {
				log.Printf("Released 3")
				forgeUi(BUTTON3, 20, 10, false, "  ")
			}
		}
	}
}

func handlerButton4(buttonChannel chan bool) {
	for {
		select {
		case pushed := <-buttonChannel:
			if pushed {
				log.Printf("Pushed 4")
				forgeUi(BUTTON4, 30, 10, false, "B4")
			} else {
				log.Printf("Released 4")
				forgeUi(BUTTON4, 30, 10, false, "  ")
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
		forgeUi(TOUCH, int(data[0]*4)+80, int(data[1]*4), false, "*")
	}
}

func displayMenu(names []string, index int) {
	name_0 := ""
	name := names[index]
	name_1 := ""
	if index == 0 {
		name_0 = names[len(names)-1]
		name_1 = names[index+1]
	} else if index == len(names)-1 {
		name_0 = names[index-1]
		name_1 = names[0]
	}else {
		name_0 = names[index-1]
		name_1 = names[index+1]
	}
	if menuIndex == 0 {
		forgeUi(MENU_LINE1, 0, 15, true, "%v", name_0)
	}else {
		forgeUi(MENU_LINE1, 0, 15, false, "%v", name_0)
	}
	if menuIndex == 1 {
		forgeUi(MENU_LINE2, 0, 30, true, "%v", name)
	}else {
		forgeUi(MENU_LINE2, 0, 30, false, "%v", name)
	}
	if menuIndex == 2 {
		forgeUi(MENU_LINE3, 0, 45, true, "%v", name_1)
	}else {
		forgeUi(MENU_LINE3, 0, 45, false, "%v", name_1)
	}
}

func selectSF() {
	sf_id = sfIndex + 1
	available_presets = synth.SFGetPresetsName(sf_id)
	if programIndex > len(available_presets)-1 {
		programIndex = len(available_presets)-1
	}
	log.Printf("SFINDEX : %v\n", sfIndex)
	log.Printf("Soundfont : %v\n", available_soundfonts[sfIndex])
	displayMenu(fluidsynth2.GetPresetsName(available_presets), programIndex)
}

func encodersHandler(encsChannel chan []byte) {
	previous_encs := make([]byte, 4)
	for encs := range encsChannel {
		// Change program with encoder 0
		if previous_encs[0] != encs[0] {
			if encs[0] > previous_encs[0] {
				if menuIndex < 2 {
					menuIndex++
				}else {
					sfIndex++
				}
			} else {
				if menuIndex != 0 {
					menuIndex--
				}else {
					sfIndex--
				}
			}
			if sfIndex > len(available_soundfonts)-1 {
				sfIndex = 0
			}
			if sfIndex < 0 {
				sfIndex = len(available_soundfonts) - 1
			}
			displayMenu(available_soundfonts, sfIndex)
			previous_encs[0] = encs[0]
		}

		// Change program with encoder 1
		if previous_encs[1] != encs[1] {
			if encs[1] > previous_encs[1] {
				if menuIndex < 2 {
					menuIndex++
				}else {
					programIndex++
				}
			} else {
				if menuIndex != 0 {
					menuIndex--
				}else {
					programIndex--
				}
			}
			if programIndex > len(available_presets)-1 {
				programIndex = 0
			}
			if programIndex < 0 {
				programIndex = len(available_presets) - 1
			}
			displayMenu(fluidsynth2.GetPresetsName(available_presets), programIndex)
			previous_encs[1] = encs[1]
		}

		// Change program with encoder 2
		if previous_encs[2] != encs[2] && encs[2]/2 < 126 {
			go synth.SetPoly(int(encs[2]))
			log.Println("Polyphony : ", encs[3])
			previous_encs[2] = encs[2]
		}

		// Change program with encoder 3
		if previous_encs[3] != encs[3] && encs[3]/2 < 126 {
			if encs[3] > previous_encs[3] {
				genTypeIndex++
			} else {
				genTypeIndex--
			}
			if int(genTypeIndex) > len(fluidsynth2.FluidGenTypeMappings)-1 {
				genTypeIndex = 0
			}
			if int(genTypeIndex) < 0 {
				genTypeIndex = fluidsynth2.FluidGenType(len(fluidsynth2.FluidGenTypeMappings) - 1)
			}
			forgeUi(POT1, 0, 60, false, fluidsynth2.FluidGenTypeMappings[genTypeIndex])
			previous_encs[3] = encs[3]
			// forgeUi(CENTER, 0, 30, false, "%d", int(encs[3]))
		}
	}
}

func potentiometersHandler(potsChannel chan []byte) {
	previous_pots := make([]byte, 6)
	for pots := range potsChannel {
		if previous_pots[0] != pots[0] && (previous_pots[0]+2 < pots[0] || previous_pots[0]+2 > pots[0]) {
			go synth.SetGain(math.Round((float64(pots[0])/126.0)*100) / 100)
			// go s.SetNum("synth.gain", math.Round((float64(pots[0])/126.0)*100)/100)
			forgeUiBarX(POT1, 0, 64, int(pots[0]/2), "Volume")
			previous_pots[0] = pots[0]
		}
		if previous_pots[1] != pots[1] && (previous_pots[1]+2 < pots[1] || previous_pots[1]+2 > pots[1]) {
			// forgeUi(POT2, 20, 54, false, "%0.2f", float64(pots[1])/255.0)
			// if pots[1] == 0 {
			// 	// result := synth.EnableReverb(-1, 0)
			// 	go s.SetInt("synth.reverb.active", 0)
			// 	// forgeUi(POT2, 20, 54, "%d", result)
			// }else{
			// 	// result := synth.EnableReverb(-1, 1)
			// 	// forgeUi(POT2, 20, 54, "%d", result)
			// 	go s.SetInt("synth.reverb.active", 1)
			// }
			go synth.SetReverbLevel(1, float64(pots[2])/10.0)
			// currentVoices[0].VoiceReverbSendSet(float64(pots[5])/255.0)
			// for _, voice := range currentVoices {
			// 	voice.VoiceReverbSendSet((float64(pots[5]) / 255.0) * 100)
			// }
			// synth.MidiReverbSendSet(float64(pots[5])/255.0)
			forgeUiBarX(POT1, 0, 64, int(pots[1]/2), "ReverbLevel")
			previous_pots[1] = pots[1]
		}
		if previous_pots[2] != pots[2] && (previous_pots[2]+2 < pots[2] || previous_pots[2]+2 > pots[2]) {
			// forgeUi(POT3, 40, 64, false, "%0.2f", float64(pots[2])/255.0)
			previous_pots[2] = pots[2]
			// go synth.SetReverbLevel(1, float64(pots[2])/10.0)
			// go s.SetNum("synth.reverb.level", float64(pots[2])/255.0)
			forgeUiBarX(POT1, 0, 64, int(pots[2]/2), "Filter Cutoff")
			go synth.MidiFilterCutoffSet(float64(pots[2]/2))
		}
		if previous_pots[3] != pots[3] && (previous_pots[3]+2 < pots[3] || previous_pots[3]+2 > pots[3]) {
			// forgeUi(POT4, 60, 54, false, "%0.2f", float64(pots[3])/10.0)
			previous_pots[3] = pots[3]
			forgeUiBarX(POT1, 0, 64, int(pots[3]/2), "Filter Res")
			go synth.MidiFilterResoSet(float64(pots[3]/2))
			// go synth.SetReverbWidth(1, float64(pots[3])/10.0)
			// go s.SetNum("synth.reverb.width", float64(pots[3])/255.0)
		}
		if previous_pots[4] != pots[4] && (previous_pots[4]+2 < pots[4] || previous_pots[4]+2 > pots[4]) {
			// forgeUi(POT5, 80, 64, false, "%d", float64(pots[4])/255.0)
			go synth.SetReverbDamp(1, float64(pots[4])/255.0)
			// go s.SetNum("synth.reverb.damp", float64(pots[4])/255.0)
			previous_pots[4] = pots[4]
		}
		if previous_pots[5] != pots[5] && (previous_pots[5]+2 < pots[5] || previous_pots[5]+2 > pots[5]) {
			// forgeUi(POT6, 100, 54, false, "%d", int(pots[5])/2)
			go synth.MidiGenSet(genTypeIndex, float64(pots[5]/2))
			forgeUiBarX(POT1, 0, 64, int(pots[5]/2), fluidsynth2.FluidGenTypeMappings[genTypeIndex])
			previous_pots[5] = pots[5]
		}
	}
}

func load_sfont(soundfont string) {
	sf_id = synth.SFLoad("soundfonts/"+soundfont, true)
	fmt.Println("Soundfont id : ", sf_id)
	fmt.Println("Soundfont name : ", soundfont)
	presets := synth.SFGetPresetsName(sf_id)
	for _, preset_name := range presets {
		if preset_name.Name != "" {
			available_presets = append(available_presets, preset_name)
		}
	}
	fmt.Println("Available presets : ", len(available_presets))
}

func load_available_preset() {
	presets := synth.SFGetPresetsName(sf_id)
	for _, preset_name := range presets {
		if preset_name.Name != "" {
			available_presets = append(available_presets, preset_name)
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

	// Set settings
	s.SetString("audio.driver", audio_driver)
	s.SetInt("audio.jack.autoconnect", 1)
	s.SetString("midi.driver", midi_driver)
	s.SetInt("midi.autoconnect", 1)
	s.SetString("midi.portname", "FluidSynth")

	// Create the synth
	synth = fluidsynth2.NewSynth(s)

	// Set Gain
	s.SetNum("synth.gain", 0.20)

	// Set polyphony
	synth.SetPoly(200)

	// Create midi driver
	midi_driver := fluidsynth2.NewFluidMidiDriver(s, synth)
	midi_router := fluidsynth2.NewFluidMidiRouter(s, synth)

	// Create audio driver
	audio_driver := fluidsynth2.NewAudioDriver(s, synth)

	// Get all available soundfont files in the soundfonts folder
	files, err := ioutil.ReadDir("./soundfonts")
	if err != nil {
		log.Fatal(err)
	}

	for index, f := range files {
		fmt.Println("Load : ", int((100/len(files))*index))
		load_sfont(f.Name())
		available_soundfonts = append(available_soundfonts, f.Name())
		forgeUiBarX(POT1, 0, 60, int((100/len(files))*index+(100/len(files))), "Loading SF ...")
	}

	sf_id = 1
	sfIndex = sf_id - 1
	load_available_preset()
	programIndex = 0
	synth.ProgramSelect(0,
		sf_id,
		available_presets[programIndex].Bank,
		available_presets[programIndex].Num)

	log.Println("Effect channels : ", synth.CountEffectsChannels())
	log.Println("Effect groups : ", synth.CountEffectsGroups())

	// Example of how to play from memory
	// player := fluidsynth2.NewPlayer(synth)
	// dat, err := ioutil.ReadFile("Super Mario 64 - Medley.mid")
	// if err != nil {
	// 	panic(err)
	// }

	// player.AddMem(dat)

	// player.SetBPM(300)
	// player.SetTempo(300)

	// player.Play()
	// player.Join()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			audio_driver.Close()
			midi_driver.Close()
			midi_router.Close()
			log.Println(sig)
			os.Exit(0)
		}
	}()
}

func playMidi() {
	// for {
	for j := 0; j < 10; j++ {
		/* Generate a random key */
		rand.Seed(time.Now().UnixNano())
		min := 30
		max := 60
		// fmt.Println(rand.Intn(max - min + 1) + min)
		note := rand.Intn(max-min+1) + min
		/* Play a note */
		synth.NoteOn(0, uint8(note), 80)
		/* Sleep for 1 second */
		time.Sleep(175 * time.Millisecond)
		/* Stop the note */
		synth.NoteOff(0, uint8(note))
	}
	// }
}

func initLaunchpad() {
	pad, err := launchpad.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer pad.Close()

	pad.Clear()

	ch := pad.Listen()
	for {
		select {
		case hit := <-ch:
			pad.Light(hit.X, hit.Y, 3, 5)
		}
	}
}
