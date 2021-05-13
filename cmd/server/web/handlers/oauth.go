package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/http"
)

func init() {

	go func() {
		<-state.DiscordReady
		oauth.Conf = &oauth2.Config{
			ClientID:     state.Application.ID,
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
		routes.WebRoute{
			Pattern:  "/auth/admin",
			Handlers: []gin.HandlerFunc{authAdmin},
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
	context.Redirect(http.StatusTemporaryRedirect, oauth.Conf.AuthCodeURL(oauthState))
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
	token, err := oauth.Conf.Exchange(context, context.Request.FormValue("code"))
	if err != nil {
		panic(err)
	}
	oauth.StoreToken(context, token)
	_ = session.Save()
	context.Redirect(http.StatusTemporaryRedirect, "/")
}

func authAdmin(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user.ID != state.Multiplexer.Administrator.ID {
		context.HTML(http.StatusForbidden, "error.tmpl", datatypes.H{
			"Title":    "Forbidden",
			"Subtitle": "This place is only for the system administrator.",
			"Message":  "What do you want...?",
		})
		return
	}
	context.HTML(http.StatusOK, "admin.tmpl", datatypes.H{
		"Title": "FreeNitori System",
	})
}
