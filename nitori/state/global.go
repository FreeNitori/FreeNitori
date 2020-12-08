// Variables containing important information.
package state

// Version information
const Version = "v0.0.1-rewrite"

// Channels
var (
	InviteURL string
	ExitCode = make(chan int)
	DiscordReady = make(chan bool)
)