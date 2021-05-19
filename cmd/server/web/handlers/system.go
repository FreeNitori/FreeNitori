package handlers

import (
	"encoding/json"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/ws"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

var infoPayload string

// PopulateInfoPayload populates /api/info payload.
func PopulateInfoPayload() {
	payload, err := json.Marshal(struct {
		Version   string `json:"nitori_version"`
		Revision  string `json:"nitori_revision"`
		InviteURL string `json:"invite_url"`
	}{
		Version:   state.Version(),
		Revision:  state.Revision(),
		InviteURL: state.InviteURL,
	})
	if err != nil {
		panic(err)
	}
	infoPayload = string(payload)
}

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
		routes.WebRoute{
			Pattern:  "/api/nitori/stats",
			Handlers: []gin.HandlerFunc{apiNitoriStats},
		},
		routes.WebRoute{
			Pattern:  "/api/nitori/logs",
			Handlers: []gin.HandlerFunc{apiNitoriLogs},
		},
		routes.WebRoute{
			Pattern:  "/api/nitori/broadcast",
			Handlers: []gin.HandlerFunc{apiNitoriBroadcast},
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
		routes.WebRoute{
			Pattern:  "/api/nitori/broadcast",
			Handlers: []gin.HandlerFunc{apiNitoriBroadcastUpdate},
		},
	)
}

func api(context *gin.Context) {
	context.JSON(http.StatusOK, datatypes.H{
		"nitori": "Rand!",
	})
}

func apiInfo(context *gin.Context) {
	context.String(http.StatusOK, infoPayload)
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
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, datatypes.H{"error": datatypes.PermissionDeniedAPI})
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
		context.JSON(http.StatusInternalServerError, datatypes.H{"error": err.Error()})
		return
	}
	state.RawSession.State.User = user
	context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
}

func apiNitoriAction(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, datatypes.H{"error": datatypes.PermissionDeniedAPI})
		return
	}
	var action struct {
		Action string `json:"action"`
	}
	err := context.BindJSON(&action)
	if err != nil {
		context.JSON(http.StatusBadRequest, datatypes.H{"error": datatypes.BadRequestAPI})
		return
	}
	switch action.Action {
	case "restart":
		context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
		go func() { state.ExitCode <- -1 }()
	case "shutdown":
		context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
		go func() { state.ExitCode <- 0 }()
	default:
		context.JSON(http.StatusBadRequest, datatypes.H{"error": datatypes.BadRequestAPI})
	}
}

func apiNitoriStats(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, datatypes.H{"error": datatypes.PermissionDeniedAPI})
		return
	}
	context.JSON(http.StatusOK, config.Stats())
}

func apiNitoriLogs(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsOperator(user.ID) {
		context.JSON(http.StatusForbidden, datatypes.H{"error": datatypes.PermissionDeniedAPI})
		return
	}
	err := ws.WS.HandleRequest(context.Writer, context.Request)
	if err != nil {
		log.Debugf("Error while handling log web socket, %s", err)
		return
	}
}

func apiNitoriBroadcast(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsOperator(user.ID) {
		context.JSON(http.StatusForbidden, datatypes.H{"error": datatypes.PermissionDeniedAPI})
		return
	}
	message, err := config.GetBroadcastMessage()
	if err != nil && err != badger.ErrKeyNotFound {
		context.JSON(http.StatusInternalServerError, datatypes.H{"error": err.Error()})
		return
	}
	if err == badger.ErrKeyNotFound || message == "" {
		context.JSON(http.StatusOK, datatypes.H{"content": "There is no content in the broadcast buffer."})
		return
	}
	context.JSON(http.StatusOK, datatypes.H{"content": message})
}

func apiNitoriBroadcastUpdate(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, datatypes.H{"error": datatypes.PermissionDeniedAPI})
		return
	}
	var broadcastPayload struct {
		Alert   bool   `json:"alert"`
		Content string `json:"content"`
	}
	err := context.BindJSON(&broadcastPayload)
	if err != nil {
		context.JSON(http.StatusBadRequest, datatypes.H{"error": datatypes.BadRequestAPI})
		return
	}
	err = config.SetBroadcastMessage(broadcastPayload.Content)
	if err != nil {
		context.JSON(http.StatusInternalServerError, datatypes.H{"error": err.Error()})
		return
	}
	if broadcastPayload.Alert {
		now := time.Now().UTC()
		for _, user := range state.Multiplexer.Operator {
			id, err := state.RawSession.UserChannelCreate(user.ID)
			if err != nil {
				log.Errorf("Error while creating user channel for operator %s, %s", user.ID, err)
				continue
			}
			_, err = state.RawSession.ChannelMessageSend(id.ID,
				fmt.Sprintf("Broadcast by system administrator %s (%s) at `%s`:\n"+
					"```"+
					"%s"+
					"```"+
					"View this broadcast at %sauth/operator.",
					state.Multiplexer.Administrator.Username,
					state.Multiplexer.Administrator.ID,
					now, broadcastPayload.Content,
					config.Config.WebServer.BaseURL))
			if err != nil {
				log.Errorf("Error while broadcasting to operator %s, %s", user.ID, err)
			}
		}
	}
	context.JSON(http.StatusOK, datatypes.H{"state": "ok"})
}
