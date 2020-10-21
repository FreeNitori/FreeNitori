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
	switch state.ProcessType {
	case state.Supervisor:
		return append([]byte("[SV]"), format...), err
	case state.ChatBackend:
		return append([]byte("[CB]"), format...), err
	case state.WebServer:
		return append([]byte("[WS]"), format...), err
	case state.InteractiveConsole:
		return append([]byte("[VT]"), format...), err
	default:
		panic("invalid process type")
	}
}
