package paging

import (
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
)

func init() {
	state.Multiplexer.MessageReactionAdd = append(state.Multiplexer.MessageReactionAdd, handlePage)
}

func handlePage(context *multiplexer.Context) {
	reactionAdd, ok := context.Event.(*discordgo.MessageReactionAdd)
	if !ok {
		return
	}
	if reactionAdd.GuildID == "" {
		return
	}
	if reactionAdd.UserID == state.RawSession.State.User.ID {
		return
	}

	// Do not handle if message does not exist
	message, ok := instances[reactionAdd.MessageID]
	if !ok {
		return
	}

	// Do not handle if not author or overriding user
	context.User.ID = reactionAdd.UserID
	if reactionAdd.UserID != message.Invoker.ID && !context.HasPermission(discordgo.PermissionManageMessages) {
		return
	}

	switch reactionAdd.Emoji.Name {
	case "⬅️":
		// Page left
		if message.hasPermission {
			_ = message.Session.MessageReactionRemove(reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.Name, reactionAdd.UserID)
		}

		if message.Page == 0 {
			return
		}
		message.Page -= 1
		_, _ = message.Session.ChannelMessageEditEmbed(message.Message.ChannelID, message.Message.ID, message.Pages[message.Page].MessageEmbed)
	case "➡️":
		// Page right
		if message.hasPermission {
			_ = message.Session.MessageReactionRemove(reactionAdd.ChannelID, reactionAdd.MessageID, reactionAdd.Emoji.Name, reactionAdd.UserID)
		}
		if message.Page == len(message.Pages)-1 {
			return
		}
		message.Page += 1
		_, _ = message.Session.ChannelMessageEditEmbed(message.Message.ChannelID, message.Message.ID, message.Pages[message.Page].MessageEmbed)
	default:
		// Do not handle if reactions do not match
		return
	}
}
