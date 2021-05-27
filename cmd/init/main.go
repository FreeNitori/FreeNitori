package main

import (
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

type out struct{}

func (out) Write(p []byte) (n int, err error) {
	if ioutil.WriteFile("/dev/console", p, 0644) == nil {
		return len(p), ioutil.WriteFile("/dev/tty1", p, 0644)
	} else {
		panic("cannot write to console")
	}
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

var err error

func main() {
	// Only run as init
	if os.Getpid() != 1 {
		println("This program must be ran as PID 0.")
		os.Exit(9)
	}

	log.Info("FreeNitori System Management and Initialization program starting.")

	startServer("/bin/freenitori")

	// Shutdown
	err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		panic(err)
	}

}

func startServer(path string) {
	// Start server
	s := exec.Command(path)
	s.Stdout = out{}
	s.Stderr = out{}
	s.Stdin = os.Stdin
	err = s.Run()
	if err != nil {
		log.Errorf("Unable to start server, %s.", err)
		err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
		if err != nil {
			panic(err)
		}
	}
}
