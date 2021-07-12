package main

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	log "git.randomchars.net/FreeNitori/Log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

var (
	writer    io.Writer
	hook      *logrusHook
	formatter = logrus.JSONFormatter{}
)

type logrusHook struct {
	lock *sync.Mutex
}

func (hook *logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *logrusHook) Fire(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()
	var message []byte
	var err error
	message, err = formatter.Format(entry)
	if err != nil {
		fmt.Printf("Error formatting log message, %s", err)
		return err
	}
	_, err = writer.Write(message)
	return err
}

func setupLogRotate() {
	var err error
	writer, err = rotatelogs.New(
		config.System.LogPath+"/freenitori.%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(config.System.LogPath+"/freenitori.log"),
		rotatelogs.WithRotationTime(86400*time.Second),
		rotatelogs.WithMaxAge(604800*time.Second))
	if err != nil {
		log.Fatalf("Error initializing log rotation, %s", err)
		os.Exit(1)
	}
	hook = &logrusHook{lock: new(sync.Mutex)}
	log.Instance.AddHook(hook)
}
