package handlers

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

// Define the categories here
var SystemCategory = multiplexer.NewCategory("System",
	"System-related utilities.")
var ManualsCategory = multiplexer.NewCategory("Manuals",
	"The operation manual pager utility.")
var ExperienceCategory = multiplexer.NewCategory("Experience",
	"Chat experience and ranking system.")

var Categories = []*multiplexer.CommandCategory{SystemCategory, ManualsCategory, ExperienceCategory}
