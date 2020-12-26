// Variables containing important information.
package state

// Information
var (
	version   = "unknown"
	revision  = "unknown"
	InviteURL string
)

func Version() string  { return version }
func Revision() string { return revision }

// Channels
var (
	ExitCode     = make(chan int)
	DiscordReady = make(chan bool)
)
