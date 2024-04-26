package audio

import (
	"fmt"
	"github.com/gen2brain/malgo"
	"io"
)

type StateType uint32

const (
	Starting StateType = 1
	Ready    StateType = 2
	Playing  StateType = 3
	Stopped  StateType = 4
	Error    StateType = 255 // Implies closed.
)

type State struct {
	State StateType
	Err   error // Only when State is Error
}

// Playback this
func Playback(reader io.Reader, channels uint32, sampleRate uint32, status chan<- State, ctrl <-chan struct{}) {
	status <- State{Starting, nil}
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		fmt.Printf("ALSA LOG <%v>\n", message)
	})
	if err != nil {
		status <- State{Error, err}
		return
	}
	defer func() {
		_ = ctx.Uninit()
		ctx.Free()
	}()

	deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	deviceConfig.Playback.Format = malgo.FormatS16
	deviceConfig.Playback.Channels = channels
	deviceConfig.SampleRate = sampleRate
	deviceConfig.Alsa.NoMMap = 1

	// This is the function that's used for sending more data to the audio for playback.
	onSamples := func(pOutputSample, pInputSamples []byte, framecount uint32) {
		io.ReadFull(reader, pOutputSample)
	}

	deviceCallbacks := malgo.DeviceCallbacks{
		Data: onSamples,
	}
	device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
	if err != nil {
		status <- State{Error, err}
		return
	}
	defer device.Uninit()

	status <- State{Ready, nil}
	err = device.Start()
	status <- State{Playing, nil}
	// Currently it either returns an error or nothing until it's done.
	device.IsStarted()

	<-ctrl // Temporary hack until we build a proper stop method

	defer func() {
		status <- State{Stopped, nil}
	}()
}
