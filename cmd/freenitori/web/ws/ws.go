// Package ws implements websocket operations.
package ws

import (
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/sirupsen/logrus"
	"go/types"
	"gopkg.in/olahol/melody.v1"
)

var WS = melody.New()

type websocketLogBroadcastHook types.Nil

func (w websocketLogBroadcastHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (w websocketLogBroadcastHook) Fire(entry *logrus.Entry) error {
	if err := WS.Broadcast([]byte(entry.Message)); err != nil {
		log.Errorf("Error broadcasting log entry, %s", err)
	}
	return nil
}

func init() {
	log.Instance.AddHook(websocketLogBroadcastHook{})
}
