package main

import (
	"fmt"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/args"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/communication"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/state"
	"github.com/dgraph-io/badger/v2"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var err error

func main() {
	vars.ProcessType = vars.Supervisor

	// Print version information and stuff
	log.Infof("Starting FreeNitori %s", vars.Version)

	// Check for an existing instance
	if _, err := os.Stat(config.Config.System.Socket); os.IsNotExist(err) {
	} else {
		_, err := net.Dial("unix", config.Config.System.Socket)
		if err != nil {
			err = syscall.Unlink(config.Config.System.Socket)
			if err != nil {
				log.Error(fmt.Sprintf("Unable to remove hanging socket, %s", err))
				os.Exit(1)
			}
		} else {
			log.Error("Another instance of FreeNitori is already running.")
			os.Exit(1)
		}
	}

	// Start RPC server
	go func() {
		err = communication.Server()
		if err != nil {
			log.Fatalf("Failed to start RPC server, %s", err)
		}
	}()

	// Open the database
	dbOptions := badger.DefaultOptions(config.Config.System.Database)
	dbOptions.Logger = log.Logger
	state.Database, err = badger.Open(dbOptions)
	if err != nil {
		log.Fatalf("Failed to open database, %s", err)
		os.Exit(1)
	}
	defer func() { _ = state.Database.Close() }()

	// Create the chat backend process
	state.ChatBackendProcess, err =
		os.StartProcess(config.Config.System.ChatBackend, []string{config.Config.System.ChatBackend, "-a", config.TokenOverride, "-c", config.NitoriConfPath}, &state.ProcessAttributes)
	if err != nil {
		log.Errorf("Failed to create chat backend process, %s", err)
		os.Exit(1)
	}

	// Create web server process
	state.WebServerProcess, err =
		os.StartProcess(config.Config.System.WebServer, []string{config.Config.System.WebServer, "-a", config.TokenOverride, "-c", config.NitoriConfPath}, &state.ProcessAttributes)
	if err != nil {
		log.Errorf("Failed to create web server process, %s", err)
		os.Exit(1)
	}
	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt, os.Kill)
	go func() {
		for {
			currentSignal := <-signalChannel
			switch currentSignal {
			case syscall.SIGUSR2:
				vars.ExitCode <- 0
				break
			default:
				// Cleanup stuffs
				fmt.Print("\n")
				log.Info("Gracefully terminating...")
				_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
				_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
				_ = state.SocketListener.Close()
				_ = syscall.Unlink(config.Config.System.Socket)
				vars.ExitCode <- 0
				break
			}
		}
	}()

	// Exit if there's something on that channel
	exitCode := <-vars.ExitCode
	os.Exit(exitCode)
}
