package handlers

import (
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/oauth"
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
		routes.WebRoute{
			Pattern:  "/api/nitori",
			Handlers: []gin.HandlerFunc{apiNitori},
		},
	)
	routes.PostRoutes = append(routes.PostRoutes,
		routes.WebRoute{
			Pattern:  "/api/nitori",
			Handlers: []gin.HandlerFunc{apiNitoriUpdate},
		},
		routes.WebRoute{
			Pattern:  "/api/nitori/action",
			Handlers: []gin.HandlerFunc{apiNitoriAction},
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

func apiNitori(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.UserInfo{
		Name:          state.RawSession.State.User.Username,
		ID:            state.RawSession.State.User.ID,
		CreationTime:  config.CreationTime(state.RawSession.State.User.ID),
		AvatarURL:     state.RawSession.State.User.AvatarURL("4096"),
		Discriminator: state.RawSession.State.User.Discriminator,
		Bot:           state.RawSession.State.User.Bot,
	})
}

func apiNitoriUpdate(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user.ID != state.Multiplexer.Administrator.ID {
		context.JSON(http.StatusForbidden, datatypes.H{"error": "permission denied"})
		return
	}
	var newInfo datatypes.UserInfo
	err := context.BindJSON(&newInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, datatypes.H{"error": "invalid json"})
		return
	}
	user, err = state.RawSession.UserUpdate("", "", newInfo.Name, "", "")
	if err != nil {
		context.JSON(http.StatusInternalServerError, datatypes.H{"error": err})
		return
	}
	state.RawSession.State.User = user
	context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
}

func apiNitoriAction(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user.ID != state.Multiplexer.Administrator.ID {
		context.JSON(http.StatusForbidden, datatypes.H{"error": "permission denied"})
		return
	}
	var action struct {
		Action string `json:"action"`
	}
	err := context.BindJSON(&action)
	if err != nil {
		context.JSON(http.StatusBadRequest, datatypes.H{"error": "invalid json"})
		return
	}
	switch action.Action {
	case "restart":
		context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
		state.ExitCode <- -1
	case "shutdown":
		context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
		state.ExitCode <- 0
	default:
		context.JSON(http.StatusBadRequest, datatypes.H{"error": "invalid action"})
	}
}
