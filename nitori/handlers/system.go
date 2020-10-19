package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	"strconv"
)

func init() {
	SystemCategory.Register(about, "about", []string{"info", "kappa", "information"}, "Display system information.")
	SystemCategory.Register(reboot, "reboot", []string{"shutdown", "halt", "restart"}, "")
	SystemCategory.Register(invite, "invite", []string{"authorize", "oauth"}, "Display authorization URL.")
}

func about(context *multiplexer.Context) {
	embed := formatter.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = ChatBackend.KappaColor
	embed.AddField("Homepage", config.Config.WebServer.BaseURL, true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
	embed.AddField("License", "GNU General Public License v3.0", false)
	if ChatBackend.Administrator != nil {
		embed.AddField("Administrator", ChatBackend.Administrator.Username+"#"+ChatBackend.Administrator.Discriminator, true)
	}
	switch len(ChatBackend.Operator) {
	case 0:
	case 1:
		embed.AddField("Operator", ChatBackend.Operator[0].Username+"#"+ChatBackend.Operator[0].Discriminator, true)
	default:
		var usernames string
		for i, user := range ChatBackend.Operator {
			switch i {
			case 0:
				usernames += user.Username + "#" + user.Discriminator
			default:
				usernames += ", " + user.Username + "#" + user.Discriminator
			}
		}
		embed.AddField("Operators", usernames, true)
	}
	embed.SetThumbnail(context.Session.State.User.AvatarURL("256"))
	embed.SetFooter("A Discord utility by RandomChars", "https://static.randomchars.net/img/RandomChars.png")
	context.SendEmbed(embed)
}

func resetGuild(context *multiplexer.Context) {
	if !context.IsOperator() {
		context.SendMessage(ChatBackend.OperatorOnly)
		return
	}
}

func reboot(context *multiplexer.Context) {
	if !context.IsAdministrator() {
		context.SendMessage(ChatBackend.AdminOnly)
		return
	}
	switch context.Fields[0] {
	case "reboot", "restart":
		context.SendMessage("Rebooting chat backend.")
		_ = state.IPCConnection.Call("IPC.Restart", []string{"ChatBackend"}, nil)
		state.ExitCode <- 0
		return
	case "halt", "shutdown":
		context.SendMessage("Performing complete shutdown.")
		if context.Fields[0] == "shutdown" {
			_ = state.IPCConnection.Call("IPC.Shutdown", []string{"ChatBackend"}, nil)
			state.ExitCode <- 0
			return
		}
	}
}

func invite(context *multiplexer.Context) {
	embed := formatter.NewEmbed("Invite", fmt.Sprintf("Click [this](%s) to invite Nitori.", state.InviteURL))
	context.SendEmbed(embed)
}
