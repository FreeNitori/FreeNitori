package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Shard {
		for _, handler := range state.EventHandlers {
			state.RawSession.AddHandler(handler)
		}
	}

	// Add the routes
	for _, handlerMeta := range multiplexer.Commands {
		multiplexer.Router.Route(
			handlerMeta.Pattern,
			handlerMeta.AliasPatterns,
			handlerMeta.Description,
			handlerMeta.Handler,
			handlerMeta.Category)
	}
}
