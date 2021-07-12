package main

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
)

var path string

func init() {
	var err error
	if path, err = os.Executable(); err != nil {
		println("error looking up executable path, " + err.Error() + ", restart will not be operative")
	}
}

func cleanup() {
	log.Info("Cleaning up...")

	// Close RPC listener
	if rpcListener != nil {
		if err := (*rpcListener).Close(); err != nil {
			log.Errorf("Error closing RPC listener, %s", err)
		}
	}

	// Shutdown web server
	if err := web.Close(); err != nil {
		log.Errorf("Error closing web server, %s", err)
	}

	// Close Discord
	if err := discord.Close(); err != nil {
		log.Errorf("Error closing Discord, %s", err)
	}

	// Close database
	if err := database.Database.Close(); err != nil {
		log.Errorf("Error closing database, %s", err)
	}
}
