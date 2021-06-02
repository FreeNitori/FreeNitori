package paging

import (
	embedutil "git.randomchars.net/FreeNitori/EmbedUtil"
	"github.com/bwmarrin/discordgo"
	"time"
)

type PagedMessage struct {
	Pages         []embedutil.Embed
	Page          int
	Session       *discordgo.Session
	Invoker       *discordgo.User
	Message       *discordgo.Message
	hasPermission bool
}

var instances = map[string]*PagedMessage{}

func RegisterMessage(message *PagedMessage) {
	permissions, err := message.Session.State.UserChannelPermissions(message.Session.State.User.ID, message.Message.ChannelID)
	message.hasPermission = err == nil && (int(permissions)&discordgo.PermissionManageMessages == discordgo.PermissionManageMessages)
	if message.hasPermission {
		_ = message.Session.MessageReactionAdd(message.Message.ChannelID, message.Message.ID, "⬅️")
		_ = message.Session.MessageReactionAdd(message.Message.ChannelID, message.Message.ID, "➡️")
	}
	instances[message.Message.ID] = message
	go timeoutCounter(message.Message.ID)
}

func timeoutCounter(id string) {
	time.Sleep(time.Minute)
	delete(instances, id)
}
