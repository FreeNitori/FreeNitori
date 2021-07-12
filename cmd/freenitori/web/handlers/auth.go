package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord/snowflake"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/structs"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"github.com/gin-gonic/gin"
	"net/http"
)

func init() {
	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/api/auth",
			Handlers: []gin.HandlerFunc{apiAuth},
		},
		routes.WebRoute{
			Pattern:  "/api/auth/user",
			Handlers: []gin.HandlerFunc{apiAuthUser},
		},
	)
}

func apiAuth(context *gin.Context) {
	context.JSON(http.StatusOK, structs.H{
		"authorized": oauth.GetToken(context) != nil,
	})
}

func apiAuthUser(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil {
		context.JSON(http.StatusOK, structs.H{
			"authorized":    false,
			"operator":      false,
			"administrator": false,
			"user":          structs.UserInfo{},
		})
		return
	}

	context.JSON(http.StatusOK, structs.H{
		"authorized":    true,
		"operator":      state.Multiplexer.IsOperator(user.ID),
		"administrator": state.Multiplexer.IsAdministrator(user.ID),
		"user": structs.UserInfo{
			Name:          user.Username,
			ID:            user.ID,
			AvatarURL:     user.AvatarURL("4096"),
			Discriminator: user.Discriminator,
			CreationTime:  snowflake.CreationTime(user.ID),
			Bot:           user.Bot,
		},
	})
}
