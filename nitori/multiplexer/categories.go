package multiplexer

var SystemCategory = NewCategory("System",
	"System-related utilities.")
var ManualsCategory = NewCategory("Manuals",
	"The operation manual pager utility.")

var Categories = []*CommandCategory{SystemCategory, ManualsCategory}
