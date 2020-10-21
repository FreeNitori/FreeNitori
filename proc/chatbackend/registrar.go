package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
)

func init() {
	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range ChatBackend.EventHandlers {
			ChatBackend.RawSession.AddHandler(handler)
		}
	}

	// Add the event handlers
	for _, handlerInfo := range multiplexer.Commands {
		multiplexer.Router.Route(
			handlerInfo.Pattern,
			handlerInfo.AliasPatterns,
			handlerInfo.Description,
			handlerInfo.Handler,
			handlerInfo.Category)
	}
}
