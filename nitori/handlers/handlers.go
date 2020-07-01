package handlers

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

// Some structures to save some registering work
type Handlers struct{}
type HandlerMetadata struct {
	Pattern       string
	AliasPatterns []string
	Description   string
	Category      *multiplexer.CommandCategory
	Handler       multiplexer.CommandHandler
}

// Define all the handlers here
var Handler Handlers
var AllHandlers = []HandlerMetadata{
	{"man",
		[]string{"help", "?"},
		"An interface to the on-line reference manuals.",
		multiplexer.ManualsCategory,
		Handler.Manuals},
	{"about",
		[]string{"info", "kappa", "information"},
		"Show some information about the kappa.",
		multiplexer.SystemCategory,
		Handler.About},
}

// Static messages
var InvalidArgument = "Invalid argument."
var KappaColor = 0x3492c4
