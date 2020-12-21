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
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "about",
		AliasPatterns: []string{"info", "kappa", "information"},
		Description:   "Display system information.",
		Category:      multiplexer.SystemCategory,
		Handler:       about,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "invite",
		AliasPatterns: []string{"authorize", "oauth"},
		Description:   "Display authorization URL.",
		Category:      multiplexer.SystemCategory,
		Handler:       invite,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "reset-guild",
		AliasPatterns: []string{},
		Description:   "Reset current guild configuration.",
		Category:      multiplexer.SystemCategory,
		Handler:       resetGuild,
	})
	multiplexer.Router.Route(&multiplexer.Route{
		Pattern:       "shutdown",
		AliasPatterns: []string{"poweroff", "reboot", "restart"},
		Description:   "",
		Category:      multiplexer.SystemCategory,
		Handler:       shutdown,
	})
}

func about(context *multiplexer.Context) {
	embed := embedutil.NewEmbed(context.Session.State.User.Username,
		"Open source, general purpose Discord utility.")
	embed.Color = vars.KappaColor
	embed.AddField("Homepage", config.Config.WebServer.BaseURL, true)
	embed.AddField("Version", state.Version, true)
	embed.AddField("Commit Hash", state.Revision, true)
	embed.AddField("Processed Messages", strconv.Itoa(config.GetTotalMessages()), true)
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
	if map[string]bool{"reboot": true, "restart": true, "shutdown": false, "poweroff": false}[context.Fields[0]] {
		context.SendMessage("Attempting restart...")
		state.ExitCode <- -1
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
