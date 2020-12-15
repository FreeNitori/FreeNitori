// Integrated command handlers.
package multiplexer

// Define the categories here
var AudioCategory = NewCategory("Audio",
	"Audio related utilities.")
var ExperienceCategory = NewCategory("Experience",
	"Chat experience and ranking system.")
var ManualsCategory = NewCategory("Manuals",
	"The operation manual pager utility.")
var ModerationCategory = NewCategory("Moderation",
	"Chat moderation utilities.")
var SystemCategory = NewCategory("System",
	"System-related utilities.")

var Categories = []*CommandCategory{ExperienceCategory, ManualsCategory, ModerationCategory, AudioCategory, SystemCategory}
