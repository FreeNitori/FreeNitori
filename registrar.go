package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Shard {
		multiplexer.RawSession.AddHandler(multiplexer.Router.OnMessageCreate)
	}

	// Add the routes
	for _, handlerMeta := range handlers.AllHandlers {
		multiplexer.Router.Route(
			handlerMeta.Pattern,
			handlerMeta.AliasPatterns,
			handlerMeta.Description,
			handlerMeta.Handler,
			handlerMeta.Category)
	}
}
