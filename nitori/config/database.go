package config

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var prefixes = []string{"conf", "exp", "rank", "exp_bl", "lastfm", "ra_metadata"}
var CustomizableMessages = map[string]string{
	"levelup": "Congratulations $USER on reaching level $LEVEL.",
}

// ResetGuild deletes all database values that belongs to a specific guild.
func ResetGuild(gid string) {
	for _, prefix := range prefixes {
		err := database.HDel(fmt.Sprintf("%s.%s", prefix, gid))
		if err != nil {
			log.Errorf("Error while resetting guild %s key %s, %s", gid, fmt.Sprintf("%s.%s", prefix, gid), err)
		}
	}
}

// getMessage gets a guild-specific string.
func getMessage(gid string, key string) (string, error) {
	return database.HGet("conf."+gid, "message."+key)
}

// setMessage sets a guild-specific string
func setMessage(gid string, key string, message string) error {
	if len(message) > 2048 {
		return &MessageOutOfBounds{}
	}
	if message == "" {
		return database.HDel("conf."+gid, "message."+key)
	}
	return database.HSet("conf."+gid, "message."+key, message)
}

// GetCustomizableMessage gets a guild-specific message within predefined messages, returning default if not present.
func GetCustomizableMessage(gid string, key string) (string, error) {
	defaultMessage, ok := CustomizableMessages[key]
	if !ok {
		return "", &MessageOutOfBounds{}
	}
	message, err := getMessage(gid, key)
	if err != nil {
		return "", err
	}
	if message == "" {
		return defaultMessage, nil
	}
	return message, nil
}

// SetCustomizableMessage sets a guild-specific message string within predefined messages.
func SetCustomizableMessage(gid string, key string, message string) error {
	_, ok := CustomizableMessages[key]
	if !ok {
		return &MessageOutOfBounds{}
	}
	err := setMessage(gid, key, message)
	return err
}

// GetTotalMessages gets the total amount of messages processed.
func GetTotalMessages() int {
	messageAmount, err := database.HGet("nitori", "total_messages")
	if err != nil {
		log.Warnf("Failed to obtain total amount of messages processed, %s", err)
		return 0
	}
	if messageAmount == "" {
		return 0
	}
	amountInteger, err := strconv.Atoi(messageAmount)
	if err != nil {
		log.Warnf("Malformed amount of messages processed, %s", err)
		return 0
	}
	return amountInteger
}

// AdvanceTotalMessages advances the total messages processed counter.
func AdvanceTotalMessages() error {
	return database.HSet("nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1))
}

// GetPrefix gets the command prefix of a guild and returns the default if none is set.
func GetPrefix(gid string) string {
	prefix, err := database.HGet("conf."+gid, "prefix")
	if err != nil {
		log.Warnf("Failed to obtain prefix in guild %s, %s", gid, err)
		return Config.System.Prefix
	}
	if prefix == "" {
		return Config.System.Prefix
	}
	return prefix
}

// SetPrefix sets the command prefix of a guild.
func SetPrefix(gid string, prefix string) error {
	return database.HSet("conf."+gid, "prefix", prefix)
}

// ResetPrefix resets the command prefix of a guild.
func ResetPrefix(gid string) error {
	return database.HDel("conf."+gid, "prefix")
}

// ExpEnabled queries whether the experience system is enabled for a guild.
func ExpEnabled(gid string) (enabled bool, err error) {
	result, err := database.HGet("conf."+gid, "exp_enable")
	if err != nil {
		return false, err
	}
	if result == "" {
		return false, nil
	}
	enabled, err = strconv.ParseBool(result)
	return
}

// ExpToggle toggles the experience system enabler.
func ExpToggle(gid string) (pre bool, err error) {
	pre, err = ExpEnabled(gid)
	switch pre {
	case true:
		err = database.HSet("conf."+gid, "exp_enable", "false")
	case false:
		err = database.HSet("conf."+gid, "exp_enable", "true")
	}
	return
}

// GetMemberExp obtains experience amount of a guild member.
func GetMemberExp(user *discordgo.User, guild *discordgo.Guild) (int, error) {
	result, err := database.HGet("exp."+guild.ID, user.ID)
	if err != nil {
		return 0, err
	}
	if result == "" {
		return 0, nil
	}
	return strconv.Atoi(result)
}

// SetMemberExp sets a member's experience amount.
func SetMemberExp(user *discordgo.User, guild *discordgo.Guild, exp int) error {
	return database.HSet("exp."+guild.ID, user.ID, strconv.Itoa(exp))
}

// GetLastfm gets a user's lastfm username.
func GetLastfm(user *discordgo.User, guild *discordgo.Guild) (string, error) {
	result, err := database.HGet("lastfm."+guild.ID, user.ID)
	if err != nil {
		return "", err
	}
	return result, err
}

// SetLastfm sets a user's lastfm username.
func SetLastfm(user *discordgo.User, guild *discordgo.Guild, username string) error {
	return database.HSet("lastfm."+guild.ID, user.ID, username)
}

// ResetLastfm resets a user's lastfm username.
func ResetLastfm(user *discordgo.User, guild *discordgo.Guild) error {
	return database.HDel("lastfm."+guild.ID, user.ID)
}
