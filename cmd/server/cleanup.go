package main

import (
	"context"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/rpc"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"os"
	"time"
)

func init() {
	execPath, _ = os.Executable()
}

var execPath string

func cleanup() {
	log.Info("Running cleanups.")

	// Close RPC connection
	if rpc.Listener != nil {
		err = rpc.Listener.Close()
		if err != nil {
			log.Errorf("Error while closing RPC listener, %s", err)
		}
	}

	// Shutdown web server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = web.Server.Shutdown(ctx)
	if err != nil {
		log.Errorf("Error while shutting down web server, %s", err)
	}

	// Close Discord sessions
	for index, shardSession := range state.ShardSessions {
		err = shardSession.Close()
		if err != nil {
			log.Errorf("Error while shutting down shard %v, %s", index, err)
		}
	}
	err = state.RawSession.Close()
	if err != nil {
		log.Errorf("Error while closing session with Discord, %s", err)
	}

	// Close database
	err = database.Database.Close()
	if err != nil {
		log.Errorf("Error while closing database, %s", err)
	}
}
