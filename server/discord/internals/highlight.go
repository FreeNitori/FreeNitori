package internals

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"github.com/bwmarrin/discordgo"
)

func init() {
	multiplexer.MessageReactionAdd = append(multiplexer.MessageReactionAdd, addReaction)
	multiplexer.MessageReactionRemove = append(multiplexer.MessageReactionRemove, removeReaction)
	overrides.RegisterComplexEntry(overrides.ComplexConfigurationEntry{
		Name:         "highlight",
		FriendlyName: "Message Highlighting",
		Description:  "Configure message highlighting system.",
		Entries: []overrides.SimpleConfigurationEntry{
			{
				Name:         "channel",
				FriendlyName: "Highlighted Message Channel",
				Description:  "Channel highlighted messages are posted to.",
				DatabaseKey:  "highlight_channel",
				Cleanup: func(context *multiplexer.Context) {
					// TODO: delete all message references
				},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					for _, channel := range context.Guild.Channels {
						if *input == channel.ID {
							return true, true
						}
					}
					return false, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					for _, channel := range context.Guild.Channels {
						if value == channel.ID {
							return channel.Name, channel.ID, true
						}
					}
					return "No channel configured", fmt.Sprintf("Configure it by issuing command `%sconf highlight channel <channelID>`.", context.Prefix()), true
				},
			},
			{
				Name:         "emoji",
				FriendlyName: "Emoji",
				Description:  "Emoji used to vote the message.",
				DatabaseKey:  "highlight_emoji",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					// TODO: check default emoji
					return false, false
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					// TODO: format stuff
					return "", "", false
				},
			},
		},
	})
}

func addReaction(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	channelID, err := config.GetGuildConfValue(add.GuildID, "highlight_channel")
	if err != nil {
		return
	}
	guild, err := session.State.Guild(add.GuildID)
	if err != nil {
		return
	}
	var channel *discordgo.Channel
	for _, c := range guild.Channels {
		if channelID == c.ID {
			channel = c
			break
		}
	}
	if channel == nil {
		return
	}
	// TODO: check for emote
}

func removeReaction(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	channelID, err := config.GetGuildConfValue(remove.GuildID, "highlight_channel")
	if err != nil {
		return
	}
	guild, err := session.State.Guild(remove.GuildID)
	if err != nil {
		return
	}
	var channel *discordgo.Channel
	for _, c := range guild.Channels {
		if channelID == c.ID {
			channel = c
			break
		}
	}
	if channel == nil {
		return
	}
	// TODO: check for emote
}
