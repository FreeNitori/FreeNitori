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
var ExperienceCategory = multiplexer.NewCategory("Experience",
	"Chat experience and ranking system.")

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
	{"level",
		[]string{"rank", "exp"},
		"Query your current experience level.",
		ExperienceCategory,
		Handler.Level},
}

// Static messages
var InvalidArgument = "Invalid argument."
var ErrorOccurred = "An error occurred while handling your request, please try again later!"
var GuildOnly = "This command can only be issued from a guild."
var FeatureDisabled = "This feature is currently disabled."
var AdminOnly = "This command is only available to system administrators!"
var KappaColor = 0x3492c4
