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

// Define the categories here
var SystemCategory = multiplexer.NewCategory("System",
	"System-related utilities.")
var ManualsCategory = multiplexer.NewCategory("Manuals",
	"The operation manual pager utility.")

var Categories = []*multiplexer.CommandCategory{SystemCategory, ManualsCategory}

// Define all the handlers here
var Handler Handlers
var AllHandlers = []HandlerMetadata{
	{"man",
		[]string{"help", "?"},
		"An interface to the on-line reference manuals.",
		ManualsCategory,
		Handler.Manuals},
	{"about",
		[]string{"info", "kappa", "information"},
		"Show some information about the kappa.",
		SystemCategory,
		Handler.About},
	{"reboot",
		[]string{"restart", "halt", "shutdown"},
		"Reboot the chat backend.",
		SystemCategory,
		Handler.Reboot},
	{"configure",
		[]string{"config", "conf", "settings"},
		"Configure per-guild overrides..",
		SystemCategory,
		Handler.Configure},
}

// Static messages
var InvalidArgument = "Invalid argument."
var AdminOnly = "This command is only available to system administrators!"
var KappaColor = 0x3492c4
