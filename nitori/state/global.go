// Variables containing important information.
package state

// Version information
var Version = "unknown"
var Revision = "unknown"

// Channels
var (
	InviteURL    string
	ExitCode     = make(chan int)
	DiscordReady = make(chan bool)
)
