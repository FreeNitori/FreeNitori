package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"strconv"
)

func (handlers *Handlers) About(context *multiplexer.Context) {
	embed := formatter.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = KappaColor
	embed.AddField("Homepage", "Not Implemented", true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
	embed.AddField("License", "GNU General Public License v3.0", false)
	embed.AddField("System Administrator", "Not Implemented", true)
	embed.AddField("Operator", "Not Implemented", true)
	embed.SetImage(context.Session.State.User.AvatarURL("256"))
	embed.SetFooter("Discord utility by RandomChars", "https://static.randomchars.net/img/RandomChars.png")
	context.SendEmbed(embed,
		"producing system info Embed")
}
