package internals

import (
	"fmt"
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/overrides"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"unicode"
)

func init() {
	state.Multiplexer.Route(&multiplexer.Route{
		Pattern:       "configure",
		AliasPatterns: []string{"conf", "settings", "set"},
		Description:   "Configure per-guild overrides.",
		Category:      multiplexer.SystemCategory,
		Handler:       configure,
	})
	state.Multiplexer.GuildMemberRemove = append(state.Multiplexer.GuildMemberRemove, func(context *multiplexer.Context) {
		if context.User.ID == context.Session.State.User.ID {
			config.ResetGuild(context.Guild.ID)
		}
	})
	state.Multiplexer.GuildDelete = append(state.Multiplexer.GuildDelete, func(context *multiplexer.Context) {
		config.ResetGuild(context.Guild.ID)
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
					context.SendMessage(multiplexer.InvalidArgument)
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
					context.SendMessage(multiplexer.InvalidArgument)
					return
				}
				context.SendMessage("Message `" + context.Fields[2] + "` has been reset.")
			case 2:
				embed := embedutil.New("Messages", "Configurable messages.")
				for identifier := range config.CustomizableMessages {
					message, err := config.GetCustomizableMessage(context.Guild.ID, identifier)
					if !context.HandleError(err) {
						return
					}
					embed.AddField(identifier, message, false)
				}
				context.SendEmbed("", embed)
			}
		},
	})
}

func configure(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(multiplexer.GuildOnly)
		return
	}
	if !context.HasPermission(discordgo.PermissionAdministrator) {
		context.SendMessage(multiplexer.PermissionDenied)
		return
	}
	switch length := len(context.Fields); length {
	case 1:
		embed := embedutil.New("Configurator", "Configure per-guild overrides.")
		embed.Color = multiplexer.KappaColor
		for _, entry := range overrides.GetSimpleEntries() {
			embed.AddField(entry.Name, entry.Description, false)
		}
		for _, entry := range overrides.GetComplexEntries() {
			embed.AddField(entry.Name, entry.Description, false)
		}
		for _, entry := range overrides.GetCustomEntries() {
			embed.AddField(entry.Name, entry.Description, false)
		}
		context.SendEmbed("", embed)
	case 2:
		for _, entry := range overrides.GetSimpleEntries() {
			if context.Fields[1] == entry.Name {
				embed := embedutil.New(entry.FriendlyName, entry.Description)
				embed.Color = multiplexer.KappaColor
				value, err := config.GetGuildConfValue(context.Guild.ID, entry.DatabaseKey)
				if !context.HandleError(err) {
					return
				}
				title, description, ok := entry.Format(context, value)
				if !ok {
					return
				}
				embed.AddField(title, description, true)
				context.SendEmbed("", embed)
				return
			}
		}
		for _, entry := range overrides.GetComplexEntries() {
			if context.Fields[1] == entry.Name {
				embed := embedutil.New(entry.FriendlyName, entry.Description)
				embed.Color = multiplexer.KappaColor
				for _, subEntry := range entry.Entries {
					embed.AddField(subEntry.Name, subEntry.Description, false)
				}
				for _, subEntry := range entry.CustomEntries {
					embed.AddField(subEntry.Name, subEntry.Description, false)
				}
				context.SendEmbed("", embed)
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
			context.SendMessage(multiplexer.InvalidArgument)
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
					context.SendMessage(fmt.Sprintf("Successfully reset value of `%s`.", entry.Name))
					return
				}
				input := context.StitchFields(2)
				valid, ok := entry.Validate(context, &input)
				if !ok {
					return
				}
				if !valid {
					context.SendMessage(multiplexer.InvalidArgument)
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
				for _, subEntry := range entry.CustomEntries {
					if context.Fields[2] == subEntry.Name {
						subEntry.Handler(context)
						return
					}
				}
				for _, subEntry := range entry.Entries {
					if context.Fields[2] == subEntry.Name {
						if len(context.Fields) == 3 {
							embed := embedutil.New(subEntry.FriendlyName, subEntry.Description)
							embed.Color = multiplexer.KappaColor
							value, err := config.GetGuildConfValue(context.Guild.ID, subEntry.DatabaseKey)
							if !context.HandleError(err) {
								return
							}
							title, description, ok := subEntry.Format(context, value)
							if !ok {
								return
							}
							embed.AddField(title, description, true)
							context.SendEmbed("", embed)
							return
						}
						if context.Fields[3] == "reset" {
							err := config.ResetGuildConfValue(context.Guild.ID, subEntry.DatabaseKey)
							if !context.HandleError(err) {
								return
							}
							subEntry.Cleanup(context)
							context.SendMessage(fmt.Sprintf("Successfully reset value of `%s.%s`.", entry.Name, subEntry.Name))
							return
						}
						input := context.StitchFields(3)
						valid, ok := subEntry.Validate(context, &input)
						if !ok {
							return
						}
						if !valid {
							context.SendMessage(multiplexer.InvalidArgument)
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
				break
			}
		}
		context.SendMessage(multiplexer.InvalidArgument)
	}
}
