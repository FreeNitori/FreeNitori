package discord

import (
	_ "git.randomchars.net/RandomChars/FreeNitori/cmd/server/discord/internals"
	"git.randomchars.net/RandomChars/FreeNitori/cmd/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/routes"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range multiplexer.EventHandlers {
			vars.RawSession.AddHandler(handler)
		}
	}
}
