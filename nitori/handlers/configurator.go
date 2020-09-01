package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"unicode"
)

var err error

func (*CommandHandlers) Configure(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
		return
	}
	if !context.HasPermission(discordgo.PermissionAdministrator) {
		context.SendMessage(state.PermissionDenied)
		return
	}
	if len(context.Fields) == 1 {
		embed := formatter.NewEmbed("Configurator", "Configure per-guild overrides.")
		embed.Color = state.KappaColor
		embed.AddField("experience", "Toggle experience system enablement.", false)
		embed.AddField("message", "Configure customizable messages.", false)
		embed.AddField("prefix", "Configure command prefix.", false)
		context.SendEmbed(embed)
		return
	}
	switch context.Fields[1] {
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
			if !context.HandleError(err, config.Debug) {
				return
			}
			context.SendMessage("Successfully updated prefix.")
		case 2:
			err = config.ResetPrefix(context.Guild.ID)
			if !context.HandleError(err, config.Debug) {
				return
			}
			context.SendMessage("Successfully reset prefix.")
		default:
			context.SendMessage("Invalid syntax, please check your parameters and try again.")
		}
	case "experience":
		pre, err := config.ExpToggle(context.Guild.ID)
		if !context.HandleError(err, config.Debug) {
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
			message := context.Fields[3]
			for i := 4; i < len(context.Fields); i++ {
				message += " " + context.Fields[i]
			}
			err := config.SetCustomizableMessage(context.Guild.ID, context.Fields[2], message)
			switch err.(type) {
			default:
				if !context.HandleError(err, config.Debug) {
					return
				}
			case *config.MessageOutOfBounds:
				context.SendMessage(state.InvalidArgument)
				return
			}
			context.SendMessage("Message `" + context.Fields[2] + "` has been set.")
		case 3:
			err := config.SetCustomizableMessage(context.Guild.ID, context.Fields[2], "")
			switch err.(type) {
			default:
				if !context.HandleError(err, config.Debug) {
					return
				}
			case *config.MessageOutOfBounds:
				context.SendMessage(state.InvalidArgument)
				return
			}
			context.SendMessage("Message `" + context.Fields[2] + "` has been reset.")
		case 2:
			embed := formatter.NewEmbed("Messages", "Configurable messages.")
			for identifier := range config.CustomizableMessages {
				message, err := config.GetCustomizableMessage(context.Guild.ID, identifier)
				if !context.HandleError(err, config.Debug) {
					return
				}
				embed.AddField(identifier, message, false)
			}
			context.SendEmbed(embed)
		}
	}
}
