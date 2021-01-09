package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"unicode"
)

var err error

func init() {
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "configure",
		AliasPatterns: []string{"conf", "settings", "set"},
		Description:   "Configure per-guild overrides.",
		Category:      multiplexer.SystemCategory,
		Handler:       configure,
	})
	multiplexer.GuildMemberRemove = append(multiplexer.GuildMemberRemove, func(session *discordgo.Session, remove *discordgo.GuildMemberRemove) {
		if remove.User.ID == session.State.User.ID {
			config.ResetGuild(remove.GuildID)
		}
	})
	multiplexer.GuildDelete = append(multiplexer.GuildDelete, func(session *discordgo.Session, delete *discordgo.GuildDelete) {
		config.ResetGuild(delete.ID)
	})
}

var SimpleEntries []SimpleConfigurationEntry
var ComplexEntries []ComplexConfigurationEntry

type SimpleConfigurationEntry struct {
	Name        string
	Description string
	DatabaseKey string
	Validator   func(context *multiplexer.Context, input string) bool
	Formatter   func(context *multiplexer.Context, value string) (string, string)
}

type ComplexConfigurationEntry struct {
	Name        string
	Description string
	Entries     []SimpleConfigurationEntry
}

func configure(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(vars.GuildOnly)
		return
	}
	if !context.HasPermission(discordgo.PermissionAdministrator) {
		context.SendMessage(vars.PermissionDenied)
		return
	}
	switch len(context.Fields) {
	case 1:
		embed := embedutil.NewEmbed("Configurator", "Configure per-guild overrides.")
		embed.Color = vars.KappaColor
		for _, entry := range SimpleEntries {
			embed.AddField(entry.Name, entry.Description, false)
		}
		for _, entry := range ComplexEntries {
			embed.AddField(entry.Name, entry.Description, false)
		}
		context.SendEmbed(embed)
	case 2:
		for _, entry := range SimpleEntries {
			if context.Fields[1] == entry.Name {
				embed := embedutil.NewEmbed(entry.Name, entry.Description)
				embed.Color = vars.KappaColor
				value, err := config.GetGuildConfValue(context.Guild.ID, entry.DatabaseKey)
				if !context.HandleError(err) {
					return
				}
				title, description := entry.Formatter(context, value)
				embed.AddField(title, description, true)
				context.SendEmbed(embed)
				return
			}
		}
	default:
		context.SendMessage(vars.InvalidArgument)
	}
}

func configureOld(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(vars.GuildOnly)
		return
	}
	if !context.HasPermission(discordgo.PermissionAdministrator) {
		context.SendMessage(vars.PermissionDenied)
		return
	}
	if len(context.Fields) == 1 {
		embed := embedutil.NewEmbed("Configurator", "Configure per-guild overrides.")
		embed.Color = vars.KappaColor
		embed.AddField("experience", "Toggle all experience system.", false)
		embed.AddField("highlight", "Configure message highlighting system.", false)
		embed.AddField("message", "Configure customizable messages.", false)
		embed.AddField("prefix", "Configure command prefix.", false)
		context.SendEmbed(embed)
		return
	}
	switch context.Fields[1] {
	case "highlight":
		switch len(context.Fields) {
		case 3:
			switch context.Fields[2] {
			case "channel":
				embed := embedutil.NewEmbed("Highlight Channel", "Configure channel for highlighted messages.")
				id, err := config.GetHighlightChannelID(context.Guild)
				if !context.HandleError(err) {
					return
				}
				if id == 0 {
					embed.AddField("No channel was configured", fmt.Sprintf("Configure a message by appending channel ID after this command."), false)
				} else {
					ok := false
					for _, channel := range context.Guild.Channels {
						if strconv.Itoa(id) == channel.ID {
							ok = true
							embed.AddField(channel.Name, channel.ID, false)
							break
						}
					}
					if !ok {
						if !context.HandleError(config.ResetHighlightChannelID(context.Guild)) {
							return
						}
						embed.AddField("No channel was configured", fmt.Sprintf("Configure a message by appending channel ID after this command."), false)
					}
				}
				context.SendEmbed(embed)
			case "emote":
				// TODO: also help message
			case "trigger":
				// TODO: even more help message
			}
		case 2:
			embed := embedutil.NewEmbed("Message highlighting", "Configure message highlighting related stuff.")
			embed.Color = vars.KappaColor
			embed.AddField("channel", "Configure channel for highlighted messages.", false)
			embed.AddField("emote", "Configure emote used for highlighting a message.", false)
			embed.AddField("trigger", "Configure amount of reactions to trigger highlighting.", false)
			context.SendEmbed(embed)
		}
	case "prefix":
		switch len(context.Fields) {
		case 3:
			var newPrefix = context.Fields[2]

			// Add a space if last character is a letter
			if unicode.IsLetter([]rune(newPrefix[len(newPrefix)-1:])[0]) {
				newPrefix += " "
			}

			// Actually set the prefix
			err = config.SetPrefix(context.Guild.ID, newPrefix)
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully updated prefix.")
		case 2:
			err = config.ResetPrefix(context.Guild.ID)
			if !context.HandleError(err) {
				return
			}
			context.SendMessage("Successfully reset prefix.")
		default:
			context.SendMessage(vars.InvalidArgument)
		}
	case "experience":
		pre, err := config.ExpToggle(context.Guild.ID)
		if !context.HandleError(err) {
			return
		}
		switch pre {
		case false:
			context.SendMessage("Chat experience system has been enabled.")
		case true:
			context.SendMessage("Chat experience system has been disabled.")
		}
	case "message":
		switch len(context.Fields) {
		default:
			err := config.SetCustomizableMessage(context.Guild.ID, context.Fields[2], context.StitchFields(3))
			switch err.(type) {
			default:
				if !context.HandleError(err) {
					return
				}
			case *config.MessageOutOfBounds:
				context.SendMessage(vars.InvalidArgument)
				return
			}
			context.SendMessage("Message `" + context.Fields[2] + "` has been set.")
		case 3:
			err := config.SetCustomizableMessage(context.Guild.ID, context.Fields[2], "")
			switch err.(type) {
			default:
				if !context.HandleError(err) {
					return
				}
			case *config.MessageOutOfBounds:
				context.SendMessage(vars.InvalidArgument)
				return
			}
			context.SendMessage("Message `" + context.Fields[2] + "` has been reset.")
		case 2:
			embed := embedutil.NewEmbed("Messages", "Configurable messages.")
			for identifier := range config.CustomizableMessages {
				message, err := config.GetCustomizableMessage(context.Guild.ID, identifier)
				if !context.HandleError(err) {
					return
				}
				embed.AddField(identifier, message, false)
			}
			context.SendEmbed(embed)
		}
	}
}
