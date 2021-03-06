package ui

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var mw *walk.MainWindow

// Serve serves the GUI.
func Serve() {
	size := Size{Width: 800, Height: 400}

	err := MainWindow{
		AssignTo: &mw,
		Title:    "FreeNitori " + state.Version(),
		MinSize:  size,
		Layout: VBox{
			MarginsZero: true,
		},
		Children: []Widget{
			TextEdit{AssignTo: &logEdit, ReadOnly: true, TextColor: walk.Color(0xffffff), Background: SolidColorBrush{Color: walk.Color(0x000000)}},
			PushButton{Text: "Restart", OnClicked: func() {
				state.Exit <- -1
				_ = mw.Close()
			}},
		},
	}.Create()

	if err != nil {
		log.Fatalf("Error creating GUI, %s", err)
		state.Exit <- 1
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
	logEdit.AppendText(earlyBuffer)
	windowInitFinish = true
	mw.Run()
	state.Exit <- 0
}
