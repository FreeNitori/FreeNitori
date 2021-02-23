package internals

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/db"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/emoji"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
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
			{
				Name:         "channel",
				FriendlyName: "Highlighted Message Channel",
				Description:  "Channel highlighted messages are posted to.",
				DatabaseKey:  "highlight_channel",
				Cleanup:      func(context *multiplexer.Context) { config.ResetGuildMap(context.Guild.ID, "highlight") },
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if channel := context.GetChannel(*input); channel != nil {
						*input = channel.ID
						config.ResetGuildMap(context.Guild.ID, "highlight")
						return true, true
					}
					return false, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if channel := context.GetChannel(value); channel != nil {
						return channel.Name, channel.ID, true
					}
					return "No channel configured", fmt.Sprintf("Configure it by issuing command `%sconf highlight channel <channel>`.", context.Prefix()), true
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
		CustomEntries: []overrides.CustomConfigurationEntry{
			{
				Name:        "inhibit",
				Description: "Inhibits highlighting of specified message.",
				Handler: func(context *multiplexer.Context) {
					if context.Message.MessageReference == nil {
						context.SendMessage("Please reply to the message to inhibit.")
						return
					}
					binding, err := db.HighlightGetBinding(context.Message.MessageReference.GuildID, context.Message.MessageReference.MessageID)
					if !context.HandleError(err) {
						return
					}
					if binding == "-" {
						err = db.HighlightUnbindMessage(context.Message.MessageReference.GuildID, context.Message.MessageReference.MessageID)
						if !context.HandleError(err) {
							return
						}
						context.SendMessage("Successfully uninhibited the message.")
						return
					}
					err = db.HighlightBindMessage(context.Message.MessageReference.GuildID, context.Message.MessageReference.MessageID, "-")
					if !context.HandleError(err) {
						return
					}
					context.SendMessage("Successfully inhibited the message.")
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

	message, err := session.ChannelMessage(reaction.ChannelID, reaction.MessageID)
	if err != nil {
		return
	}

	if message.Author.ID == session.State.User.ID {
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

	binding, err := db.HighlightGetBinding(guild.ID, message.ID)
	if err != nil {
		return
	}
	if binding == "-" {
		return
	}

	sufficient, reactions := func() (bool, *discordgo.MessageReactions) {
		for _, react := range message.Reactions {
			if react.Emoji.Name == e {
				if react.Count >= amount {
					return true, react
				}
				return false, nil
			}
		}
		return false, nil
	}()

	if sufficient {
		content := fmt.Sprintf("**%d | **%s", reactions.Count, fmt.Sprintf("<#%s>", message.ChannelID))
		embed := embedutil.New("", message.Content)
		for _, attachment := range message.Attachments {
			if attachment.Width != 0 && attachment.Height != 0 {
				embed.SetImage(attachment.URL, attachment.ProxyURL)
			}
			embed.AddField("Attachment", fmt.Sprintf("[%s](%s)", attachment.Filename, attachment.URL), false)
		}
		embed.SetAuthor(message.Author.Username+"#"+message.Author.Discriminator, message.Author.AvatarURL("128"))
		embed.SetFooter(fmt.Sprintf("Author: %s", message.Author.ID))
		embed.Color = state.KappaColor
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
			err = db.HighlightBindMessage(guild.ID, message.ID, highlight.ID)
			if err != nil {
				return
			}
		} else {
			_, err = session.ChannelMessageEditComplex(&discordgo.MessageEdit{
				Content:         &content,
				Embed:           embed.MessageEmbed,
				AllowedMentions: nil,
				ID:              binding,
				Channel:         channel.ID,
			})

			if fmt.Sprint(err) == "HTTP 404 Not Found, {\"message\": \"Unknown Message\", \"code\": 10008}" {
				err = db.HighlightUnbindMessage(guild.ID, message.ID)
				if err != nil {
					return
				}
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
				err = db.HighlightBindMessage(guild.ID, message.ID, highlight.ID)
				if err != nil {
					return
				}
			}
		}
	} else {
		if binding != "" {
			_ = session.ChannelMessageDelete(channelID, binding)
			err = db.HighlightUnbindMessage(guild.ID, binding)
			if err != nil {
				return
			}
		}
	}
}

func addReaction(session *discordgo.Session, add *discordgo.MessageReactionAdd) {
	handleHighlightReaction(session, add.MessageReaction)
}

func removeReaction(session *discordgo.Session, remove *discordgo.MessageReactionRemove) {
	handleHighlightReaction(session, remove.MessageReaction)
}
