package log

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Info logs message on info level.
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof logs message on info level.
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Debug logs message on debug level.
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf logs message on debug level.
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Warn logs message on warn level.
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf logs message on warn level.
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error logs message on error level.
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf logs message on error level.
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal logs message on fatal level.
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf logs message on fatal level.
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// SetLevel sets the logger level.
func SetLevel(level logrus.Level) {
	Logger.SetLevel(level)
}

// GetLevel returns the logger level.
func GetLevel() logrus.Level {
	return Logger.GetLevel()
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param.
func IsLevelEnabled(level logrus.Level) bool {
	return Logger.IsLevelEnabled(level)
}

// DiscordGoLogger overrides logger of discordgo library.
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
