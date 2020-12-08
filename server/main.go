// Server program.
package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/binaries/confdefault"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/database"
	"git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord"
	"git.randomchars.net/RandomChars/FreeNitori/server/web"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var err error

func init() {
	// Print version information and stuff
	log.Infof("FreeNitori %s early initialization.", state.Version)

	// Check for configuration, or generate default config if nonexistent
	if config.Config == nil {
		defaultConfigFile, err := confdefault.Asset("nitori.conf")
		if err != nil {
			log.Fatalf("Failed to extract the default configuration file, %s", err)
			os.Exit(1)
		}
		err = ioutil.WriteFile("nitori.conf", defaultConfigFile, 0644)
		if err != nil {
			log.Fatalf("Failed to write the default configuration file, %s", err)
			os.Exit(1)
		}
		log.Fatalf("Generated default configuration file at ./nitori.conf, " +
			"please edit it before starting FreeNitori.")
		os.Exit(1)
	}

	// Check for existence of plugin directory
	_, err := os.Stat("plugins")
	if os.IsNotExist(err) {
		err = os.Mkdir("plugins", 0755)
		if err != nil {
			log.Fatalf("Failed to create plugin directory, %s", err)
			os.Exit(1)
		}
	}

	// Initialize database
	err = database.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize database, %s", err)
		os.Exit(1)
	}

	// Initialize web services
	err = web.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize web services, %s", err)
		_ = vars.Database.Close()
		os.Exit(1)
	}

	// Initialize discord-related services
	err = discord.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize discord-related services, %s", err)
		_ = vars.Database.Close()
		os.Exit(1)
	}
}

func main() {

	// Cleanup after exit
	defer cleanup()

	// Start Discord and web services routines
	go discord.Serve()
	go web.Serve()

	// Print thing
	log.Info("Begin late initialization.")

	// Late initialization of discord-related services
	err = discord.LateInitialize()
	if err != nil {
		log.Fatalf("Failed to initialize discord-related services, %s", err)
		_ = vars.Database.Close()
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
		os.Exit(exitCode)
	}
}
