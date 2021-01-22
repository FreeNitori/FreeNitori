package handlers

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/datatypes"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/oauth"
	"git.randomchars.net/RandomChars/FreeNitori/server/web/routes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
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
			Pattern:  "/auth/logout",
			Handlers: []gin.HandlerFunc{authLogout},
		},
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

func authLogout(context *gin.Context) {
	oauth.RemoveToken(context)
	context.Redirect(http.StatusTemporaryRedirect, "/")
}

func authLogin(context *gin.Context) {
	session := sessions.Default(context)
	oauthState := uuid.New().String()
	session.Set("state", oauthState)
	_ = session.Save()
	context.Redirect(http.StatusTemporaryRedirect, oauthConf.AuthCodeURL(oauthState))
}

func authCallback(context *gin.Context) {
	session := sessions.Default(context)
	if context.Request.FormValue("state") != session.Get("state") {
		context.HTML(http.StatusBadRequest, "error.tmpl", datatypes.H{
			"Title":    datatypes.BadRequest,
			"Subtitle": "State doesn't seem to match.",
			"Message":  "Trying to be sneaky...?",
		})
		return
	}
	session.Delete("state")
	token, err := oauthConf.Exchange(context, context.Request.FormValue("code"))
	if err != nil {
		panic(err)
	}
	oauth.StoreToken(context, token)
	_ = session.Save()
	context.Redirect(http.StatusTemporaryRedirect, "/")
}
