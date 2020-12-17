package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	dbVars "git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	dcVars "git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/rpc"
	"syscall"
)

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
