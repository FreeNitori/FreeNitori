package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"gopkg.in/macaron.v1"
	"net/http"
)

// FIXME: Make these static

func init() {
	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/guild/:gid/leaderboard",
			Handlers: []macaron.Handler{leaderboard},
		})
}

func leaderboard(context *macaron.Context) {
	guild := discord.FetchGuild(context.Params("gid"))
	if guild == nil {
		context.HTML(http.StatusNotFound, "error", datatypes.H{
			"Title":    datatypes.NoSuchFileOrDirectory,
			"Subtitle": "This guild doesn't seem to exist.",
			"Message":  "Maybe you got the wrong URL?",
		})
		return
	}
	expEnabled, err := config.ExpEnabled(guild.ID)
	if err != nil {
		context.HTML(http.StatusInternalServerError, "error", datatypes.H{
			"Title":    datatypes.InternalServerError,
			"Subtitle": "Failed to fetch experience system enablement status.",
			"Message":  "Nitori taking a nap?",
		})
		return
	}
	if !expEnabled {
		context.HTML(http.StatusServiceUnavailable, "error", datatypes.H{
			"Title":    datatypes.ServiceUnavailable,
			"Subtitle": "This feature is disabled in your guild.",
			"Message":  "Moderators don't like Nitori?",
		})
		return
	}
	context.HTML(http.StatusOK, "leaderboard", datatypes.H{
		"GuildName": guild.Name,
		"GuildIcon": guild.IconURL(),
	})
}
