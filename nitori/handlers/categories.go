package handlers

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

// Define the categories here
var ExperienceCategory = multiplexer.NewCategory("Experience",
	"Chat experience and ranking system.")
var ManualsCategory = multiplexer.NewCategory("Manuals",
	"The operation manual pager utility.")
var ModerationCategory = multiplexer.NewCategory("Moderation",
	"Chat moderation utilities.")
var SystemCategory = multiplexer.NewCategory("System",
	"System-related utilities.")

var Categories = []*multiplexer.CommandCategory{ExperienceCategory, ManualsCategory, ModerationCategory, SystemCategory}
