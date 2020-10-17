package log

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

func SetLevel(level logrus.Level) {
	Logger.SetLevel(level)
}

func GetLevel() logrus.Level {
	return Logger.GetLevel()
}

func IsLevelEnabled(level logrus.Level) bool {
	return Logger.IsLevelEnabled(level)
}

func DiscordGoLogger(msgL, _ int, format string, a ...interface{}) {
	var level logrus.Level
	switch msgL {
	case discordgo.LogDebug:
		level = logrus.DebugLevel
	case discordgo.LogInformational:
		level = logrus.InfoLevel
	case discordgo.LogWarning:
		level = logrus.WarnLevel
	case discordgo.LogError:
		level = logrus.ErrorLevel
	}
	Logger.Log(level, fmt.Sprintf(format, a...))
}
