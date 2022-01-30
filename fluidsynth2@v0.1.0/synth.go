package fluidsynth2

// #cgo pkg-config: fluidsynth
// #include <fluidsynth.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)


type Synth struct {
	ptr *C.fluid_synth_t
}

func NewSynth(settings Settings) Synth {
	return Synth{C.new_fluid_synth(settings.ptr)}
}

func (s *Synth) Close() {
	C.delete_fluid_synth(s.ptr)
}

func (s *Synth) SFLoad(path string, resetPresets bool) int {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	creset := cbool(resetPresets)
	cfont_id, _ := C.fluid_synth_sfload(s.ptr, cpath, creset)
	return int(cfont_id)
}

type Preset struct {
	ptr *C.fluid_preset_t
}

type Soundfont struct {
	ptr *C.fluid_sfont_t
}

func (s *Synth) SFGetPresetName(channel int) string {
	preset := Preset{C.fluid_synth_get_channel_preset(s.ptr, C.int(channel))}
	preset_name := C.fluid_preset_get_name(preset.ptr)
	presetString := C.GoString(preset_name)
	// C.free(unsafe.Pointer(preset_name))
	return presetString
}

func (s *Synth) SFGetPresetsName(sf_id int) []string {
	presets_name := make([]string, 200)
	soundfont := Soundfont{C.fluid_synth_get_sfont_by_id(s.ptr, C.int(sf_id))}
	C.fluid_sfont_iteration_start(soundfont.ptr)
	preset := Preset{C.fluid_sfont_iteration_next(soundfont.ptr)}
	for preset.ptr != nil {
		preset_name := C.fluid_preset_get_name(preset.ptr)
		presetString := C.GoString(preset_name)
		presets_name = append(presets_name, presetString)
		preset = Preset{C.fluid_sfont_iteration_next(soundfont.ptr)}
	}
	return presets_name
}

func (s *Synth) NoteOn(channel, note, velocity uint8) {
	C.fluid_synth_noteon(s.ptr, C.int(channel), C.int(note), C.int(velocity))
}

func (s *Synth) NoteOff(channel, note uint8) {
	C.fluid_synth_noteoff(s.ptr, C.int(channel), C.int(note))
}

func (s *Synth) ProgramChange(channel, program uint8) {
	C.fluid_synth_program_change(s.ptr, C.int(channel), C.int(program))
}

/* EFFECTS */
func (s *Synth) EnableReverb(fx_group int, on int) int {
	return int(C.fluid_synth_reverb_on(s.ptr, C.int(fx_group), C.int(on)))
}

func (s *Synth) SetReverbRoomsize(fx_group int, roomsize float64){
	C.fluid_synth_set_reverb_group_roomsize(s.ptr, C.int(fx_group), C.double(roomsize))
}

func (s *Synth) SetReverbLevel(fx_group int, level float64){
	C.fluid_synth_set_reverb_group_level(s.ptr, C.int(fx_group), C.double(level))
}

func (s *Synth) SetReverbWidth(fx_group int, width float64){
	C.fluid_synth_set_reverb_group_width(s.ptr, C.int(fx_group), C.double(width))
}

/* WriteS16 synthesizes signed 16-bit samples. It will fill as much of the provided
slices as it can without overflowing 'left' or 'right'. For interleaved stereo, have both
'left' and 'right' share a backing array and use lstride = rstride = 2. ie:
    synth.WriteS16(samples, samples[1:], 2, 2)
*/
func (s *Synth) WriteS16(left, right []int16, lstride, rstride int) {
	nframes := (len(left) + lstride - 1) / lstride
	rframes := (len(right) + rstride - 1) / rstride
	if rframes < nframes {
		nframes = rframes
	}
	C.fluid_synth_write_s16(s.ptr, C.int(nframes), unsafe.Pointer(&left[0]), 0, C.int(lstride), unsafe.Pointer(&right[0]), 0, C.int(rstride))
}

func (s *Synth) WriteFloat(left, right []float32, lstride, rstride int) {
	nframes := (len(left) + lstride - 1) / lstride
	rframes := (len(right) + rstride - 1) / rstride
	if rframes < nframes {
		nframes = rframes
	}
	C.fluid_synth_write_float(s.ptr, C.int(nframes), unsafe.Pointer(&left[0]), 0, C.int(lstride), unsafe.Pointer(&right[0]), 0, C.int(rstride))
}

type TuningId struct {
	Bank, Program uint8
}

/* ActivateKeyTuning creates/modifies a specific tuning bank/program */
func (s *Synth) ActivateKeyTuning(id TuningId, name string, tuning [128]float64, apply bool) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	C.fluid_synth_activate_key_tuning(s.ptr, C.int(id.Bank), C.int(id.Program), n, (*C.double)(&tuning[0]), cbool(apply))
}

/* ActivateTuning switches a midi channel onto the specified tuning bank/program */
func (s *Synth) ActivateTuning(channel uint8, id TuningId, apply bool) {
	C.fluid_synth_activate_tuning(s.ptr, C.int(channel), C.int(id.Bank), C.int(id.Program), cbool(apply))
}
