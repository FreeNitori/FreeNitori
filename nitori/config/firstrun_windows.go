package config

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"os"
)

func firstRun(edit bool) {
	var message = "Configuration file created, please edit before restarting Nitori."
	if edit {
		message = "Please edit the configuration file before starting Nitori."
	}
	// Show first run dialogue.
	size := Size{Width: 400, Height: 210}
	var mw *walk.MainWindow

	err := MainWindow{
		AssignTo: &mw,
		Title:    "FreeNitori " + state.Version(),
		Size:     size,
		Layout: Flow{
			Margins: Margins{
				Left:   100,
				Top:    100,
				Right:  100,
				Bottom: 100,
			},
		},
		Children: []Widget{
			TextLabel{Text: message},
		},
	}.Create()

	if err != nil {
		log.Fatalf("Unable to create first run dialog, %s", err)
		state.ExitCode <- 1
	}

	hIcon := win.LoadImage(
		win.GetModuleHandle(nil),
		win.MAKEINTRESOURCE(2),
		win.IMAGE_ICON,
		0, 0,
		win.LR_DEFAULTSIZE)
	if hIcon != 0 {
		win.SendMessage(mw.Handle(), win.WM_SETICON, 1, uintptr(hIcon))
		win.SendMessage(mw.Handle(), win.WM_SETICON, 0, uintptr(hIcon))
	}

	mw.Run()
	os.Exit(0)
}
