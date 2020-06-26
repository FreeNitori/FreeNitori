package main

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

var Router = multiplexer.New()

func init() {
	// Registers the router's handler to handle all incoming messages.
	Session.AddHandler(Router.OnMessageCreate)

	// Registers the help command.
	// Router.Route("help", "Displays the help message.", Router.Help)
}
