package oauth

import (
	"github.com/bwmarrin/discordgo"
	"golang.org/x/oauth2"
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
