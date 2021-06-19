// Server program.
package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/db"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/extension"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/plugin"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/rpc"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/ui"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
	"os/signal"
	"syscall"
)

var err error

func init() {

	// Print initial message and stuff
	log.Infof("FreeNitori %s (%s) early initialization.", state.Version(), state.Revision())

	// Initialize plugins
	err = plugin.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize plugins, %s", err)
		os.Exit(1)
	}

	// Initialize RPC
	err = rpc.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize RPC server, %s", err)
		os.Exit(1)
	}

	// Initialize database
	err = db.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize database, %s", err)
		os.Exit(1)
	}

	// Initialize web services
	err = web.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize web services, %s", err)
		_ = database.Database.Close()
		os.Exit(1)
	}

	// Initialize Discord-related services
	err = discord.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize Discord services, %s", err)
		_ = database.Database.Close()
		os.Exit(1)
	}

	// Initialize UI
	err = ui.Initialize()
	if err != nil {
		log.Fatalf("Unable to initialize UI, %s", err)
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
	go ui.Serve()

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
		log.Fatalf("Unable to register extensions, %s", err)
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

	// Block until exit code is received, and exit with that code or restart if it is -1.
	exitCode := <-state.ExitCode
	if exitCode != 0 {
		if exitCode == -1 {
			cleanup()
			restart()
		}
		abnormalExit()
		os.Exit(exitCode)
	}
}
