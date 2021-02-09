package db

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/database"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger/v2"
	"strconv"
)

func init() {
	config.Prefixes = append(config.Prefixes, "exp", "rank", "exp_bl", "lastfm", "ra_metadata", "highlight")
	config.CustomizableMessages["levelup"] = "Congratulations $USER on reaching level $LEVEL."
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

// GetRankBinds gets binding of a role to a level.
func GetRankBinds(guild *discordgo.Guild) (map[string]string, error) {
	result, err := database.Database.HGetAll("rank." + guild.ID)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, err
	}
	return result, err
}

// GetRankBind gets binding of a role to a level.
func GetRankBind(guild *discordgo.Guild, level int) (string, error) {
	result, err := database.Database.HGet("rank."+guild.ID, strconv.Itoa(level))
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return "", nil
		}
		return "", err
	}
	return result, err
}

// SetRankBind binds a role to a level.
func SetRankBind(guild *discordgo.Guild, level int, role *discordgo.Role) error {
	return database.Database.HSet("rank."+guild.ID, strconv.Itoa(level), role.ID)
}

// UnsetRankBind unbinds a role from a level.
func UnsetRankBind(guild *discordgo.Guild, level string) error {
	err := database.Database.HDel("rank."+guild.ID, []string{level})
	if err == badger.ErrKeyNotFound {
		return nil
	}
	return err
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