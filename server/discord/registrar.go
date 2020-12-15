package discord

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	_ "git.randomchars.net/RandomChars/FreeNitori/server/discord/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
)

func init() {
	log.Info("Registering event handlers.")

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range multiplexer.EventHandlers {
			vars.RawSession.AddHandler(handler)
		}
	}

	// Add the event handlers
	for _, route := range multiplexer.Commands {
		multiplexer.Router.Route(route)
		log.Debugf("Registered route with pattern '%s'.", route.Pattern)
	}
}
