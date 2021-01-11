package routes

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/overrides"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"github.com/bwmarrin/discordgo"
	"unicode"
)

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
	overrides.RegisterSimpleEntry(overrides.SimpleConfigurationEntry{
		Name:         "prefix",
		FriendlyName: "Command Prefix",
		Description:  "Configure command prefix.",
		DatabaseKey:  "prefix",
		Cleanup:      func(context *multiplexer.Context) {},
		Validate: func(context *multiplexer.Context, input *string) (bool, bool) {

			// Does not exceed length of 16
			if len(*input) > 16 {
				return false, true
			}

			// Add a space if last character is a letter
			if unicode.IsLetter([]rune((*input)[len(*input)-1:])[0]) {
				*input += " "
			}

			return true, true
		},
		Format: func(context *multiplexer.Context, value string) (string, string, bool) {
			if value == "" {
				return "Current prefix", config.Config.System.Prefix, true
			}
			return "Current prefix", value, true
		},
	})
	overrides.RegisterCustomEntry(overrides.CustomConfigurationEntry{
		Name:        "message",
		Description: "Configure custom messages.",
		Handler: func(context *multiplexer.Context) {
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
		},
	})
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
	switch length := len(context.Fields); length {
	case 1:
		embed := embedutil.NewEmbed("Configurator", "Configure per-guild overrides.")
		embed.Color = vars.KappaColor
		for _, entry := range overrides.GetSimpleEntries() {
			embed.AddField(entry.Name, entry.Description, false)
		}
		for _, entry := range overrides.GetComplexEntries() {
			embed.AddField(entry.Name, entry.Description, false)
		}
		for _, entry := range overrides.GetCustomEntries() {
			embed.AddField(entry.Name, entry.Description, false)
		}
		context.SendEmbed(embed)
	case 2:
		for _, entry := range overrides.GetSimpleEntries() {
			if context.Fields[1] == entry.Name {
				embed := embedutil.NewEmbed(entry.FriendlyName, entry.Description)
				embed.Color = vars.KappaColor
				value, err := config.GetGuildConfValue(context.Guild.ID, entry.DatabaseKey)
				if !context.HandleError(err) {
					return
				}
				title, description, ok := entry.Format(context, value)
				if !ok {
					return
				}
				embed.AddField(title, description, true)
				context.SendEmbed(embed)
				return
			}
		}
		for _, entry := range overrides.GetComplexEntries() {
			if context.Fields[1] == entry.Name {
				embed := embedutil.NewEmbed(entry.FriendlyName, entry.Description)
				embed.Color = vars.KappaColor
				for _, subEntry := range entry.Entries {
					embed.AddField(subEntry.Name, subEntry.Description, false)
				}
				context.SendEmbed(embed)
				return
			}
		}
		fallthrough
	default:
		if length > 1 {
			for _, entry := range overrides.GetCustomEntries() {
				if context.Fields[1] == entry.Name {
					entry.Handler(context)
					return
				}
			}
		}

		if length < 3 {
			context.SendMessage(vars.InvalidArgument)
			return
		}

		for _, entry := range overrides.GetSimpleEntries() {
			if context.Fields[1] == entry.Name {
				if context.Fields[2] == "reset" {
					err := config.ResetGuildConfValue(context.Guild.ID, entry.DatabaseKey)
					if !context.HandleError(err) {
						return
					}
					entry.Cleanup(context)
					context.SendMessage(fmt.Sprintf("Successfully reset value of `%s`.", entry.DatabaseKey))
					return
				}
				input := context.StitchFields(2)
				valid, ok := entry.Validate(context, &input)
				if !ok {
					return
				}
				if !valid {
					context.SendMessage(vars.InvalidArgument)
					return
				}
				err := config.SetGuildConfValue(context.Guild.ID, entry.DatabaseKey, input)
				if !context.HandleError(err) {
					return
				}
				context.SendMessage(fmt.Sprintf("Successfully set value of `%s` to `%s`.", entry.Name, input))
				return
			}
		}

		for _, entry := range overrides.GetComplexEntries() {
			if context.Fields[1] == entry.Name {
				for _, subEntry := range entry.Entries {
					if context.Fields[2] == subEntry.Name {
						if len(context.Fields) == 3 {
							embed := embedutil.NewEmbed(subEntry.FriendlyName, subEntry.Description)
							embed.Color = vars.KappaColor
							value, err := config.GetGuildConfValue(context.Guild.ID, subEntry.DatabaseKey)
							if !context.HandleError(err) {
								return
							}
							title, description, ok := subEntry.Format(context, value)
							if !ok {
								return
							}
							embed.AddField(title, description, true)
							context.SendEmbed(embed)
							return
						} else {
							if context.Fields[3] == "reset" {
								err := config.ResetGuildConfValue(context.Guild.ID, subEntry.DatabaseKey)
								if !context.HandleError(err) {
									return
								}
								subEntry.Cleanup(context)
								context.SendMessage(fmt.Sprintf("Successfully reset value of `%s`.", subEntry.DatabaseKey))
								return
							}
							input := context.StitchFields(3)
							valid, ok := subEntry.Validate(context, &input)
							if !ok {
								return
							}
							if !valid {
								context.SendMessage(vars.InvalidArgument)
								return
							}
							err := config.SetGuildConfValue(context.Guild.ID, subEntry.DatabaseKey, input)
							if !context.HandleError(err) {
								return
							}
							context.SendMessage(fmt.Sprintf("Successfully set value of `%s.%s` to `%s`.", entry.Name, subEntry.Name, input))
							return
						}
					}
				}
				break
			}
		}
		context.SendMessage(vars.InvalidArgument)
	}
}
