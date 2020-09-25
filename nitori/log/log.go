package log

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var Formatter = formatter{logrus.TextFormatter{
	ForceColors:               false,
	DisableColors:             false,
	ForceQuote:                false,
	DisableQuote:              false,
	EnvironmentOverrideColors: false,
	DisableTimestamp:          false,
	FullTimestamp:             true,
	TimestampFormat:           "",
	DisableSorting:            true,
	SortingFunc:               nil,
	DisableLevelTruncation:    false,
	PadLevelText:              false,
	QuoteEmptyFields:          false,
	FieldMap:                  nil,
	CallerPrettyfier:          nil,
}}

type formatter struct {
	logrus.TextFormatter
}

func init() {
	logger.SetFormatter(&Formatter)
	switch config.Debug {
	case true:
		logger.SetLevel(logrus.DebugLevel)
	case false:
		logger.SetLevel(logrus.InfoLevel)
	}
}

func (formatter *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	format, err := formatter.TextFormatter.Format(entry)
	switch {
	case state.StartChatBackend:
		return append([]byte("[CB]"), format...), err
	case state.StartWebServer:
		return append([]byte("[WS]"), format...), err
	case !state.StartWebServer && !state.StartChatBackend:
		return append([]byte("[SV]"), format...), err
	default:
		panic("invalid start parameters")
	}
}
