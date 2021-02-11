package multiplexer

import "github.com/bwmarrin/discordgo"

// Router is the command router instance.
var Router = New()

// Slice of event handlers.
var (
	EventHandlers         []interface{}
	NotTargeted           []func(context *Context)
	Ready                 []func(session *discordgo.Session, ready *discordgo.Ready)
	GuildMemberAdd        []func(session *discordgo.Session, add *discordgo.GuildMemberAdd)
	GuildMemberRemove     []func(session *discordgo.Session, remove *discordgo.GuildMemberRemove)
	GuildDelete           []func(session *discordgo.Session, delete *discordgo.GuildDelete)
	MessageCreate         []func(session *discordgo.Session, create *discordgo.MessageCreate)
	MessageDelete         []func(session *discordgo.Session, delete *discordgo.MessageDelete)
	MessageUpdate         []func(session *discordgo.Session, update *discordgo.MessageUpdate)
	MessageReactionAdd    []func(session *discordgo.Session, add *discordgo.MessageReactionAdd)
	MessageReactionRemove []func(session *discordgo.Session, remove *discordgo.MessageReactionRemove)
)
