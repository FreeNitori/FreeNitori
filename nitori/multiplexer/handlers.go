// Integrated command handlers.
package multiplexer

// Define the categories here
var (
	AudioCategory = NewCategory("Audio",
		"Audio related utilities.")
	ExperienceCategory = NewCategory("Experience",
		"Chat experience and ranking system.")
	ManualsCategory = NewCategory("Manuals",
		"The operation manual pager utility.")
	ModerationCategory = NewCategory("Moderation",
		"Chat moderation utilities.")
	SystemCategory = NewCategory("System",
		"System-related utilities.")
)

var Categories = []*CommandCategory{AudioCategory, ExperienceCategory, ManualsCategory, ModerationCategory, SystemCategory}
