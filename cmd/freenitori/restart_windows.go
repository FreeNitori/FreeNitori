package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
)

func abort() {
	<-state.Exit
}

func restart() {
	if _, err := os.Stat(path); err != nil {
		log.Fatalf("Error stat executable path, %s", err)
		os.Exit(1)
	}
	log.Infof("Executable found at %s.", path)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory, %s", err)
		os.Exit(1)
	}
	log.Infof("Current working directory is %s.", wd)

	_, err = os.StartProcess(path, []string{}, &os.ProcAttr{
		Dir:   wd,
		Env:   nil,
		Files: []*os.File{os.Stderr, os.Stdin, os.Stdout},
		Sys:   nil,
	})
	if err != nil {
		log.Fatalf("Error creating new process, %s", err)
		os.Exit(1)
	}
}
