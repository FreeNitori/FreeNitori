package config

import (
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/database"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/dgraph-io/badger/v3"
	"strconv"
)

// Prefixes holds a slice of strings that prefixes guild configuration hashmaps.
var Prefixes = []string{"conf"}

// CustomizableMessages maps customizable messages to their default values.
var CustomizableMessages = map[string]string{}

// ResetGuild deletes all db values that belongs to a specific guild.
func ResetGuild(gid string) {
	for _, prefix := range Prefixes {
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

// GetTotalMessages gets the message counter.
func GetTotalMessages() int {
	messageAmount, err := database.Database.HGet("nitori", "total_messages")
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return 0
		}
		log.Warnf("Unable to obtain message counter, %s", err)
		return 0
	}
	if messageAmount == "" {
		return 0
	}
	amountInteger, err := strconv.Atoi(messageAmount)
	if err != nil {
		log.Warnf("Malformed message counter, %s", err)
		return 0
	}
	return amountInteger
}

// AdvanceTotalMessages advances the message counter.
func AdvanceTotalMessages() error {
	return database.Database.HSet("nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1))
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
