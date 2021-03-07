// +build !windows

package main

import (
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
	log.Infof("Program found at %s, re-executing...", execPath)
	err = syscall.Exec(execPath, os.Args, os.Environ())
	if err != nil {
		log.Fatalf("Error occurred while re-executing, %s", err)
		os.Exit(1)
	}
}
