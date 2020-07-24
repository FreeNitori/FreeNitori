package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"github.com/bwmarrin/discordgo"
)

func (context *Context) SendMessage(message string) *discordgo.Message {
	var err error

	resultMessage, err := context.Session.ChannelMessageSend(context.Message.ChannelID, message)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		Logger.Error(fmt.Sprintf("Error while sending message to guild %s, %s", context.Message.GuildID, err))
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			"Something went wrong and I am very confused! Please try again!")
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
		Logger.Error(fmt.Sprintf("Error while sending embed to guild %s, %s", context.Message.GuildID, err))
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			"Something went wrong and I am very confused! Please try again!")
		return nil
	}
	return resultMessage
}
