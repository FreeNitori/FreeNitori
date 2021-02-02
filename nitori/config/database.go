package config

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger/v2"
	"strconv"
)

var prefixes = []string{"conf", "exp", "rank", "exp_bl", "lastfm", "ra_metadata", "highlight"}
var CustomizableMessages = map[string]string{
	"levelup": "Congratulations $USER on reaching level $LEVEL.",
}

// ResetGuild deletes all db values that belongs to a specific guild.
func ResetGuild(gid string) {
	for _, prefix := range prefixes {
		err := database.Database.HDel(fmt.Sprintf("%s.%s", prefix, gid), []string{})
		if err != nil {
			if err == badger.ErrKeyNotFound {
				continue
			}
			log.Errorf("Error while resetting guild %s key %s, %s", gid, fmt.Sprintf("%s.%s", prefix, gid), err)
		}
	}
}

// ResetGuildMap deletes a map that belongs to a specific guild.
func ResetGuildMap(gid, key string) {
	err := database.Database.HDel(fmt.Sprintf("%s.%s", key, gid), []string{})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return
		}
		log.Errorf("Error while resetting guild %s key %s, %s", gid, fmt.Sprintf("%s.%s", key, gid), err)
	}

}

// getMessage gets a guild-specific string.
func getMessage(gid string, key string) (string, error) {
	message, err := database.Database.HGet("conf."+gid, "message."+key)
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
		return database.Database.HDel("conf."+gid, []string{"message." + key})
	}
	return database.Database.HSet("conf."+gid, "message."+key, message)
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
	messageAmount, err := database.Database.HGet("nitori", "total_messages")
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
	return database.Database.HSet("nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1))
}

// GetPrefix gets the command prefix of a guild and returns the default if none is set.
func GetPrefix(gid string) string {
	prefix, err := database.Database.HGet("conf."+gid, "prefix")
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

// ExpEnabled queries whether the experience system is enabled for a guild.
func ExpEnabled(gid string) (enabled bool, err error) {
	result, err := database.Database.HGet("conf."+gid, "exp_enable")
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

// HighlightBindMessage binds a message with the highlight message.
func HighlightBindMessage(gid, message, highlight string) error {
	return database.Database.HSet("highlight."+gid, message, highlight)
}

// HighlightUnbindMessage unbinds a message with the highlight message.
func HighlightUnbindMessage(gid, message string) error {
	err := database.Database.HDel("highlight."+gid, []string{message})
	if err == badger.ErrKeyNotFound {
		return nil
	}
	return err
}

// HighlightGetBinding gets the binding of a message.
func HighlightGetBinding(gid, message string) (string, error) {
	value, err := database.Database.HGet("highlight."+gid, message)
	if err == badger.ErrKeyNotFound {
		return "", nil
	}
	return value, err
}

// GetGuildConfValue gets a configuration value for a specific guild
func GetGuildConfValue(id, key string) (string, error) {
	result, err := database.Database.HGet("conf."+id, key)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return "", nil
		}
		return "", err
	}
	return result, nil
}

// SetGuildConfValue sets a configuration value for a specific guild
func SetGuildConfValue(id, key, value string) error {
	return database.Database.HSet("conf."+id, key, value)
}

// ResetGuildConfValue resets a configuration value for a specific guild
func ResetGuildConfValue(id, key string) error {
	err := database.Database.HDel("conf."+id, []string{key})
	if err == badger.ErrKeyNotFound {
		return nil
	}
	return err
}

// GetMemberExp obtains experience amount of a guild member.
func GetMemberExp(user *discordgo.User, guild *discordgo.Guild) (int, error) {
	result, err := database.Database.HGet("exp."+guild.ID, user.ID)
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
	return database.Database.HSet("exp."+guild.ID, user.ID, strconv.Itoa(exp))
}

// GetLastfm gets a user's lastfm username.
func GetLastfm(user *discordgo.User, guild *discordgo.Guild) (string, error) {
	result, err := database.Database.HGet("lastfm."+guild.ID, user.ID)
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
	return database.Database.HSet("lastfm."+guild.ID, user.ID, username)
}

// ResetLastfm resets a user's lastfm username.
func ResetLastfm(user *discordgo.User, guild *discordgo.Guild) error {
	err := database.Database.HDel("lastfm."+guild.ID, []string{user.ID})
	if err == badger.ErrKeyNotFound {
		return nil
	}
	return err
}
