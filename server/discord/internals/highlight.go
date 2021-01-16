package internals

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/emoji"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"github.com/bwmarrin/discordgo"
	"strconv"
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
			{
				Name:         "amount",
				FriendlyName: "Minimum Requirement",
				Description:  "Minimum amount of reactions required for highlighting.",
				DatabaseKey:  "highlight_amount",
				Cleanup:      func(context *multiplexer.Context) { config.ResetGuildMap(context.Guild.ID, "highlight") },
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					amount, err := strconv.Atoi(*input)
					if err != nil {
						return false, true
					}
					if amount > 16 || amount < 1 {
						return false, true
					}
					return true, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if value == "" {
						return "Not configured", fmt.Sprintf(fmt.Sprintf("Configure by issuing command `%sconf highlight amount <integer>`.", context.Prefix())), true
					}
					return "Current requirement", value + " reactions", true
				},
			},
		},
	})
}

func handleHighlightReaction(session *discordgo.Session, reaction *discordgo.MessageReaction) {
	channelID, err := config.GetGuildConfValue(reaction.GuildID, "highlight_channel")
	if err != nil {
		return
	}
	amountString, err := config.GetGuildConfValue(reaction.GuildID, "highlight_amount")
	if err != nil {
		return
	}
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		return
	}
	guild, err := session.State.Guild(reaction.GuildID)
	if err != nil {
		guild, err = session.Guild(reaction.GuildID)
		if err != nil {
			return
		}
		_ = session.State.GuildAdd(guild)
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

	if reaction.Emoji.ID != "" {
		return
	}

	e, err := config.GetGuildConfValue(guild.ID, "highlight_emoji")
	if err != nil {
		return
	}
	if reaction.Emoji.Name != e {
		return
	}

	message, err := session.State.Message(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		message, err = session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
		if err != nil {
			return
		}
		_ = session.State.MessageAdd(message)
	}

	if message.Author.ID == session.State.User.ID {
		return
	}

	for _, reactions := range message.Reactions {
		if reactions.Emoji.Name == e {
			if reactions.Count >= amount {
				binding, err := config.HighlightGetBinding(guild.ID, message.ID)
				if err != nil {
					return
				}

				content := fmt.Sprintf("**%d | **%s", reactions.Count, channel.Mention())
				embed := embedutil.NewEmbed("", message.Content)
				if len(message.Attachments) > 0 {
					embed.SetImage(message.Attachments[0].URL)
				}
				embed.SetAuthor(message.Author.Username+"#"+message.Author.Discriminator, message.Author.AvatarURL("128"))
				embed.SetFooter(fmt.Sprintf("Author: %s", message.Author.ID))
				embed.Color = vars.KappaColor
				embed.AddField("Original Message", fmt.Sprintf("[Redirect](https://discord.com/channels/%s/%s/%s)", guild.ID, message.ChannelID, message.ID), false)

				if binding == "" {
					highlight, err := session.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
						Content:         content,
						Embed:           embed.MessageEmbed,
						TTS:             false,
						Files:           nil,
						AllowedMentions: nil,
						File:            nil,
					})
					if err != nil {
						return
					}
					err = config.HighlightBindMessage(guild.ID, message.ID, highlight.ID)
					if err != nil {
						return
					}
					binding = message.ID
				}
				_, _ = session.ChannelMessageEditComplex(&discordgo.MessageEdit{
					Content:         &content,
					Embed:           embed.MessageEmbed,
					AllowedMentions: nil,
					ID:              binding,
					Channel:         channel.ID,
				})
			}
			break
		}
	}
}

func addReaction(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	handleHighlightReaction(session, add.MessageReaction)
}

func removeReaction(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	handleHighlightReaction(session, remove.MessageReaction)
}
