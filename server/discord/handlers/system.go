package handlers

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"strconv"
)

func init() {
	multiplexer.SystemCategory.Register(about, "about", []string{"info", "kappa", "information"}, "Display system information.")
	multiplexer.SystemCategory.Register(invite, "invite", []string{"authorize", "oauth"}, "Display authorization URL.")
	multiplexer.SystemCategory.Register(resetGuild, "reset-guild", []string{}, "Reset current guild configuration.")
	multiplexer.SystemCategory.Register(shutdown, "shutdown", []string{}, "")
}

func about(context *multiplexer.Context) {
	embed := embedutil.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = vars.KappaColor
	embed.AddField("Homepage", config.Config.WebServer.BaseURL, true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
	embed.AddField("License", "GNU General Public License v3.0", false)
	if vars.Administrator != nil {
		embed.AddField("Administrator", vars.Administrator.Username+"#"+vars.Administrator.Discriminator, true)
	}
	switch len(vars.Operator) {
	case 0:
	case 1:
		embed.AddField("Operator", vars.Operator[0].Username+"#"+vars.Operator[0].Discriminator, true)
	default:
		var usernames string
		for i, user := range vars.Operator {
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
		context.SendMessage(vars.OperatorOnly)
		return
	}
	config.ResetGuild(context.Guild.ID)
	context.SendMessage("Guild configuration has been reset.")
}

func shutdown(context *multiplexer.Context) {
	if !context.IsAdministrator() {
		context.SendMessage(vars.AdminOnly)
		return
	}
	context.SendMessage("Performing complete shutdown.")
	log.Info("Shutdown requested via Discord command.")
	state.ExitCode <- 0
}

func invite(context *multiplexer.Context) {
	embed := embedutil.NewEmbed("Invite", fmt.Sprintf("Click [this](%s) to invite Nitori.", state.InviteURL))
	context.SendEmbed(embed)
}
