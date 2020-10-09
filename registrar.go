package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
)

func init() {

	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.System.Shard {
		for _, handler := range ChatBackend.EventHandlers {
			ChatBackend.RawSession.AddHandler(handler)
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
