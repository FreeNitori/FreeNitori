package internals

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/emoji"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"github.com/bwmarrin/discordgo"
	"unicode/utf8"
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
				Cleanup:      func(context *multiplexer.Context) { config.ResetGuildMap(context.Guild.ID, "highlight") },
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					for _, channel := range context.Guild.Channels {
						if *input == channel.ID {
							config.ResetGuildMap(context.Guild.ID, "highlight")
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
				Cleanup:      func(context *multiplexer.Context) { config.ResetGuildMap(context.Guild.ID, "highlight") },
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if utf8.RuneCountInString(*input) != 1 {
						return false, true
					}
					var key string
					for _, r := range []rune(*input) {
						key += fmt.Sprintf("%X", r)
					}
					_, ok := emoji.Emojis[key]
					if ok {
						config.ResetGuildMap(context.Guild.ID, "highlight")
					}
					return ok, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if value == "" {
						return "No emoji is configured", fmt.Sprintf("Configure it by issuing command `%sconf highlight emoji <emoji>`.", context.Prefix()), true
					}
					return "Current emoji", value, true
				},
			},
		},
	})
}

func processReaction(session *discordgo.Session, reaction *discordgo.MessageReaction) {
	channelID, err := config.GetGuildConfValue(reaction.GuildID, "highlight_channel")
	if err != nil {
		return
	}
	guild, err := session.State.Guild(reaction.GuildID)
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
	if reaction.Emoji.ID != "null" {
		return
	}
	e, err := config.GetGuildConfValue(guild.ID, "highlight_emoji")
	if err != nil {
		return
	}
	if reaction.Emoji.Name != e {
		return
	}
	log.Info(reaction.Emoji.Name)
}

func addReaction(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	processReaction(session, interface{}(add).(*discordgo.MessageReaction))
	// TODO: handler thing
}

func removeReaction(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	processReaction(session, interface{}(remove).(*discordgo.MessageReaction))
	// TODO: handler thing
}
