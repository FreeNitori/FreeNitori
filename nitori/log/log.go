// Logging functions and logger object.
package log

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
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
	switch vars.ProcessType {
	case vars.Supervisor:
		return append([]byte("[SV]"), format...), err
	case vars.ChatBackend:
		return append([]byte("[CB]"), format...), err
	case vars.WebServer:
		return append([]byte("[WS]"), format...), err
	case vars.Other:
		return format, err
	default:
		panic("invalid process type")
	}
}
