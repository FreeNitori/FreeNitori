// Variables containing important information.
package state

// Information
var (
	Version   = "unknown"
	Revision  = "unknown"
	InviteURL string
)

// Channels
var (
	ExitCode     = make(chan int)
	DiscordReady = make(chan bool)
)
