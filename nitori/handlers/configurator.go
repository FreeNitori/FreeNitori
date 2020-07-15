package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"unicode"
)

var err error

func (*Handlers) Configure(context *multiplexer.Context) {
	if len(context.Fields) == 1 {
		embed := formatter.NewEmbed("Configurator", "Configure per-guild overrides.")
		embed.Color = KappaColor
		embed.AddField("prefix", "Configure command prefix.", false)
		context.SendEmbed(embed, "sending configurator help")
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
				context.SendMessage("Failed to set custom prefix, please try again later.", "generating database error message")
				return
			}
			context.SendMessage("Successfully updated prefix.", "generating prefix update success message")
		case 2:
			err = config.ResetPrefix(context.Guild.ID)
			if err != nil {
				context.SendMessage("Failed to reset prefix, please try again later.", "generating database error message")
				return
			}
			context.SendMessage("Successfully reset prefix.", "generating prefix reset success message")
		default:
			context.SendMessage("Invalid syntax, please check your parameters and try again.", "generating invalid syntax message")
		}
	}
}
