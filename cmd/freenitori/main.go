package main

import (
	"flag"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/extension"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/ui"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var v bool

func init() {
	flag.BoolVar(&v, "v", false, "Display version information and exit")
}

func main() {
	// Parse flags
	flag.Parse()

	// Display version information and exit as required
	if v {
		fmt.Printf("%s (%s)", state.Version(), state.Revision())
		os.Exit(0)
	}

	// Startup message
	log.Infof("FreeNitori %s (%s) on %s %s.", state.Version(), state.Revision(), runtime.GOOS, runtime.GOARCH)

	// Parse config
	if err := config.ReadConfig(); err != nil {
		if err == config.ErrFirstRun {
			log.Warn("Default configuration file generated, edit before next startup.")
			os.Exit(1)
		}
		log.Fatalf("Error reading config, %s", err)
		os.Exit(1)
	}

	// Setup log rotation
	setupLogRotate()

	// Check config file for first run placeholders
	config.CheckPlaceholders()

	// Start GUI
	go ui.Serve()

	// Plugin setup
	if err := setupPlugins(); err != nil {
		log.Fatalf("Error setting up plugins, %s", err)
		os.Exit(1)
	}

	// Open database
	if err := database.Database.Open(config.System.Database); err != nil {
		log.Fatalf("Error opening database, %s", err)
		os.Exit(1)
	}

	// Discord
	if err := discord.Open(); err != nil {
		log.Fatalf("Error setting up Discord, %s", err)
		os.Exit(1)
	}

	// Start RPC server
	if err := startRPC(); err != nil {
		log.Fatalf("Error starting RPC server, %s", err)
		os.Exit(1)
	}

	// Extension load
	if err := extension.Setup(); err != nil {
		log.Fatalf("Error setting up extensions, %s", err)
		os.Exit(1)
	}

	// Web server
	if err := web.Open(); err != nil {
		log.Fatalf("Error setting up web server, %s", err)
		os.Exit(1)
	}

	// Signal handling
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		var exit int
		defer func() { state.Exit <- exit }()
		for {
			s := <-sig
			switch s {
			case syscall.SIGINT:
				exit = 0
				println()
				log.Infof("Gracefully exiting.")
				return
			default:
				exit = 0
				log.Infof("Gracefully exiting.")
				return
			}
		}
	}()

	// Block on exit channel
	c := <-state.Exit
	switch c {
	case 0:
		cleanup()
	case -1:
		cleanup()
		restart()
	default:
		abort()
		os.Exit(c)
	}
}
