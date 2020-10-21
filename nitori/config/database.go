package config

import (
	"encoding/base64"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var prefixes = []string{"settings", "exp", "rank", "exp_bl", "lastfm", "ra_metadata"}

// Completely reset a specific guild's configuration
func ResetGuild(gid string) {
	for _, prefix := range prefixes {
		err := database.HDel(fmt.Sprintf("%s.%s", prefix, gid))
		if err != nil {
			log.Errorf("Error while resetting guild %s key %s, %s", gid, fmt.Sprintf("%s.%s", prefix, gid), err)
		}
	}
}

// Get a guild-specific message string
func getMessage(gid string, key string) (string, error) {
	messageEncoded, err := database.HGet("settings."+gid, "message."+key)
	if err != nil {
		return "", err
	}
	if messageEncoded == "" {
		return "", nil
	}
	message, err := base64.StdEncoding.DecodeString(messageEncoded)
	if err != nil {
		return "", err
	}
	return string(message), nil
}

// Set a guild-specific message string
func setMessage(gid string, key string, message string) error {
	if len(message) > 2048 {
		return &MessageOutOfBounds{}
	}
	if message == "" {
		err := database.HDel("settings."+gid, "message."+key)
		return err
	}
	messageEncoded := base64.StdEncoding.EncodeToString([]byte(message))
	err := database.HSet("settings."+gid, "message."+key, messageEncoded)
	return err
}

// Get amount of messages totally processed
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

// Advance the counter once
func AdvanceTotalMessages() error {
	return database.HSet("nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1))
}

// Get prefix for a guild and return the default if there is none
func GetPrefix(gid string) string {
	prefixValue, err := database.HGet("settings."+gid, "prefix")
	if err != nil {
		log.Warnf("Failed to obtain prefix in guild %s, %s", gid, err)
		return Config.System.Prefix
	}
	if prefixValue == "" {
		return Config.System.Prefix
	}
	prefixDecoded, err := base64.StdEncoding.DecodeString(prefixValue)
	if err != nil {
		log.Warnf("Malformed prefix in guild %s, %s", gid, err)
		return Config.System.Prefix
	}
	return string(prefixDecoded)
}

// Set the prefix of a guild
func SetPrefix(gid string, prefix string) error {
	prefixEncoded := base64.StdEncoding.EncodeToString([]byte(prefix))
	return database.HSet("settings."+gid, "prefix", prefixEncoded)
}

// Reset the prefix of a guild
func ResetPrefix(gid string) error {
	return database.HDel("settings."+gid, "prefix")
}

// Figure out if experience system is enabled
func ExpEnabled(gid string) (enabled bool, err error) {
	result, err := database.HGet("settings."+gid, "exp_enable")
	if err != nil {
		return false, err
	}
	if result == "" {
		return false, nil
	}
	enabled, err = strconv.ParseBool(result)
	return
}

// Toggle the experience system enabler
func ExpToggle(gid string) (pre bool, err error) {
	pre, err = ExpEnabled(gid)
	switch pre {
	case true:
		err = database.HSet("settings."+gid, "exp_enable", "false")
	case false:
		err = database.HSet("settings."+gid, "exp_enable", "true")
	}
	return
}

// Obtain experience amount of a guild member
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

// Set a member's experience amount
func SetMemberExp(user *discordgo.User, guild *discordgo.Guild, exp int) error {
	return database.HSet("exp."+guild.ID, user.ID, strconv.Itoa(exp))
}

// Get a user's lastfm username
func GetLastfm(user *discordgo.User, guild *discordgo.Guild) (string, error) {
	result, err := database.HGet("lastfm."+guild.ID, user.ID)
	if err != nil {
		return "", err
	}
	return result, err
}

// Set a user's lastfm username
func SetLastfm(user *discordgo.User, guild *discordgo.Guild, username string) error {
	return database.HSet("lastfm."+guild.ID, user.ID, username)
}

// Reset a user's lastfm username
func ResetLastfm(user *discordgo.User, guild *discordgo.Guild) error {
	return database.HDel("lastfm."+guild.ID, user.ID)
}
