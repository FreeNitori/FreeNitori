package handlers

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

// Some structures to save some registering work
type CommandHandlers struct{}
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
var ExperienceCategory = multiplexer.NewCategory("Experience",
	"Chat experience and ranking system.")

var Categories = []*multiplexer.CommandCategory{SystemCategory, ManualsCategory}

// Define all the handlers here
var CommandHandler CommandHandlers
var AllHandlers = []HandlerMetadata{
	{"man",
		[]string{"help", "?"},
		"An interface to the on-line reference manuals.",
		ManualsCategory,
		CommandHandler.Manuals},
	{"about",
		[]string{"info", "kappa", "information"},
		"Show some information about the kappa.",
		SystemCategory,
		CommandHandler.About},
	{"reboot",
		[]string{"restart", "halt", "shutdown"},
		"Reboot the chat backend.",
		SystemCategory,
		CommandHandler.Reboot},
	{"configure",
		[]string{"config", "conf", "settings"},
		"Configure per-guild overrides..",
		SystemCategory,
		CommandHandler.Configure},
	{"level",
		[]string{"rank", "exp"},
		"Query your current experience level.",
		ExperienceCategory,
		CommandHandler.Level},
}
