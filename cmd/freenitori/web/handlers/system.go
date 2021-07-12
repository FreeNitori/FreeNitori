package handlers

import (
	"encoding/json"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/db"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord/snowflake"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/stats"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/structs"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/web/ws"
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
	context.JSON(http.StatusOK, structs.H{
		"nitori": "Rand!",
	})
}

func apiInfo(context *gin.Context) {
	context.String(http.StatusOK, infoPayload)
}

func apiStats(context *gin.Context) {
	context.JSON(http.StatusOK, structs.H{
		"total_messages":  db.GetTotalMessages(),
		"guilds_deployed": strconv.Itoa(len(state.Session.State.Guilds)),
	})
}

func apiNitori(context *gin.Context) {
	context.JSON(http.StatusOK, structs.UserInfo{
		Name:          state.Session.State.User.Username,
		ID:            state.Session.State.User.ID,
		CreationTime:  snowflake.CreationTime(state.Session.State.User.ID),
		AvatarURL:     state.Session.State.User.AvatarURL("4096"),
		Discriminator: state.Session.State.User.Discriminator,
		Bot:           state.Session.State.User.Bot,
	})
}

func apiNitoriUpdate(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, structs.H{"error": structs.PermissionDeniedAPI})
		return
	}
	var newInfo structs.UserInfo
	err := context.BindJSON(&newInfo)
	if err != nil {
		context.JSON(http.StatusBadRequest, structs.H{"error": "invalid json"})
		return
	}
	user, err = state.Session.UserUpdate("", "", newInfo.Name, "", "")
	if err != nil {
		context.JSON(http.StatusInternalServerError, structs.H{"error": err.Error()})
		return
	}
	state.Session.State.User = user
	context.JSON(http.StatusOK, structs.H{"state": "ok"})
}

func apiNitoriAction(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, structs.H{"error": structs.PermissionDeniedAPI})
		return
	}
	var action struct {
		Action string `json:"action"`
	}
	err := context.BindJSON(&action)
	if err != nil {
		context.JSON(http.StatusBadRequest, structs.H{"error": structs.BadRequestAPI})
		return
	}
	switch action.Action {
	case "restart":
		context.JSON(http.StatusOK, structs.H{"state": "ok"})
		go func() { state.Exit <- -1 }()
	case "shutdown":
		context.JSON(http.StatusOK, structs.H{"state": "ok"})
		go func() { state.Exit <- 0 }()
	default:
		context.JSON(http.StatusBadRequest, structs.H{"error": structs.BadRequestAPI})
	}
}

func apiNitoriStats(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, structs.H{"error": structs.PermissionDeniedAPI})
		return
	}
	context.JSON(http.StatusOK, stats.Get())
}

func apiNitoriLogs(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsOperator(user.ID) {
		context.JSON(http.StatusForbidden, structs.H{"error": structs.PermissionDeniedAPI})
		return
	}
	err := ws.WS.HandleRequest(context.Writer, context.Request)
	if err != nil {
		log.Debugf("Error handling log web socket, %s", err)
		return
	}
}

func apiNitoriBroadcast(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsOperator(user.ID) {
		context.JSON(http.StatusForbidden, structs.H{"error": structs.PermissionDeniedAPI})
		return
	}
	message, err := db.GetBroadcastMessage()
	if err != nil && err != badger.ErrKeyNotFound {
		context.JSON(http.StatusInternalServerError, structs.H{"error": err.Error()})
		return
	}
	if err == badger.ErrKeyNotFound || message == "" {
		context.JSON(http.StatusOK, structs.H{"content": "There is no content in the broadcast buffer."})
		return
	}
	context.JSON(http.StatusOK, structs.H{"content": message})
}

func apiNitoriBroadcastUpdate(context *gin.Context) {
	user := oauth.GetSelf(context)
	if user == nil || !state.Multiplexer.IsAdministrator(user.ID) {
		context.JSON(http.StatusForbidden, structs.H{"error": structs.PermissionDeniedAPI})
		return
	}
	var broadcastPayload struct {
		Alert   bool   `json:"alert"`
		Content string `json:"content"`
	}
	if err := context.BindJSON(&broadcastPayload); err != nil {
		context.JSON(http.StatusBadRequest, structs.H{"error": structs.BadRequestAPI})
		return
	}

	if err := db.SetBroadcastMessage(broadcastPayload.Content); err != nil {
		context.JSON(http.StatusInternalServerError, structs.H{"error": err.Error()})
		return
	}

	if broadcastPayload.Alert {
		now := time.Now().UTC()
		for _, operator := range state.Multiplexer.Operator {
			if id, err := state.Session.UserChannelCreate(operator.ID); err != nil {
				log.Errorf("Error creating user channel for operator %s, %s", operator.ID, err)
				continue
			} else {
				_, err = state.Session.ChannelMessageSend(id.ID,
					fmt.Sprintf("Broadcast by system administrator %s (%s) at `%s`:\n"+
						"```"+
						"%s"+
						"```"+
						"View this broadcast at %sauth/operator.",
						state.Multiplexer.Administrator.Username,
						state.Multiplexer.Administrator.ID,
						now, broadcastPayload.Content,
						config.WebServer.BaseURL))
				if err != nil {
					log.Errorf("Error broadcasting to operator %s, %s", operator.ID, err)
				}
			}
		}
	}
	context.JSON(http.StatusOK, structs.H{"state": "ok"})
}
