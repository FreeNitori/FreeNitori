// Variables containing important information.
package state

import "time"

// Information
var (
	version   = "unknown"
	revision  = "unknown"
	start     time.Time
	InviteURL string
)

func Version() string       { return version }
func Revision() string      { return revision }
func Uptime() time.Duration { return time.Since(start) }

// Channels
var (
	ExitCode     = make(chan int)
	DiscordReady = make(chan bool)
)

func init() {
	start = time.Now()
}
