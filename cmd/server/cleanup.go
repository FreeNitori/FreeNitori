package main

import (
	"context"
	dbVars "git.randomchars.net/RandomChars/FreeNitori/cmd/server/database/vars"
	dcVars "git.randomchars.net/RandomChars/FreeNitori/cmd/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/cmd/server/rpc"
	"git.randomchars.net/RandomChars/FreeNitori/cmd/server/web"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"os"
	"syscall"
	"time"
)

func init() {
	execPath, _ = os.Executable()
}

var execPath string

func cleanup() {
	log.Info("Running cleanups.")

	// Close RPC connection
	err = rpc.Listener.Close()
	if err != nil {
		log.Errorf("Error while closing RPC listener, %s", err)
	}

	// Shutdown web server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = web.Server.Shutdown(ctx)
	if err != nil {
		log.Errorf("Error while shutting down web server, %s", err)
	}

	// Close Discord sessions
	for index, shardSession := range dcVars.ShardSessions {
		err = shardSession.Close()
		if err != nil {
			log.Errorf("Error while shutting down shard %v, %s", index, err)
		}
	}
	err = dcVars.RawSession.Close()
	if err != nil {
		log.Errorf("Error while closing session with Discord, %s", err)
	}

	// Close database
	err = dbVars.Database.Close()
	if err != nil {
		log.Errorf("Error while closing database, %s", err)
	}
}

func restart() {
	if _, err := os.Stat(execPath); err != nil {
		log.Fatalf("Failed to get executable path, %s", err)
		os.Exit(1)
	}
	log.Infof("Program found at %s, re-executing...", execPath)
	err = syscall.Exec(execPath, os.Args, os.Environ())
	if err != nil {
		log.Fatalf("Failed to re-execute, %s", err)
		os.Exit(1)
	}
}
