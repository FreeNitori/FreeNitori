package multiplexer

import "github.com/bwmarrin/discordgo"

var (
	EventHandlers         []interface{}
	Router                = New()
	NotTargeted           []func(context *Context)
	Ready                 []func(session *discordgo.Session, ready *discordgo.Ready)
	GuildMemberAdd        []func(session *discordgo.Session, add *discordgo.GuildMemberAdd)
	GuildMemberRemove     []func(session *discordgo.Session, remove *discordgo.GuildMemberRemove)
	GuildDelete           []func(session *discordgo.Session, delete *discordgo.GuildDelete)
	MessageReactionAdd    []func(session *discordgo.Session, add *discordgo.MessageReactionAdd)
	MessageReactionRemove []func(session *discordgo.Session, remove *discordgo.MessageReactionRemove)
)
