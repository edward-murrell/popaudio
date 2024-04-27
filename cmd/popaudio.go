package main

import (
	"context"
	"fmt"
	"os"
	"popaudio/internal/dlnaclient"
	"popaudio/internal/ui"
)

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		fmt.Println(`"Call with path to server eg http://server/file.mp3`)
		os.Exit(1)
	}

	dlnaClient, err := dlnaclient.GetServerByName(ctx, os.Args[1])
	if err != nil {
		println("error getting DLNA client", err)
	}
	ui.NewScanLineBrowser(dlnaClient).Run()
}
