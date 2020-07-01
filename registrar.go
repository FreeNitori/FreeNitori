package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
)

var Router = multiplexer.New()

func init() {
	// Registers the router's handler to handle all incoming messages.
	Session.AddHandler(Router.OnMessageCreate)

	// Register commands
	Router.Route(
		"man",
		"An interface to the on-line reference manuals.",
		handlers.Handler.Manuals,
		handlers.ManualsCategory)
}
