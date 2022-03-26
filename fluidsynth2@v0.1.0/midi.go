package fluidsynth2

// #cgo pkg-config: fluidsynth
// #include <fluidsynth.h>
// #include <stdlib.h>
import "C"
import "unsafe"
import "fmt"

type MidiDriver struct {
	ptr *C.fluid_midi_driver_t
}

func NewFluidMidiDriver(settings Settings, synth Synth) MidiDriver {
	return MidiDriver{C.new_fluid_midi_driver(settings.ptr, (*[0]byte)(C.fluid_synth_handle_midi_event), unsafe.Pointer(synth.ptr))}
}

func (m *MidiDriver) Close() {
	C.delete_fluid_midi_driver(m.ptr)
}	

type MidiRouter struct {
	ptr *C.fluid_midi_router_t
}

// type HandleMidiEvent struct {
// 	ptr *C.fluid_synth_handle_midi_event
// }

func NewFluidMidiRouter(settings Settings, synth Synth) MidiRouter {
	return MidiRouter{C.new_fluid_midi_router(settings.ptr, (*[0]byte)(C.fluid_synth_handle_midi_event), unsafe.Pointer(synth.ptr))}
}

func (m *MidiRouter) Close() {
	C.delete_fluid_midi_router(m.ptr)
}

type MidiEvent struct {
	ptr *C.fluid_midi_event_t
}

func handle_midi_event(data, event MidiEvent) int {
    fmt.Printf("event type: %d\n", C.fluid_midi_event_get_type(event.ptr))
	return int(C.fluid_midi_event_get_type(event.ptr))
}
