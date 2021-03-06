package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/structs"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func init() {
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
		routes.WebRoute{
			Pattern:  "/auth/operator",
			Handlers: []gin.HandlerFunc{authOperator},
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
		context.HTML(http.StatusBadRequest, "error.tmpl", structs.H{
			"Title":    structs.BadRequest,
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
	if user == nil {
		context.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	}
	if !state.Multiplexer.IsAdministrator(user.ID) {
		context.HTML(http.StatusForbidden, "error.tmpl", structs.H{
			"Title":    "Forbidden",
			"Subtitle": "This place is only for the system administrator.",
			"Message":  "What do you want...?",
		})
		return
	}
	context.HTML(http.StatusOK, "admin.tmpl", structs.H{
		"Title": "FreeNitori System",
	})
}

func authOperator(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil {
		context.Redirect(http.StatusTemporaryRedirect, "/auth/login")
		return
	}
	if !state.Multiplexer.IsOperator(user.ID) {
		context.HTML(http.StatusForbidden, "error.tmpl", structs.H{
			"Title":    "Forbidden",
			"Subtitle": "This place is only for operators.",
			"Message":  "What do you want...?",
		})
		return
	}
	context.HTML(http.StatusOK, "operator.tmpl", structs.H{
		"Title": "FreeNitori System",
	})
}
