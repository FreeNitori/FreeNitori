package main

import (
	log "git.randomchars.net/FreeNitori/FreeNitori/Log"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"os"
)

func abnormalExit() {
	<-state.ExitCode
}

func restart() {
	if _, err := os.Stat(execPath); err != nil {
		log.Fatalf("Unable to get executable path, %s", err)
		os.Exit(1)
	}
	log.Infof("Program found at %s.", execPath)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory, %s", err)
		os.Exit(1)
	}
	log.Infof("Current working directory is %s.", wd)
	_, err = os.StartProcess(execPath, []string{}, &os.ProcAttr{
		Dir:   wd,
		Env:   nil,
		Files: []*os.File{os.Stderr, os.Stdin, os.Stdout},
		Sys:   nil,
	})
	if err != nil {
		log.Fatalf("Unable to create new process, %s", err)
		os.Exit(1)
	}
}
