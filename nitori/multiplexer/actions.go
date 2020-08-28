package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
)

func (context *Context) SendMessage(message string) *discordgo.Message {
	var err error

	resultMessage, err := context.Session.ChannelMessageSend(context.Message.ChannelID, message)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		state.Logger.Error(fmt.Sprintf("Error while sending message to guild %s, %s", context.Message.GuildID, err))
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			state.ErrorOccurred)
		return nil
	}
	return resultMessage
}

func (context *Context) SendEmbed(embed *formatter.Embed) *discordgo.Message {
	var err error

	resultMessage, err := context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed.MessageEmbed)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		state.Logger.Error(fmt.Sprintf("Error while sending embed to guild %s, %s", context.Message.GuildID, err))
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			state.ErrorOccurred)
		return nil
	}
	return resultMessage
}

func (context *Context) HandleError(err error, debug bool) bool {
	if err != nil {
		context.SendMessage(state.ErrorOccurred)
		if debug {
			context.SendMessage(err.Error())
		}
		return false
	}
	return true
}

func (context *Context) HasPermission(permission int) bool {
	if context.Author.ID == config.Operator || context.Author.ID == config.Administrator {
		return true
	}
	permissions, err := context.Session.State.UserChannelPermissions(context.Author.ID, context.Message.ChannelID)
	return err == nil && (permissions&permission == permission)
}
