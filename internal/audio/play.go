package audio

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"os"
	"popaudio/internal/dlnaclient"
)

func PlayHttpPath(httpPath string) {
	httpReader, err := dlnaclient.HttpGetFile(httpPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	m, err := mp3.NewDecoder(httpReader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	channels := uint32(2)
	sampleRate := uint32(m.SampleRate())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	stateChan := make(chan State)

	ctrlChan := make(chan struct{})
	go Playback(m, channels, sampleRate, stateChan, ctrlChan)
	for {
		select {
		case state := <-stateChan:
			if state.State == Starting {
				fmt.Println("Starting")
			} else if state.State == Ready {
				fmt.Println("Playback ready")
			} else if state.State == Stopped {
				fmt.Println("Playback stopped")
			} else if state.State == Error {
				fmt.Println("Playback error", state.Err)
				return
			}
		}
	}
}
