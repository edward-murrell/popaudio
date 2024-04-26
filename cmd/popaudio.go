package main

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"os"
	"popaudio/internal/audio"
	"popaudio/internal/dlnaclient"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(`"Call with path to server file-  eg http://server/file.mp3`)
		os.Exit(1)
	}
	httpPath := os.Args[1]

	playHttpPath(httpPath)
}

func playHttpPath(httpPath string) {
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
	stateChan := make(chan audio.State)

	ctrlChan := make(chan struct{})
	go audio.Playback(m, channels, sampleRate, stateChan, ctrlChan)
	for {
		select {
		case state := <-stateChan:
			if state.State == audio.Starting {
				fmt.Println("Starting")
			} else if state.State == audio.Ready {
				fmt.Println("Playback ready")
			} else if state.State == audio.Stopped {
				fmt.Println("Playback stopped")
			} else if state.State == audio.Error {
				fmt.Println("Playback error", state.Err)
				return
			}
		}
	}
}
