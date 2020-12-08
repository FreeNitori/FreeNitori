package config

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	dbVars "git.randomchars.net/RandomChars/FreeNitori/server/database/vars"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger/v2"
	"strconv"
)

var prefixes = []string{"conf", "exp", "rank", "exp_bl", "lastfm", "ra_metadata"}
var CustomizableMessages = map[string]string{
	"levelup": "Congratulations $USER on reaching level $LEVEL.",
}

// ResetGuild deletes all database values that belongs to a specific guild.
func ResetGuild(gid string) {
	for _, prefix := range prefixes {
		err := dbVars.Database.HDel(fmt.Sprintf("%s.%s", prefix, gid), []string{})
		if err != nil {
			if err == badger.ErrKeyNotFound {
				continue
			}
			log.Errorf("Error while resetting guild %s key %s, %s", gid, fmt.Sprintf("%s.%s", prefix, gid), err)
		}
	}
}

// getMessage gets a guild-specific string.
func getMessage(gid string, key string) (string, error) {
	message, err := dbVars.Database.HGet("conf."+gid, "message."+key)
	if err == badger.ErrKeyNotFound {
		return "", nil
	}
	return message, err
}

// setMessage sets a guild-specific string
func setMessage(gid string, key string, message string) error {
	if len(message) > 2048 {
		return &MessageOutOfBounds{}
	}
	if message == "" {
		return dbVars.Database.HDel("conf."+gid, []string{"message."+key})
	}
	return dbVars.Database.HSet("conf."+gid, "message."+key, message)
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
	messageAmount, err := dbVars.Database.HGet("nitori", "total_messages")
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 0
		}
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
	return dbVars.Database.HSet("nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1))
}

// GetPrefix gets the command prefix of a guild and returns the default if none is set.
func GetPrefix(gid string) string {
	prefix, err := dbVars.Database.HGet("conf."+gid, "prefix")
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return Config.System.Prefix
		}
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
	return dbVars.Database.HSet("conf."+gid, "prefix", prefix)
}

// ResetPrefix resets the command prefix of a guild.
func ResetPrefix(gid string) error {
	return dbVars.Database.HDel("conf."+gid, []string{"prefix"})
}

// ExpEnabled queries whether the experience system is enabled for a guild.
func ExpEnabled(gid string) (enabled bool, err error) {
	result, err := dbVars.Database.HGet("conf."+gid, "exp_enable")
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return false, nil
		}
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
		err = dbVars.Database.HSet("conf."+gid, "exp_enable", "false")
	case false:
		err = dbVars.Database.HSet("conf."+gid, "exp_enable", "true")
	}
	return
}

// GetMemberExp obtains experience amount of a guild member.
func GetMemberExp(user *discordgo.User, guild *discordgo.Guild) (int, error) {
	result, err := dbVars.Database.HGet("exp."+guild.ID, user.ID)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 0, nil
		}
		return 0, err
	}
	if result == "" {
		return 0, nil
	}
	return strconv.Atoi(result)
}

// SetMemberExp sets a member's experience amount.
func SetMemberExp(user *discordgo.User, guild *discordgo.Guild, exp int) error {
	return dbVars.Database.HSet("exp."+guild.ID, user.ID, strconv.Itoa(exp))
}

// GetLastfm gets a user's lastfm username.
func GetLastfm(user *discordgo.User, guild *discordgo.Guild) (string, error) {
	result, err := dbVars.Database.HGet("lastfm."+guild.ID, user.ID)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return "", nil
		}
		return "", err
	}
	return result, err
}

// SetLastfm sets a user's lastfm username.
func SetLastfm(user *discordgo.User, guild *discordgo.Guild, username string) error {
	return dbVars.Database.HSet("lastfm."+guild.ID, user.ID, username)
}

// ResetLastfm resets a user's lastfm username.
func ResetLastfm(user *discordgo.User, guild *discordgo.Guild) error {
	err := dbVars.Database.HDel("lastfm."+guild.ID, []string{user.ID})
	if err == badger.ErrKeyNotFound {
		return nil
	}
	return err
}
