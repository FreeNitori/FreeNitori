// Plugin example.
package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/multiplexer"
)

//goland:noinspection GoUnusedGlobalVariable
var CommandRoute = multiplexer.Route{
	Pattern:       "feels",
	AliasPatterns: []string{"feelsgreat"},
	Description:   "Sends a feels great emote.",
	Category:      *handlers.SystemCategory,
	Handler: func(context *multiplexer.Context) {
		context.SendMessage("<:FeelsKappa:713635502741520384>")
	},
}
