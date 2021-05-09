package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func init() {
	routes.GetRoutes = append(routes.GetRoutes,
		routes.WebRoute{
			Pattern:  "/api",
			Handlers: []gin.HandlerFunc{api},
		},
		routes.WebRoute{
			Pattern:  "/api/info",
			Handlers: []gin.HandlerFunc{apiInfo},
		},
		routes.WebRoute{
			Pattern:  "/api/stats",
			Handlers: []gin.HandlerFunc{apiStats},
		},
	)
}

func api(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"status": "OK!",
	})
}

func apiInfo(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"nitori_version":  state.Version(),
		"nitori_revision": state.Revision(),
		"invite_url":      state.InviteURL,
	})
}

func apiStats(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"total_messages":  config.GetTotalMessages(),
		"guilds_deployed": strconv.Itoa(len(state.RawSession.State.Guilds)),
	})
}
