package main

import (
	"flag"
	"git.randomchars.net/RandomChars/FreeNitori/cmd/cli/client"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

var err error
var socketPath string
var action string
var exitCode = make(chan int)

func init() {
	flag.StringVar(&action, "c", "", "Execute command.")
	flag.StringVar(&socketPath, "s", "/tmp/nitori", "Nitori socket path")
	flag.Parse()
}

func main() {

	// Dial the RPC server
	client.Client, err = rpc.DialHTTP("unix", socketPath)
	if err != nil {
		log.Errorf("Unable to connect to Nitori, %s", err)
		os.Exit(1)
	}
	defer func() { _ = client.Client.Close() }()

	if action != "" {
		switch action {
		case "shutdown":
			_ = client.Client.Call("N.Shutdown", []int{}, nil)
		case "restart":
			_ = client.Client.Call("N.Restart", []int{}, nil)
		}
		os.Exit(0)
	}

	// Initialize the shell
	go shell()

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt, os.Kill)
	go func() {
		<-signalChannel
		exitCode <- 0
	}()

	// Exit if there's something on that channel
	exitCode := <-exitCode
	os.Exit(exitCode)
}
