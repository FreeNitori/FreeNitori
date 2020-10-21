package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/communication"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	SuperVisor "git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
	"github.com/dgraph-io/badger/v2"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

var err error

func main() {
	state.ProcessType = state.Supervisor

	// IPC functions
	var IPCFunctions = new(communication.IPC)

	// Print version information and stuff
	log.Infof("Starting FreeNitori %s", state.Version)

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

	// Initialize the socket
	_ = rpc.Register(IPCFunctions)
	rpc.HandleHTTP()
	SuperVisor.SocketListener, err = net.Listen("unix", config.Config.System.Socket)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to listen on the socket, %s", err))
		os.Exit(1)
	}
	go func() { _ = http.Serve(SuperVisor.SocketListener, nil) }()

	// Open the database
	dbOptions := badger.DefaultOptions(config.Config.System.Database)
	dbOptions.Logger = log.Logger
	SuperVisor.Database, err = badger.Open(dbOptions)
	if err != nil {
		log.Fatalf("Failed to open database, %s", err)
		os.Exit(1)
	}
	defer func() { _ = SuperVisor.Database.Close() }()

	// Create the chat backend process
	SuperVisor.ChatBackendProcess, err =
		os.StartProcess(config.Config.System.ChatBackend, []string{config.Config.System.ChatBackend, "-a", ChatBackend.RawSession.Token, "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
	if err != nil {
		log.Errorf("Failed to create chat backend process, %s", err)
		os.Exit(1)
	}

	// Create web server process
	SuperVisor.WebServerProcess, err =
		os.StartProcess(config.Config.System.WebServer, []string{config.Config.System.WebServer, "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
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
				state.ExitCode <- 0
				return
			default:
				// Cleanup stuffs
				fmt.Print("\n")
				log.Info("Gracefully terminating...")
				_ = SuperVisor.ChatBackendProcess.Signal(syscall.SIGUSR2)
				_ = SuperVisor.WebServerProcess.Signal(syscall.SIGUSR2)
				_ = SuperVisor.SocketListener.Close()
				_ = syscall.Unlink(config.Config.System.Socket)
				state.ExitCode <- 0
				return
			}
		}
	}()

	// Exit if there's something on that channel
	exitCode := <-state.ExitCode
	os.Exit(exitCode)
}
