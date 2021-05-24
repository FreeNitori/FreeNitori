package main

import (
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type out struct{}

func (out) Write(p []byte) (n int, err error) {
	return len(p), ioutil.WriteFile("/dev/tty1", p, 0644)
}

func init() {
	log.Instance.SetOutput(out{})
	log.Instance.Formatter = &logrus.TextFormatter{
		ForceColors:               true,
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
	}
}

func main() {
	log.Info("Nitori!")
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
