// Server program.
package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/database"
	"git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord"
	"git.randomchars.net/RandomChars/FreeNitori/server/extension"
	"git.randomchars.net/RandomChars/FreeNitori/server/rpc"
	"git.randomchars.net/RandomChars/FreeNitori/server/web"
	"os"
	"os/signal"
	"syscall"
)

var err error

func init() {
	if config.VersionStartup {
		println(state.Version() + " (" + state.Revision() + ")")
		os.Exit(0)
	}

	// Print initial message and stuff
	log.Infof("FreeNitori %s (%s) early initialization.", state.Version(), state.Revision())

	// Check for existence of plugin directory
	_, err := os.Stat("plugins")
	if os.IsNotExist(err) {
		err = os.Mkdir("plugins", 0755)
		if err != nil {
			log.Fatalf("Unable to create plugin directory, %s", err)
			os.Exit(1)
		}
	}

	// Initialize RPC
	err = rpc.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize RPC server, %s", err)
		os.Exit(1)
	}

	// Initialize database
	err = database.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize database, %s", err)
		os.Exit(1)
	}

	// Initialize web services
	err = web.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize web services, %s", err)
		_ = vars.Database.Close()
		os.Exit(1)
	}

	// Initialize Discord-related services
	err = discord.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize Discord services, %s", err)
		_ = vars.Database.Close()
		os.Exit(1)
	}
}

func main() {

	// Cleanup after exit
	defer cleanup()

	// Start service routines
	go rpc.Serve()
	go discord.Serve()
	go web.Serve()

	// Print thing
	log.Info("Begin late initialization.")

	// Late initialization of Discord-related services
	err = discord.LateInitialize()
	if err != nil {
		log.Fatalf("Unable to initialize Discord services, %s", err)
		cleanup()
		os.Exit(1)
	}

	// Load Discord extensions
	err = extension.FindExtensions()
	if err != nil {
		log.Fatalf("Unable to find extensions, %s", err)
		cleanup()
		os.Exit(1)
	}
	err = extension.RegisterHandlers()
	if err != nil {
		log.Fatalf("Unable to register event handlers, %s", err)
		cleanup()
		os.Exit(1)
	}

	// Print thing
	log.Info("Late initialization completed.")

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt, os.Kill)
	go func() {
		var exit int
		defer func() { state.ExitCode <- exit }()
		for {
			currentSignal := <-signalChannel
			switch currentSignal {
			case os.Interrupt:
				exit = 0
				println()
				log.Info("Gracefully exiting.")
				return
			default:
				exit = 0
				log.Info("Gracefully exiting.")
				return
			}
		}
	}()

	// Exit if there's something on that channel
	exitCode := <-state.ExitCode
	if exitCode != 0 {
		if exitCode == -1 {
			cleanup()
			restart()
		}
		os.Exit(exitCode)
	}
}
