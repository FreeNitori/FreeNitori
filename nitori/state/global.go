// Variables containing important information.
package state

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

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

// Static messages
const InvalidArgument = "Invalid argument."
const ErrorOccurred = "Something went wrong and I am very confused! Please try again!"
const GuildOnly = "This command can only be issued from a guild."
const FeatureDisabled = "This feature is currently disabled."
const AdminOnly = "This command is only available to system administrators!"
const OperatorOnly = "This command is only available to operators!"
const PermissionDenied = "You are not allowed to issue this command!"
const MissingUser = "Specified user does not exist."
const LackingPermission = "Lacking permission to perform specified action."
const KappaColor = 0x3492c4

// Important users
var Administrator *discordgo.User
var Operator []*discordgo.User

// Session information
var RawSession, _ = discordgo.New()
var ShardSessions []*discordgo.Session
var Application *discordgo.Application

func init() {
	start = time.Now()
}
