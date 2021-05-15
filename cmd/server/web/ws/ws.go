package ws

import (
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/sirupsen/logrus"
	"go/types"
	"gopkg.in/olahol/melody.v1"
)

var err error

var WS = melody.New()

type websocketLogBroadcastHook types.Nil

func (w websocketLogBroadcastHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (w websocketLogBroadcastHook) Fire(entry *logrus.Entry) error {
	err = WS.Broadcast([]byte(entry.Message))
	if err != nil {
		log.Errorf("Error while broadcasting log entry, %s", err)
	}
	return nil
}

func init() {
	log.Instance.AddHook(websocketLogBroadcastHook{})
}
