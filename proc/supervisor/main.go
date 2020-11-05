// Supervisor program that house the database and mediate communication between server programs.
package main

import (
	"fmt"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/args"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/communication"
	"git.randomchars.net/RandomChars/FreeNitori/proc/supervisor/state"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"plugin"
	"strings"
	"syscall"
)

var err error

func init() {
	vars.ProcessType = vars.Supervisor
	func() {
		stat, err := os.Stat("plugins")
		if os.IsNotExist(err) {
			err = os.Mkdir("plugins", 0755)
			if err != nil {
				log.Fatalf("Failed to create plugin directory, %s", err)
				os.Exit(1)
			}
			return
		}
		if !stat.IsDir() {
			log.Fatal("Plugin path is not a directory.")
			os.Exit(1)
		}
		pluginPaths, err := ioutil.ReadDir("plugins/")
		if err != nil {
			log.Fatalf("Unable to read plugin directory, %s", err)
			os.Exit(1)
		}
		for _, path := range pluginPaths {
			if !strings.HasSuffix(path.Name(), ".so") {
				continue
			}
			pl, err := plugin.Open("plugins/" + path.Name())
			if err != nil {
				log.Warnf("Error while loading plugin %s, %s", path.Name(), err)
				continue
			}
			symbol, err := pl.Lookup("Database")
			if err != nil {
				continue
			}
			db, ok := symbol.(state.DatabaseBackend)
			if !ok {
				log.Warnf("No DatabaseBackend found in %s.", path.Name())
				continue
			}
			if state.Database != nil {
				log.Warnf("Already loaded database backend %s, skipping plugin %s.", state.Database.DBType(), path.Name())
				continue
			}
			state.Database = db
			log.Infof("Loaded plugin %s implementing database backend %s.", path.Name(), db.DBType())
		}
	}()
	if state.Database == nil {
		log.Fatal("Please place a database driver in the plugins directory.")
		os.Exit(1)
	}
}

func main() {
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
			log.Fatalf("RPC server encountered an error, %s", err)
		}
		defer func() { _ = state.SocketListener.Close() }()
	}()

	// Open the database
	err = state.Database.Open(config.Config.System.Database)
	if err != nil {
		log.Fatalf("Failed to open database, %s", err)
		os.Exit(1)
	}
	defer func() { _ = state.Database.Close() }()

	// Create chat backend process
	state.ChatBackendProcess, err =
		os.StartProcess(config.Config.System.ChatBackend, append([]string{config.Config.System.ChatBackend}, state.ServerArgs...), &state.ProcessAttributes)
	if err != nil {
		log.Errorf("Failed to create chat backend process, %s", err)
		os.Exit(1)
	}

	// Create web server process
	state.WebServerProcess, err =
		os.StartProcess(config.Config.System.WebServer, append([]string{config.Config.System.WebServer}, state.ServerArgs...), &state.ProcessAttributes)
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
