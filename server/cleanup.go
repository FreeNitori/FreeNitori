package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	dbVars "git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	dcVars "git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
)

func cleanup()  {
	log.Info("Cleaning up...")
	// Close Discord sessions
	for _, shardSession := range dcVars.ShardSessions {
		_ = shardSession.Close()
	}
	_ = dcVars.RawSession.Close()

	// Close database
	_ = dbVars.Database.Close()
}
