package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
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
	context.JSON(http.StatusOK, datatypes.H{
		"authorized": oauth.GetToken(context) != nil,
	})
}

func apiAuthUser(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil {
		context.JSON(http.StatusOK, datatypes.H{
			"authorized":    false,
			"administrator": false,
			"user":          datatypes.UserInfo{},
		})
		return
	}

	context.JSON(http.StatusOK, datatypes.H{
		"authorized":    true,
		"administrator": user.ID == state.Multiplexer.Administrator.ID,
		"user": datatypes.UserInfo{
			Name:          user.Username,
			ID:            user.ID,
			AvatarURL:     user.AvatarURL("4096"),
			Discriminator: user.Discriminator,
			CreationTime:  config.CreationTime(user.ID),
			Bot:           user.Bot,
		},
	})
}
