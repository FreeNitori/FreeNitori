package discord

import (
	// Run all init functions from internals.
	_ "git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/internals"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		state.Multiplexer.SessionRegisterHandlers(state.RawSession)
	}
}
