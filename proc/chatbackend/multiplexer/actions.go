package multiplexer

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/embedutil"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/state"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

var numericalRegex *regexp.Regexp

func init() {
	numericalRegex, _ = regexp.Compile("[^0-9]+")
}

// SendMessage sends a text message in the current channel and returns the message.
func (context *Context) SendMessage(message string) *discordgo.Message {
	resultMessage, err := context.Session.ChannelMessageSend(context.Message.ChannelID, message)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		log.Errorf("Error while sending message to guild %s, %s", context.Message.GuildID, err)
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			state.ErrorOccurred)
		return nil
	}
	return resultMessage
}

// SendEmbed sends an embedutil message in the current channel and returns the message.
func (context *Context) SendEmbed(embed *embedutil.Embed) *discordgo.Message {
	resultMessage, err := context.Session.ChannelMessageSendEmbed(context.Message.ChannelID, embed.MessageEmbed)
	if err != nil {
		if err == discordgo.ErrUnauthorized {
			return nil
		}
		log.Errorf("Error while sending embedutil to guild %s, %s", context.Message.GuildID, err)
		_, _ = context.Session.ChannelMessageSend(context.Message.ChannelID,
			state.ErrorOccurred)
		return nil
	}
	return resultMessage
}

// HandleError handles a returned error and send the information of it if in debug mode.
func (context *Context) HandleError(err error) bool {
	if err != nil {
		log.Errorf("Error occurred while executing command, %s", err)
		context.SendMessage(state.ErrorOccurred)
		if log.GetLevel() == logrus.DebugLevel {
			context.SendMessage(err.Error())
		}
		return false
	}
	return true
}

// HasPermission checks a user for a permission.
func (context *Context) HasPermission(permission int) bool {
	// Override check for operators and system administrators
	if context.Author.ID == state.Administrator.ID {
		return true
	} else {
		for _, user := range state.Operator {
			if context.Author.ID == user.ID {
				return true
			}
		}
	}

	// Check against the user
	permissions, err := context.Session.State.UserChannelPermissions(context.Author.ID, context.Message.ChannelID)
	return err == nil && (permissions&permission == permission)
}

// IsOperator checks of a user is an operator.
func (context *Context) IsOperator() bool {
	if context.Author.ID == strconv.Itoa(config.Config.System.Administrator) {
		return true
	}
	for _, id := range config.Config.System.Operator {
		if context.Author.ID == strconv.Itoa(id) {
			return true
		}
	}
	return false
}

// IsAdministrator checks of a user is the system administrator.
func (context *Context) IsAdministrator() bool {
	return context.Author.ID == strconv.Itoa(config.Config.System.Administrator)
}

// GetMember gets a member from a string representing it.
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

// StitchFields stitches together fields of the message.
func (context *Context) StitchFields(start int) string {
	message := context.Fields[1]
	for i := start + 1; i < len(context.Fields); i++ {
		message += " " + context.Fields[i]
	}
	return message
}

// GenerateGuildPrefix returns the command prefix of a context.
func (context *Context) GenerateGuildPrefix() string {
	switch context.IsPrivate {
	case true:
		return config.Config.System.Prefix
	case false:
		return config.GetPrefix(context.Guild.ID)
	}
	return ""
}

// GetVoiceState returns the voice state of a user if found.
func (context *Context) GetVoiceState() (*discordgo.VoiceState, bool) {
	if context.IsPrivate {
		return nil, false
	}
	for _, voiceState := range context.Guild.VoiceStates {
		if voiceState.UserID == context.Author.ID {
			return voiceState, true
		}
	}
	return nil, false
}

// MakeVoiceConnection returns the voice connection to a user's voice channel if join-able.
func (context *Context) MakeVoiceConnection() (*discordgo.VoiceConnection, error) {
	if context.IsPrivate {
		return nil, nil
	}
	voiceState, ok := context.GetVoiceState()
	if !ok {
		return nil, nil
	}
	return context.Session.ChannelVoiceJoin(voiceState.GuildID, voiceState.ChannelID, false, true)
}
