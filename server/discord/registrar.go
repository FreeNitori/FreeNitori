package discord

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/routes"
	_ "git.randomchars.net/RandomChars/FreeNitori/server/discord/internals"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range multiplexer.EventHandlers {
			vars.RawSession.AddHandler(handler)
		}
	}
}
