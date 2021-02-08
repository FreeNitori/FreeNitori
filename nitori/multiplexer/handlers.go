package multiplexer

// Predefined command categories.
var (
	AudioCategory = NewCategory("Audio",
		"Audio related utilities.")
	ExperienceCategory = NewCategory("Experience",
		"Chat experience and ranking system.")
	ManualsCategory = NewCategory("Manuals",
		"The operation manual pager utility.")
	MediaCategory = NewCategory("Media",
		"Media related utilities.")
	ModerationCategory = NewCategory("Moderation",
		"Chat moderation utilities.")
	SystemCategory = NewCategory("System",
		"System-related utilities.")
)

// Categories is a slice of pointers to command categories.
var Categories = []*CommandCategory{AudioCategory, ExperienceCategory, ManualsCategory, MediaCategory, ModerationCategory, SystemCategory}
