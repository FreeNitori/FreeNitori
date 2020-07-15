package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"strconv"
)

func (*Handlers) About(context *multiplexer.Context) {
	embed := formatter.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = KappaColor
	embed.AddField("Homepage", config.BaseURL, true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
	embed.AddField("License", "GNU General Public License v3.0", false)
	embed.SetThumbnail(context.Session.State.User.AvatarURL("256"))
	embed.SetFooter("A Discord utility by RandomChars", "https://static.randomchars.net/img/RandomChars.png")
	context.SendEmbed(embed,
		"producing system info Embed")
}

func (handlers *Handlers) Reboot(context *multiplexer.Context) {
	if context.Author.ID != config.Administrator {
		context.SendMessage(AdminOnly, "generating permission denied message")
		return
	}
	switch context.Fields[0] {
	case "reboot", "restart":
		context.SendMessage("Rebooting chat backend.", "generating chat backend reboot message")
		_ = multiplexer.IPCConnection.Call("IPC.Restart", []string{"ChatBackend"}, nil)
		multiplexer.ExitCode <- 0
		return
	case "halt", "shutdown":
		context.SendMessage("Performing complete shutdown.", "generating system shutdown message")
		if context.Fields[0] == "shutdown" {
			_ = multiplexer.IPCConnection.Call("IPC.Shutdown", []string{"ChatBackend"}, nil)
			multiplexer.ExitCode <- 0
			return
		}
	}
}
