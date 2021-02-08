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

// Discord OAuth2 scope.
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

// Endpoint returns discord OAuth2 endpoints.
func Endpoint() oauth2.Endpoint {
	return endpoint
}

// GetToken gets a token stored in session in a context.
func GetToken(context *gin.Context) *oauth2.Token {
	session := sessions.Default(context)
	tokenJSON := session.Get("token")
	var token = &oauth2.Token{}
	err := jsoniter.UnmarshalFromString(fmt.Sprint(tokenJSON), token)
	if err != nil {
		RemoveToken(context)
		return nil
	}
	return token
}

// StoreToken stores a token in session in a context.
func StoreToken(context *gin.Context, token *oauth2.Token) {
	session := sessions.Default(context)
	tokenJSON, err := jsoniter.Marshal(token)
	if err != nil {
		panic(err)
	}
	session.Set("token", string(tokenJSON))
	_ = session.Save()
}

// RemoveToken removes a token in session in a context.
func RemoveToken(context *gin.Context) {
	session := sessions.Default(context)
	session.Delete("token")
	_ = session.Save()
}

// Client returns pointer to an http.Client.
func Client(context *gin.Context, conf *oauth2.Config) *http.Client {
	return conf.Client(context, GetToken(context))
}
