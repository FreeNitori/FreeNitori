package ui

import (
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/lxn/walk"
	"github.com/sirupsen/logrus"
	"go/types"
)

var earlyBuffer string
var windowInitFinish = false
var logEdit *walk.TextEdit

type windowLogViewHook types.Nil

func (w windowLogViewHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (w windowLogViewHook) Fire(entry *logrus.Entry) error {
	if windowInitFinish {
		logEdit.AppendText(entry.Message + "\r\n")
	} else {
		earlyBuffer += entry.Message + "\r\n"
	}
	return nil
}

func init() {
	log.Instance.AddHook(windowLogViewHook{})
}
