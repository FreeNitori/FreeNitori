// Program to interact with Nitori through a shell prompt.
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
var action string
var err error

func init() {
	flag.StringVar(&action, "c", "", "Execute command.")
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
	if len(os.Args) == 2 {
		if os.Args[1] == "shutdown" {
			_ = vars.RPCConnection.Call("R.Shutdown", []int{vars.ProcessType}, nil)
			os.Exit(0)
		}
	}

	// Initialize the shell prompt
	go initShell()

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
