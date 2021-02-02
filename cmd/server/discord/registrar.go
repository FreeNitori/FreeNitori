package discord

import (
	_ "git.randomchars.net/RandomChars/FreeNitori/cmd/server/discord/internals"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/routes"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range multiplexer.EventHandlers {
			state.RawSession.AddHandler(handler)
		}
	}
}
