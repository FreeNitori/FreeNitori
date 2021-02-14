package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func init() {
	writer, err := rotatelogs.New(
		config.Config.System.LogPath+"/freenitori.%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(config.Config.System.LogPath+"/freenitori.log"),
		rotatelogs.WithRotationTime(86400*time.Second),
		rotatelogs.WithMaxAge(604800*time.Second))
	if err != nil {
		log.Fatalf("Unable to initialize disk-based logging, %s", err)
		os.Exit(1)
	}
	log.Logger.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
			logrus.ErrorLevel: writer,
			logrus.FatalLevel: writer,
			logrus.InfoLevel:  writer,
			logrus.PanicLevel: writer,
			logrus.TraceLevel: writer,
			logrus.WarnLevel:  writer,
		},
		&logrus.JSONFormatter{}))
}
