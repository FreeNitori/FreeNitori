package internals

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/sessioning"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func init() {
	multiplexer.GuildMemberAdd = append(multiplexer.GuildMemberAdd, welcomeHandler)
	multiplexer.GuildMemberRemove = append(multiplexer.GuildMemberRemove, removeHandler)
	overrides.RegisterComplexEntry(overrides.ComplexConfigurationEntry{
		Name:         "greet",
		FriendlyName: "Greeter",
		Description:  "Configure member greeter.",
		Entries: []overrides.SimpleConfigurationEntry{
			{
				Name:         "channel",
				FriendlyName: "Greeting Channel",
				Description:  "Channel to send all greeting messages to.",
				DatabaseKey:  "greet_channel",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if channel := context.GetChannel(*input); channel != nil {
						*input = channel.ID
						return true, true
					}
					return false, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if channel := context.GetChannel(value); channel != nil {
						return channel.Name, channel.ID, true
					}
					return "No channel configured", fmt.Sprintf("Configure it by issuing command `%sconf greet channel <channel>`.", context.Prefix()), true
				},
			},
			{
				Name:         "welcome-message",
				FriendlyName: "Welcome Message",
				Description:  "Message sent on user join, setting this will enable welcomes.",
				DatabaseKey:  "welcome_message",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if len(*input) > 2000 {
						return false, true
					}
					return true, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if value == "" {
						return "No welcome message is set", fmt.Sprintf(
							"Configure it by issuing command `%sconf greet welcome-message <message>`.\n"+
								"Valid placeholders: \n"+
								"$USERNAME: Username of the user.\n"+
								"$DISCRIMINATOR: Discriminator/tag of the user.\n"+
								"$MENTION: Mention the user.", context.Prefix()), true
					}
					return "Current message", value, true
				},
			},
			{
				Name:         "welcome-url",
				FriendlyName: "Welcome Image URL",
				Description:  "Image embedded in message sent on user join.",
				DatabaseKey:  "welcome_url",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if len(*input) > 2000 {
						return false, true
					}
					return true, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if value == "" {
						return "No welcome image URL is set", fmt.Sprintf("Configure it by issuing command `%sconf greet welcome-url <url>`.", context.Prefix()), true
					}
					return "Current URL", value, true
				},
			},
			{
				Name:         "goodbye-message",
				FriendlyName: "Goodbye Message",
				Description:  "Message sent on user leave, setting this will enable goodbyes.",
				DatabaseKey:  "goodbye_message",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if len(*input) > 2000 {
						return false, true
					}
					return true, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if value == "" {
						return "No goodbye message is set", fmt.Sprintf(
							"Configure it by issuing command `%sconf greet goodbye-message <message>`.\n"+
								"Valid placeholders: \n"+
								"$USERNAME: Username of the user.\n"+
								"$DISCRIMINATOR: Discriminator/tag of the user.\n"+
								"$MENTION: Mention the user.", context.Prefix()), true
					}
					return "Current message", value, true
				},
			},
			{
				Name:         "goodbye-url",
				FriendlyName: "Welcome Image URL",
				Description:  "Image embedded in message sent on user leave.",
				DatabaseKey:  "goodbye_url",
				Cleanup:      func(context *multiplexer.Context) {},
				Validate: func(context *multiplexer.Context, input *string) (bool, bool) {
					if len(*input) > 2000 {
						return false, true
					}
					return true, true
				},
				Format: func(context *multiplexer.Context, value string) (string, string, bool) {
					if value == "" {
						return "No goodbye image URL is set", fmt.Sprintf("Configure it by issuing command `%sconf greet goodbye-url <url>`.", context.Prefix()), true
					}
					return "Current URL", value, true
				},
			},
		},
		CustomEntries: nil,
	})
}

func welcomeHandler(session *discordgo.Session, add *discordgo.GuildMemberAdd) {
	if add.Member.User.Bot {
		return
	}
	var embed embedutil.Embed
	channelID, err := config.GetGuildConfValue(add.GuildID, "greet_channel")
	if err != nil {
		return
	}
	if channelID == "" {
		return
	}
	if sessioning.FetchChannel(sessioning.FetchGuild(add.GuildID), "", channelID) == nil {
		return
	}
	message, err := config.GetGuildConfValue(add.GuildID, "welcome_message")
	if err != nil {
		return
	}
	if message == "" {
		return
	}
	url, err := config.GetGuildConfValue(add.GuildID, "welcome_url")
	if err != nil {
		return
	}
	if url != "" {
		embed = embedutil.New("", "")
		embed.Color = state.KappaColor
		embed.SetImage(url)
	}
	context := &multiplexer.Context{Message: &discordgo.Message{ChannelID: channelID}, Session: session}
	context.SendEmbed(strings.NewReplacer(
		"$USERNAME", add.User.Username,
		"$DISCRIMINATOR", add.User.Discriminator,
		"$MENTION", add.User.Mention()).Replace(message), embed)
}

func removeHandler(session *discordgo.Session, remove *discordgo.GuildMemberRemove) {
	if remove.Member.User.Bot {
		return
	}
	var embed embedutil.Embed
	channelID, err := config.GetGuildConfValue(remove.GuildID, "greet_channel")
	if err != nil {
		return
	}
	if channelID == "" {
		return
	}
	if sessioning.FetchChannel(sessioning.FetchGuild(remove.GuildID), "", channelID) == nil {
		return
	}
	message, err := config.GetGuildConfValue(remove.GuildID, "goodbye_message")
	if err != nil {
		return
	}
	if message == "" {
		return
	}
	url, err := config.GetGuildConfValue(remove.GuildID, "goodbye_url")
	if err != nil {
		return
	}
	if url != "" {
		embed = embedutil.New("", "")
		embed.Color = state.KappaColor
		embed.SetImage(url)
	}
	context := &multiplexer.Context{Message: &discordgo.Message{ChannelID: channelID}, Session: session}
	context.SendEmbed(strings.NewReplacer(
		"$USERNAME", remove.User.Username,
		"$DISCRIMINATOR", remove.User.Discriminator,
		"$MENTION", remove.User.Mention()).Replace(message), embed)
}
