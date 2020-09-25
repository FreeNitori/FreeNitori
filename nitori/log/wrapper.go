package log

import "github.com/sirupsen/logrus"

func Info(args ...interface{}) {
	logger.Info(args)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
}

func Warn(args ...interface{}) {
	logger.Warn(args)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args)
}

func Error(args ...interface{}) {
	logger.Error(args)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args)
}

func SetLevel(level logrus.Level) {
	logger.SetLevel(level)
}

func GetLevel() logrus.Level {
	return logger.GetLevel()
}

func IsLevelEnabled(level logrus.Level) bool {
	return logger.IsLevelEnabled(level)
}
