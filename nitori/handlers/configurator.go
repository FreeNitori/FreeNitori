package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"unicode"
)

var err error

func (*CommandHandlers) Configure(context *multiplexer.Context) {
	if context.IsPrivate {
		context.SendMessage(state.GuildOnly)
		return
	}
	if len(context.Fields) == 1 {
		embed := formatter.NewEmbed("Configurator", "Configure per-guild overrides.")
		embed.Color = state.KappaColor
		embed.AddField("prefix", "Configure command prefix.", false)
		embed.AddField("experience", "Toggle experience system enablement.", false)
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
			if err != nil {
				context.SendMessage("Failed to set custom prefix, please try again later.")
				return
			}
			context.SendMessage("Successfully updated prefix.")
		case 2:
			err = config.ResetPrefix(context.Guild.ID)
			if err != nil {
				context.SendMessage("Failed to reset prefix, please try again later.")
				return
			}
			context.SendMessage("Successfully reset prefix.")
		default:
			context.SendMessage("Invalid syntax, please check your parameters and try again.")
		}
	case "experience":
		pre, err := config.ExpToggle(context.Guild.ID)
		if err != nil {
			state.Logger.Warning(fmt.Sprintf("Failed to toggle experience enabler, %s", err))
			context.SendMessage(state.ErrorOccurred)
			return
		}
		switch pre {
		case false:
			context.SendMessage("Chat experience system has been enabled.")
		case true:
			context.SendMessage("Chat experience system has been disabled.")
		}
	}
}
