package ui

import (
	"fmt"
	"os"
	"popaudio/internal/audio"
	"popaudio/internal/dlnaclient"
	"popaudio/pkg/model"
	"strings"
)

type ScanLineBrowser struct {
	dlna *dlnaclient.DlnaClient
}

func NewScanLineBrowser(dlna *dlnaclient.DlnaClient) *ScanLineBrowser {
	return &ScanLineBrowser{
		dlna: dlna,
	}
}

func (b *ScanLineBrowser) Run() {
	root, err := b.dlna.GetRoot()
	if err != nil {
		println("error getting root", err)
		os.Exit(1) // Give up
	}

	b.Loop(root)
}

func (b *ScanLineBrowser) Loop(root *model.DIDLLite) {
	entity := root
	var input string
	for {
		b.PrintContainer(entity)
		_, err := fmt.Scanf("%s", &input)
		if err != nil {
			fmt.Printf("could not read input: %s", err.Error())
		}
		if input == "q" {
			os.Exit(1)
		}
		// horrible for loop over entity.Container to find if matches and browse that loop or play it
		for _, c := range entity.Containers {
			if c.ID == input {
				newEntity, err := b.dlna.Browse(&c)
				if err != nil {
					fmt.Printf("error getting sub container: %s", err.Error())
				} else {
					entity = newEntity
					break
				}
			}
		}

		for _, i := range entity.Items {
			if i.ID == input {
				audio.PlayHttpPath(i.Res.Text)
			}
		}

	}
}

func (b *ScanLineBrowser) PrintContainer(entity *model.DIDLLite) {
	println("Folders")
	for _, c := range entity.Containers {
		fmt.Printf("%s: %s %s (%s)\n", c.ID, c.Title, c.Class, c.ChildCount)
	}
	println("Music")
	for _, i := range entity.Items {
		if strings.Contains(i.Res.ProtocolInfo, "MP3") { // Temporarily filter out anything not mp3
			fmt.Printf("%s: %s by %s from %s\n", i.ID, i.Title, i.Artist, i.Album)
		}
	}
}
