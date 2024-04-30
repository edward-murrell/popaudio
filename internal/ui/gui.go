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
	dlna *dlnaclient.DlnaClient
}

func NewFyneDisplay(dlna *dlnaclient.DlnaClient) *FyneDisplay {
	return &FyneDisplay{
		dlna: dlna,
	}
}

func (b *FyneDisplay) Run() {
	var containerList *[]model.Container
	var itemList *[]model.Item
	root, err := b.dlna.GetRoot()

	if err != nil {
		println("error getting root", err)
		os.Exit(1) // Give up
	}
	containerList = &root.Containers
	itemList = &root.Items

	myApp := app.New()
	myWindow := myApp.NewWindow("List Widget")
	myWindow.FullScreen()
	listFolders := widget.NewList(
		func() int {
			return len(*containerList)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			x := *containerList
			o.(*widget.Button).SetText(x[i].Title)
			o.(*widget.Button).OnTapped = func() {
				b.BrowseContainer(&x[i], &containerList, &itemList)
			}
		})
	listItems := widget.NewList(
		func() int {
			return len(*itemList)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("template", func() {})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			items := *itemList
			o.(*widget.Button).SetText(items[i].Title)
			o.(*widget.Button).OnTapped = func() {
				audio.PlayHttpPath(items[i].Res.Text)
			}

		})

	myWindow.SetContent(container.New(layout.NewHBoxLayout(), listFolders, listItems))
	myWindow.ShowAndRun()
}

func (b *FyneDisplay) BrowseContainer(con *model.Container, conList **[]model.Container, itemList **[]model.Item) {
	fmt.Printf("requesting ID: %s", con.ID)
	o, err := b.dlna.Browse(con)

	if err != nil {
		println("Error: ", err)
		return
	}
	println("done")
	*conList = &o.Containers
	*itemList = &o.Items
}
