// package main

// import (
// 	"fmt"
// 	"io/ioutil"

// 	"math/rand"
//     "time"

// 	"github.com/coral/fluidsynth2"
// )

// const audio_driver string = "jack"
// const midi_driver string = "jack"
// const soundfont string = "Super_Italo_DiscoFont_Director_s_Cut.sf2"
// const RAND_MAX int = 126


// func main() {

// 	s := fluidsynth2.NewSettings()
// 	fmt.Println("\nAvaliable audio drivers:")
// 	for _, value := range s.GetOptions("audio.driver") {
// 		fmt.Println(value)
// 	}

// 	fmt.Println("\nAvaliable midi drivers:")
// 	for _, value := range s.GetOptions("midi.driver") {
// 		fmt.Println(value)
// 	}

// 	// Easy way to set audio backend
// 	s.SetString("audio.driver", audio_driver)
// 	s.SetString("midi.driver", midi_driver)

// 	s.SetInt("audio.jack.autoconnect", 1)
// 	s.SetInt("midi.autoconnect", 1)

// 	s.SetNum("synth.gain", 0.20)

// 	synth := fluidsynth2.NewSynth(s)

// 	sf_id := synth.SFLoad(soundfont, true)
// 	fmt.Printf("Soundfont id : %d\n", sf_id)

// 	player := fluidsynth2.NewPlayer(synth)
	
// 	// player.Add("Super Mario 64 - Medley.mid")

// 	// Example of how to play from memory
// 	dat, err := ioutil.ReadFile("Super Mario 64 - Medley.mid")
// 	if err != nil {
// 		panic(err)
// 	}

// 	player.AddMem(dat)

// 	player.SetBPM(300)
// 	player.SetTempo(300)

// 	fluidsynth2.NewAudioDriver(s, synth)

// 	// player.Play()
// 	// player.Join()
	
// 	for i := 0; i < 127; i++ {
// 		synth.ProgramChange(0, uint8(i))
// 		name := synth.SFGetPresetName(0)
// 		fmt.Printf("Program change to ")
// 		fmt.Println(name)

// 		for j := 0; j < 10; j++ {
// 			/* Generate a random key */
// 			rand.Seed(time.Now().UnixNano())
// 			min := 30
// 			max := 100
// 			// fmt.Println(rand.Intn(max - min + 1) + min)
// 			note := rand.Intn(max - min + 1) + min
// 			/* Play a note */
// 			synth.NoteOn(0, uint8(note), 80)
// 			/* Sleep for 1 second */
// 			time.Sleep(175 * time.Millisecond);
// 			/* Stop the note */
// 			synth.NoteOff(0, uint8(note))
// 		}
// 	}
// }