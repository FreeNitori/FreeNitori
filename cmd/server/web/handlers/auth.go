package handlers

import (
	"encoding/json"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/datatypes"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/oauth"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/web/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
	token := oauth.GetToken(context)
	if token == nil {
		context.JSON(http.StatusOK, datatypes.H{
			"authorized": false,
			"user":       datatypes.UserInfo{},
		})
		return
	}
	client := oauth.Client(context, oauthConf)
	response, err := client.Get(discordgo.EndpointUser("@me"))
	if err != nil {
		panic(err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode == http.StatusUnauthorized {
		oauth.RemoveToken(context)
		context.JSON(http.StatusOK, datatypes.H{
			"authorized": false,
			"user":       datatypes.UserInfo{},
		})
		return
	}

	var user discordgo.User
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		panic(err)
	}

	context.JSON(http.StatusOK, datatypes.H{
		"authorized": true,
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
