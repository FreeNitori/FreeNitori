// +build !windows

package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
	"syscall"
)

func abort() {}

func restart() {
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("Error stat executable path, %s", err)
		os.Exit(1)
	}

	// Set reincarnation payload for success message on next startup
	if state.Reincarnation != "" {
		log.Infof("Setting reincarnation payload: %s", state.Reincarnation)
		if err := os.Setenv("REINCARNATION", state.Reincarnation); err != nil {
			log.Errorf("Error setting reincarnation payload, %s", err)
		}
	}

	// Do execve(2)
	log.Infof("Executable found at %s, re-executing...", path)
	if err := syscall.Exec(path, os.Args, os.Environ()); err != nil {
		log.Fatalf("Error occurred while re-executing, %s", err)
		os.Exit(1)
	}
}
