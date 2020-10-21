package main

import (
	"flag"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

var socketPath string
var err error

func init() {
	flag.StringVar(&socketPath, "s", "/tmp/nitori", "Nitori socket path")
	flag.Parse()
}

func main() {
	state.ProcessType = state.InteractiveConsole

	// Dial the supervisor socket
	state.IPCConnection, err = rpc.DialHTTP("unix", socketPath)
	if err != nil {
		log.Errorf("Failed to connect to the supervisor process, %s", err)
		os.Exit(1)
	}
	defer func() { _ = state.IPCConnection.Close() }()

	// Initialize the console
	go console()

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt, os.Kill)
	go func() {
		<-signalChannel
		state.ExitCode <- 0
	}()

	// Exit if there's something on that channel
	exitCode := <-state.ExitCode
	os.Exit(exitCode)
}
