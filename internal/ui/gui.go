package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"os"
	"popaudio/internal/audio"
	"popaudio/internal/dlnaclient"
	"popaudio/pkg/model"
)

type FyneDisplay struct {
	dlna               *dlnaclient.DlnaClient
	containerList      []model.Container
	itemList           []model.Item
	containDisplayList *widget.List
	itemDisplayList    *widget.List
}

func NewFyneDisplay(dlna *dlnaclient.DlnaClient) *FyneDisplay {
	return &FyneDisplay{
		dlna: dlna,
	}
}

func (disp *FyneDisplay) Run() {
	root, err := disp.dlna.GetRoot()

	if err != nil {
		println("error getting root", err)
		os.Exit(1) // Give up
	}
	disp.containerList = root.Containers
	disp.itemList = root.Items

	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")
	myWindow.FullScreen()
	disp.containDisplayList = widget.NewList(
		func() int {
			return len(disp.containerList)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			x := disp.containerList
			o.(*widget.Button).SetText(x[i].Title)
			o.(*widget.Button).OnTapped = func() {
				disp.BrowseContainer(&x[i])
			}
		})
	disp.itemDisplayList = widget.NewList(
		func() int {
			return len(disp.itemList)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			items := disp.itemList
			o.(*widget.Button).SetText(items[i].Title)
			o.(*widget.Button).OnTapped = func() {
				println("playing", items[i].Res.ProtocolInfo)
				go audio.PlayHttpPath(items[i].Res.Text)
			}
		})

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), disp.containDisplayList, disp.itemDisplayList))
	myWindow.SetFullScreen(true)
	myWindow.ShowAndRun()
}

func (disp *FyneDisplay) BrowseContainer(con *model.Container) {
	fmt.Printf("requesting ID: %s", con.ID)
	o, err := disp.dlna.Browse(con)

	if err != nil {
		println("Error: ", err)
		return
	}
	println("done")
	disp.containerList = o.Containers
	disp.itemList = o.Items
	disp.containDisplayList.Refresh()
	disp.itemDisplayList.Refresh()
}
