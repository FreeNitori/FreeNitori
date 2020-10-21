package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/communication"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/web"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var err error
var readyChannel = make(chan bool, 1)

func main() {
	state.ProcessType = state.WebServer

	// Dial the supervisor socket
	err = communication.InitializeIPC()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
		os.Exit(1)
	}

	// Initialize and start the server
	web.Initialize()
	go func() {
		<-readyChannel
		err = web.Engine.Run(fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)))
		if err != nil {
			log.Error(fmt.Sprintf("Failed to start web server, %s", err))
			_ = state.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
			os.Exit(1)
		}
	}()

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt, os.Kill)
	go func() {
		for {
			currentSignal := <-signalChannel
			switch currentSignal {
			case syscall.SIGUSR1:
				// Go to the supervisor to fetch further instructions or set ready status
				if !state.Initialized {
					readyChannel <- true
					state.Initialized = true
				}
			case syscall.SIGUSR2:
				state.ExitCode <- 0
				return
			default:
				// Cleanup stuffs
				if currentSignal != os.Interrupt {
					// Only write the packet if SIGUSR2 was not sent or the program was not interrupted
					_ = state.IPCConnection.Call("IPC.Restart", []string{"WebServer"}, nil)
				}
				state.ExitCode <- 0
				return
			}
		}
	}()

	// Tell the Supervisor and exit if there's something on that channel
	exitCode := <-state.ExitCode
	if exitCode != 0 {
		_ = state.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
	}
	os.Exit(exitCode)
}
