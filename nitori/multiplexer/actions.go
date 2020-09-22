package multiplexer

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/formatter"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
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

// Send a text message and return it
func (context *Context) SendMessage(message string) *discordgo.Message {
	resultMessage, err := context.Session.ChannelMessageSend(context.Message.ChannelID, message)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		log.Logger.Errorf("Error while sending message to guild %s, %s", context.Message.GuildID, err)
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			state.ErrorOccurred)
		return nil
	}
	return resultMessage
}

// Send an embed message and return it
func (context *Context) SendEmbed(embed *formatter.Embed) *discordgo.Message {
	resultMessage, err := context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed.MessageEmbed)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		log.Logger.Errorf("Error while sending embed to guild %s, %s", context.Message.GuildID, err)
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			state.ErrorOccurred)
		return nil
	}
	return resultMessage
}

// Handle error and send the stuff if in debug mode
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

// Check if user has specific permission
func (context *Context) HasPermission(permission int) bool {
	// Override check for operators and system administrators
	if context.Author.ID == config.Operator || context.Author.ID == config.Administrator {
		return true
	}

	// Check against the user
	permissions, err := context.Session.State.UserChannelPermissions(context.Author.ID, context.Message.ChannelID)
	return err == nil && (permissions&permission == permission)
}

// Get a guild member from a string
func (context *Context) GetMember(user string) *discordgo.Member {
	// Guild only function
	if context.IsPrivate {
		return nil
	}

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

// Stitch together fields of a context
func (context *Context) StitchFields(start int) string {
	message := context.Fields[1]
	for i := start + 1; i < len(context.Fields); i++ {
		message += " " + context.Fields[i]
	}
	return message
}
