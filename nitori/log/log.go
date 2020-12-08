// Logging functions and logger object.
package log

import (
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