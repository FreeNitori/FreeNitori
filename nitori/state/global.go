// Variables containing important information.
package state

import "runtime/debug"

// Version information
var Version = func() string {
	build, ok := debug.ReadBuildInfo()
	if !ok {
		panic("Failed to read build info.")
	}
	return build.Main.Version
}()

// Channels
var (
	InviteURL    string
	ExitCode     = make(chan int)
	DiscordReady = make(chan bool)
)
