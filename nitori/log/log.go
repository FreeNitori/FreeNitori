package log

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()
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
	Logger.SetFormatter(&Formatter)
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
