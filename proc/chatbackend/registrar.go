package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/state"
)

func init() {
	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range state.EventHandlers {
			state.RawSession.AddHandler(handler)
		}
	}

	// Add the event handlers
	for _, route := range multiplexer.Commands {
		multiplexer.Router.Route(route)
	}
}
