package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/oauth"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var oauthConf *oauth2.Config

func init() {

	go func() {
		<-state.DiscordReady
		oauthConf = &oauth2.Config{
			ClientID:     vars.Application.ID,
			ClientSecret: config.Config.Discord.ClientSecret,
			Endpoint:     oauth.Endpoint(),
			RedirectURL:  config.Config.WebServer.BaseURL + "auth/callback",
			Scopes:       []string{oauth.ScopeIdentify, oauth.ScopeGuilds},
		}
	}()

	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/auth/login",
			Handlers: []gin.HandlerFunc{authLogin},
		},
		routes.WebRoute{
			Pattern:  "/auth/callback",
			Handlers: []gin.HandlerFunc{authCallback},
		},
	)
}

func authLogin(context *gin.Context) {
	// TODO: login
}

func authCallback(context *gin.Context) {
	// TODO: callback
}
