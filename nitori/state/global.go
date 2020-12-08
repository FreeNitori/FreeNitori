// Variables containing important information.
package state

// Version information
const Version = "v0.1.0"

// Channels
var (
	InviteURL string
	ExitCode = make(chan int)
	DiscordReady = make(chan bool)
)