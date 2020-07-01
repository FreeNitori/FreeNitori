package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
)

var Router = multiplexer.New()

func init() {

	// Registers the router's handler to handle all incoming messages.
	Session.AddHandler(Router.OnMessageCreate)

	// Register all route handlers
	for _, handlerMeta := range handlers.AllHandlers {
		Router.Route(
			handlerMeta.Pattern,
			handlerMeta.Description,
			handlerMeta.Handler,
			handlerMeta.Category)
	}
}
