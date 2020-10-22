package main

import (
	"flag"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
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
	vars.ProcessType = vars.Other

	// Dial the supervisor socket
	vars.RPCConnection, err = rpc.DialHTTP("unix", socketPath)
	if err != nil {
		log.Errorf("Failed to connect to the supervisor process, %s", err)
		os.Exit(1)
	}
	defer func() { _ = vars.RPCConnection.Close() }()

	// Initialize the shell prompt
	go shell()

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt, os.Kill)
	go func() {
		<-signalChannel
		vars.ExitCode <- 0
	}()

	// Exit if there's something on that channel
	exitCode := <-vars.ExitCode
	os.Exit(exitCode)
}
