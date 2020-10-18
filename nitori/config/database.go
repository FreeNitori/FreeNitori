package config

import (
	"encoding/base64"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// Completely reset a specific guild's configuration
func ResetGuild(gid string) {
	Redis.Del(RedisContext, "settings."+gid)
	Redis.Del(RedisContext, "exp."+gid)
	Redis.Del(RedisContext, "rank."+gid)
	Redis.Del(RedisContext, "exp_bl."+gid)
	Redis.Del(RedisContext, "lastfm."+gid)
	Redis.Del(RedisContext, "ra_metadata."+gid)
	Redis.Del(RedisContext, "ra_table_0."+gid)
	Redis.Del(RedisContext, "ra_table_1."+gid)
	Redis.Del(RedisContext, "ra_table_2."+gid)
	Redis.Del(RedisContext, "ra_table_3."+gid)
	Redis.Del(RedisContext, "ra_table_4."+gid)
	Redis.Del(RedisContext, "ra_table_5."+gid)
	Redis.Del(RedisContext, "ra_table_6."+gid)
	Redis.Del(RedisContext, "ra_table_7."+gid)
}

// Get a guild-specific message string
func getMessage(gid string, key string) (string, error) {
	messageEncoded, err := database.HGet("settings."+gid, "message."+key)
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		log.Warnf("Failed to obtain message in guild %s, %s", gid, err)
		return "", err
	}
	if messageEncoded == "" {
		return "", nil
	}
	message, err := base64.StdEncoding.DecodeString(messageEncoded)
	if err != nil {
		log.Warnf("Malformed message in guild %s, %s", gid, err)
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
		err := Redis.HDel(RedisContext, "settings."+gid, "message."+key).Err()
		return err
	}
	messageEncoded := base64.StdEncoding.EncodeToString([]byte(message))
	err := Redis.HSet(RedisContext, "settings."+gid, "message."+key, messageEncoded).Err()
	return err
}

// Get amount of messages totally processed
func GetTotalMessages() int {
	messageAmount, err := Redis.HGet(RedisContext, "nitori", "total_messages").Result()
	if err != nil {
		if err == redis.Nil {
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

// Advance the counter once
func AdvanceTotalMessages() error {
	return Redis.HSet(RedisContext, "nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1)).Err()
}

// Get prefix for a guild and return the default if there is none
func GetPrefix(gid string) string {
	prefixValue, err := Redis.HGet(RedisContext, "settings."+gid, "prefix").Result()
	if err != nil {
		if err == redis.Nil {
			return Config.System.Prefix
		}
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
	return Redis.HSet(RedisContext, "settings."+gid, "prefix", prefixEncoded).Err()
}

// Reset the prefix of a guild
func ResetPrefix(gid string) error {
	return Redis.HDel(RedisContext, "settings."+gid, "prefix").Err()
}

// Figure out if experience system is enabled
func ExpEnabled(gid string) (enabled bool, err error) {
	result, err := Redis.HGet(RedisContext, "settings."+gid, "exp_enable").Result()
	if err != nil {
		if err == redis.Nil {
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

// Toggle the experience system enabler
func ExpToggle(gid string) (pre bool, err error) {
	pre, err = ExpEnabled(gid)
	switch pre {
	case true:
		err = Redis.HSet(RedisContext, "settings."+gid, "exp_enable", "false").Err()
	case false:
		err = Redis.HSet(RedisContext, "settings."+gid, "exp_enable", "true").Err()
	}
	return
}

// Obtain experience amount of a guild member
func GetMemberExp(user *discordgo.User, guild *discordgo.Guild) (int, error) {
	result, err := Redis.HGet(RedisContext, "exp."+guild.ID, user.ID).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	if result == "" {
		return 0, nil
	}
	return strconv.Atoi(result)
}

// Set a member's experience amount
func SetMemberExp(user *discordgo.User, guild *discordgo.Guild, exp int) error {
	return Redis.HSet(RedisContext, "exp."+guild.ID, user.ID, strconv.Itoa(exp)).Err()
}

// Get a user's lastfm username
func GetLastfm(user *discordgo.User, guild *discordgo.Guild) (string, error) {
	result, err := Redis.HGet(RedisContext, "lastfm."+guild.ID, user.ID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return result, err
}

// Set a user's lastfm username
func SetLastfm(user *discordgo.User, guild *discordgo.Guild, username string) error {
	return Redis.HSet(RedisContext, "lastfm."+guild.ID, user.ID, username).Err()
}

// Reset a user's lastfm username
func ResetLastfm(user *discordgo.User, guild *discordgo.Guild) error {
	return Redis.HDel(RedisContext, "lastfm."+guild.ID, user.ID).Err()
}
