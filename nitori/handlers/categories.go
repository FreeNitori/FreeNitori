package handlers

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

var ManualsCategory = multiplexer.NewCategory("Manuals",
	"The operation manual pager utility.")

var Categories = []*multiplexer.CommandCategory{ManualsCategory}
