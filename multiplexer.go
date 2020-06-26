package main

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

var Router = multiplexer.New()

func init() {
	// Registers the router's handler to handle all incoming messages.
	Session.AddHandler(Router.OnMessageCreate)

	// Registers the help command.
	_, _ = Router.Route("man", "An interface to the on-line reference manuals.", Router.Manuals)
}
