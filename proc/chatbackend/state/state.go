package state

import (
	"github.com/bwmarrin/discordgo"
	"github.com/shkh/lastfm-go/lastfm"
)

// Static messages
const InvalidArgument = "Invalid argument."
const ErrorOccurred = "Something went wrong and I am very confused! Please try again!"
const GuildOnly = "This command can only be issued from a guild."
const FeatureDisabled = "This feature is currently disabled."
const AdminOnly = "This command is only available to system administrators!"
const OperatorOnly = "This command is only available to operators!"
const PermissionDenied = "You are not allowed to issue this command!"
const MissingUser = "Specified user is not present."
const KappaColor = 0x3492c4

// Important users
var Administrator *discordgo.User
var Operator []*discordgo.User

// Session information
var RawSession, _ = discordgo.New()
var ShardSessions []*discordgo.Session
var LastFM *lastfm.Api
var Application *discordgo.Application
var EventHandlers []interface{}
