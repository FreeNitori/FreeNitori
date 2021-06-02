// +build !windows

package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
	"syscall"
)

func abnormalExit() {}

func restart() {
	if _, err := os.Stat(execPath); err != nil {
		log.Fatalf("Unable to get executable path, %s", err)
		os.Exit(1)
	}
	if state.Reincarnation != "" {
		log.Infof("Setting reincarnation payload %s", state.Reincarnation)
		err = os.Setenv("REINCARNATION", state.Reincarnation)
		if err != nil {
			log.Errorf("Error occurred while setting reincarnation payload, %s", err)
		}
	}
	log.Infof("Program found at %s, re-executing...", execPath)
	err = syscall.Exec(execPath, os.Args, os.Environ())
	if err != nil {
		log.Fatalf("Error occurred while re-executing, %s", err)
		os.Exit(1)
	}
}
