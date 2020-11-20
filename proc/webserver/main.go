// Server program that serves things over HTTP.
package main

import (
	"fmt"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/args"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/ipc"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var err error
var readyChannel = make(chan bool, 1)

func init() {
	vars.ProcessType = vars.WebServer
}

func main() {
	// Dial the supervisor socket
	err = ipc.InitializeIPC()
	if err != nil {
		log.Errorf("Failed to connect to the supervisor process, %s", err)
		os.Exit(1)
	}
	defer func() { _ = vars.RPCConnection.Close() }()

	// Initialize and start the server
	Initialize()
	go func() {
		<-readyChannel
		// Get some constant data
		vars.InviteURL = fetchData("inviteURL")

		// Start the server
		log.Infof("Web server listening on %s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port))
		err = Engine.Run(fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)))
		if err != nil {
			log.Error(fmt.Sprintf("Failed to start web server, %s", err))
			_ = vars.RPCConnection.Call("R.Error", []string{"WebServer"}, nil)
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
				if !vars.Initialized {
					readyChannel <- true
					vars.Initialized = true
				}
			case syscall.SIGUSR2:
				vars.ExitCode <- 0
				break
			default:
				// Cleanup stuffs
				if currentSignal != os.Interrupt {
					// Only write the packet if SIGUSR2 was not sent or the program was not interrupted
					_ = vars.RPCConnection.Call("R.Restart", []string{"WebServer"}, nil)
				}
				vars.ExitCode <- 0
				break
			}
		}
	}()

	// Tell the Supervisor and exit if there's something on that channel
	exitCode := <-vars.ExitCode
	if exitCode != 0 {
		_ = vars.RPCConnection.Call("R.Error", []string{"WebServer"}, nil)
	}
	os.Exit(exitCode)
}
