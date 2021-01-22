package oauth

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	ScopeIdentify             = "identify"
	ScopeBot                  = "bot"
	ScopeEmail                = "email"
	ScopeGuilds               = "guilds"
	ScopeGuildsJoin           = "guilds.join"
	ScopeConnections          = "connections"
	ScopeGroupDMJoin          = "gdm.join"
	ScopeMessagesRead         = "messages.read"
	ScopeRPC                  = "rpc"
	ScopeRPCAPI               = "rpc.api"
	ScopeRPCNotificationsRead = "rpc.notifications.read"
	ScopeWebhookIncoming      = "webhook.Incoming"
)

var endpoint = oauth2.Endpoint{
	AuthURL:  discordgo.EndpointOauth2 + "authorize",
	TokenURL: discordgo.EndpointOauth2 + "token",
}

func Endpoint() oauth2.Endpoint {
	return endpoint
}

func GetToken(context *gin.Context) *oauth2.Token {
	session := sessions.Default(context)
	tokenJson := session.Get("token")
	var token = &oauth2.Token{}
	err := jsoniter.UnmarshalFromString(fmt.Sprint(tokenJson), token)
	if err != nil {
		RemoveToken(context)
		return nil
	}
	return token
}

func StoreToken(context *gin.Context, token *oauth2.Token) {
	session := sessions.Default(context)
	tokenJson, err := jsoniter.Marshal(token)
	if err != nil {
		panic(err)
	}
	session.Set("token", string(tokenJson))
	_ = session.Save()
}

func RemoveToken(context *gin.Context) {
	session := sessions.Default(context)
	session.Delete("token")
	_ = session.Save()
}

func Client(context *gin.Context, conf *oauth2.Config) *http.Client {
	return conf.Client(context, GetToken(context))
}
