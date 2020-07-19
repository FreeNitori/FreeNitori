package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"github.com/bwmarrin/discordgo"
)

func (context *Context) SendMessage(message string, action string) *discordgo.Message {
	var err error

	resultMessage, err := context.Session.ChannelMessageSend(context.Message.ChannelID, message)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error while %s for guild %s, %s", action, context.Message.GuildID, err))
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			"Something went wrong and the kappa is very confused! Please try again!")
	}
	return resultMessage
}

func (context *Context) SendEmbed(embed *formatter.Embed, action string) *discordgo.Message {
	var err error

	resultMessage, err := context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed.MessageEmbed)
	if err != nil {
		Logger.Error(fmt.Sprintf("Error while %s for guild %s, %s", action, context.Message.GuildID, err))
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			"Something went wrong and the kappa is very confused! Please try again!")
	}
	return resultMessage
}
