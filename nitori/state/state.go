// Package state contains important variables.
package state

import (
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"time"
)

// Information
var (
	version  = "unknown"
	revision = "unknown"
	start    time.Time
)

// InviteURL contains the invite URL of the bot.
var InviteURL string

// Version returns the version of Nitori.
func Version() string { return version }

// Revision returns the git revision of Nitori.
func Revision() string { return revision }

// Uptime returns the duration the instance has been online.
func Uptime() time.Duration { return time.Since(start) }

// Channels
var (
	ExitCode     = make(chan int)
	DiscordReady = make(chan bool)
)

// Multiplexer is the Discord event multiplexer.
var Multiplexer = multiplexer.New()

// RawSession is the raw session with Discord.
var RawSession, _ = discordgo.New()

// ShardSessions is the slice of sessions of each shard.
var ShardSessions []*discordgo.Session

// Application is the Discord application of the instance.
var Application *discordgo.Application

func init() {
	start = time.Now()
}
