package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/state"
	"strconv"
)

func init() {
	SystemCategory.Register(about, "about", []string{"info", "kappa", "information"}, "Display system information.")
	SystemCategory.Register(invite, "invite", []string{"authorize", "oauth"}, "Display authorization URL.")
	SystemCategory.Register(reboot, "reboot", []string{"shutdown", "halt", "restart"}, "")
	SystemCategory.Register(resetGuild, "reset-guild", []string{}, "Reset current guild configuration.")
}

func about(context *multiplexer.Context) {
	embed := embedutil.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = state.KappaColor
	embed.AddField("Homepage", config.Config.WebServer.BaseURL, true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
	embed.AddField("License", "GNU General Public License v3.0", false)
	if state.Administrator != nil {
		embed.AddField("Administrator", state.Administrator.Username+"#"+state.Administrator.Discriminator, true)
	}
	switch len(state.Operator) {
	case 0:
	case 1:
		embed.AddField("Operator", state.Operator[0].Username+"#"+state.Operator[0].Discriminator, true)
	default:
		var usernames string
		for i, user := range state.Operator {
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
		context.SendMessage(state.OperatorOnly)
		return
	}
	config.ResetGuild(context.Guild.ID)
	context.SendMessage("Guild configuration has been reset.")
}

func reboot(context *multiplexer.Context) {
	if !context.IsAdministrator() {
		context.SendMessage(state.AdminOnly)
		return
	}
	switch context.Fields[0] {
	case "reboot", "restart":
		context.SendMessage("Rebooting chat backend.")
		_ = vars.RPCConnection.Call("R.Restart", []int{vars.ProcessType}, nil)
		vars.ExitCode <- 0
		return
	case "halt", "shutdown":
		context.SendMessage("Performing complete shutdown.")
		if context.Fields[0] == "shutdown" {
			_ = vars.RPCConnection.Call("R.Shutdown", []int{vars.ProcessType}, nil)
			vars.ExitCode <- 0
			return
		}
	}
}

func invite(context *multiplexer.Context) {
	embed := embedutil.NewEmbed("Invite", fmt.Sprintf("Click [this](%s) to invite Nitori.", vars.InviteURL))
	context.SendEmbed(embed)
}
