package multiplexer

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strconv"
	"strings"
)

var numericalRegex *regexp.Regexp

func init() {
	numericalRegex, _ = regexp.Compile("[^0-9]+")
}

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

// Get a guild member from a string
func (context *Context) GetMember(user string) *discordgo.Member {
	// Check if it's a mention or the string is numerical
	_, err := strconv.Atoi(user)
	if strings.HasPrefix(user, "<@") && strings.HasSuffix(user, ">") || err == nil {
		// Strip off the mention thingy
		userID := numericalRegex.ReplaceAllString(user, "")
		// Length of a real user ID after stripping off stuff
		if len(userID) == 18 {
			for _, member := range context.Guild.Members {
				if member.User.ID == userID {
					return member
				}
			}
		}
	} else {
		// Find as username or nickname
		for _, member := range context.Guild.Members {
			if member.User.Username == user || member.Nick == user {
				return member
			}
		}
	}
	return nil
}
