package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	dbVars "git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	dcVars "git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/rpc"
	"os"
	"syscall"
)

func init() {
	execPath, _ = os.Executable()
}

var execPath string

func cleanup() {
	log.Info("Running cleanups.")

	// Close RPC connection
	_ = rpc.Listener.Close()
	_ = syscall.Unlink(config.Config.System.Socket)

	// Close Discord sessions
	for _, shardSession := range dcVars.ShardSessions {
		_ = shardSession.Close()
	}
	_ = dcVars.RawSession.Close()

	// Close database
	_ = dbVars.Database.Close()
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
