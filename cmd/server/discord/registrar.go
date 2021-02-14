package discord

import (
	// Run all init functions from internals.
	_ "git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/internals"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	// Register all routes.
	_ "git.randomchars.net/FreeNitori/FreeNitori/nitori/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range multiplexer.EventHandlers {
			state.RawSession.AddHandler(handler)
		}
	}
}
