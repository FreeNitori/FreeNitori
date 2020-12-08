// Integrated command handlers.
package handlers

import "git.randomchars.net/RandomChars/FreeNitori/server/discord/multiplexer"

// Define the categories here
var AudioCategory = multiplexer.NewCategory("Audio",
	"Audio related utilities.")
var ExperienceCategory = multiplexer.NewCategory("Experience",
	"Chat experience and ranking system.")
var ManualsCategory = multiplexer.NewCategory("Manuals",
	"The operation manual pager utility.")
var ModerationCategory = multiplexer.NewCategory("Moderation",
	"Chat moderation utilities.")
var SystemCategory = multiplexer.NewCategory("System",
	"System-related utilities.")

var Categories = []*multiplexer.CommandCategory{ExperienceCategory, ManualsCategory, ModerationCategory, AudioCategory, SystemCategory}
