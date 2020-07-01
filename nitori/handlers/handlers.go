package handlers

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

// Some structures to save some registering work
type Handlers struct{}
type HandlerMetadata struct {
	Pattern     string
	Description string
	Category    *multiplexer.CommandCategory
	Handler     multiplexer.CommandHandler
}

// Define all the handlers here
var Handler Handlers
var AllHandlers = []HandlerMetadata{
	{"man",
		"An interface to the on-line reference manuals.",
		multiplexer.ManualsCategory,
		Handler.Manuals},
}

// Static messages
var InvalidArgument = "Invalid argument."
