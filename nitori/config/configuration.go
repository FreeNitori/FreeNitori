package config

import (
	"context"
	"encoding/base64"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
)

// Exported variables for usage in other classes
var Config = getConfig()
var SocketPath = Config.Section("System").Key("Socket").String()
var Prefix = Config.Section("System").Key("Prefix").String()
var Shard, _ = Config.Section("System").Key("Shard").Bool()
var ShardCount, _ = Config.Section("System").Key("ShardCount").Int()
var Presence = Config.Section("System").Key("Presence").String()
var Administrator = Config.Section("System").Key("Administrator").String()
var Operator = Config.Section("System").Key("Operator").String()
var BaseURL = Config.Section("WebServer").Key("BaseURL").String()
var Host = Config.Section("WebServer").Key("Host").String()
var Port = Config.Section("WebServer").Key("Port").String()
var Debug = getDebug()
var Redis = getDatabaseClient()
var RedisContext = context.Background()
var CustomizableMessages = map[string]string{
	"levelup": "Congratulations $USER on reaching level $LEVEL.",
}

type MessageOutOfBounds struct{}

func (err MessageOutOfBounds) Error() string {
	return "message out of bounds"
}

// Fetch configuration file object and generate one if needed
func getConfig() *ini.File {
	config, err := ini.Load("/etc/nitori.conf")
	if err != nil {
		config, err = ini.Load("nitori.conf")
		if err != nil {
			log.Printf("Error loading configuration file, %s", err)
			defaultConfigFile, err := Asset("nitori.conf")
			if err != nil {
				log.Printf("Failed to extract the default configuration file, %s", err)
				os.Exit(1)
			}
			err = ioutil.WriteFile("nitori.conf", defaultConfigFile, 0644)
			if err != nil {
				log.Printf("Failed to write the default configuration file, %s", err)
				os.Exit(1)
			}
			log.Println("Generated default configuration file at ./nitori.conf, " +
				"please edit it before restarting FreeNitori.")
			os.Exit(1)
		}
		if config.Section("System").Key("ExecutionMode").String() == "debug" {
			log.Println("Loaded configuration file from current directory.")
		}
	} else {
		if config.Section("System").Key("ExecutionMode").String() == "debug" {
			log.Println("Loaded system-wide configuration from /etc/nitori.conf.")
		}
	}
	return config
}

// Obtain a redis client using details stored in the configuration
func getDatabaseClient() *redis.Client {
	db, err := strconv.Atoi(Config.Section("Redis").Key("Database").String())
	if err != nil {
		log.Printf("Failed to read redis database configuration, %s", err)
		os.Exit(1)
	}
	return redis.NewClient(&redis.Options{
		Addr: Config.Section("Redis").Key("Host").String() +
			":" + Config.Section("Redis").Key("Port").String(),
		Password: Config.Section("Redis").Key("Password").String(),
		DB:       db,
	})
}

// Figure out if the execution mode happens to be debug
func getDebug() bool {
	executionMode := Config.Section("System").Key("ExecutionMode").String()
	var debugMode bool
	switch {
	case executionMode == "debug":
		debugMode = true
		break
	case executionMode == "production":
		debugMode = false
		break
	case true:
		log.Printf("Unknown execution mode: %s", executionMode)
		os.Exit(1)
	}
	return debugMode
}

// Get a guild-specific message string within predefined messages
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

// Set a guild-specific message string within predefined messages
func SetCustomizableMessage(gid string, key string, message string) error {
	_, ok := CustomizableMessages[key]
	if !ok {
		return &MessageOutOfBounds{}
	}
	err := setMessage(gid, key, message)
	return err
}

// Get a guild-specific message string
func getMessage(gid string, key string) (string, error) {
	messageEncoded, err := Redis.HGet(RedisContext, "settings."+gid, "message."+key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		log.Printf("Failed to obtain message in guild %s, %s", gid, err)
		return "", err
	}
	if messageEncoded == "" {
		return "", nil
	}
	message, err := base64.StdEncoding.DecodeString(messageEncoded)
	if err != nil {
		log.Printf("Malformed message in guild %s, %s", gid, err)
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
		log.Printf("Failed to obtain total amount of messages processed, %s", err)
		return 0
	}
	if messageAmount == "" {
		return 0
	}
	amountInteger, err := strconv.Atoi(messageAmount)
	if err != nil {
		log.Printf("Malformed amount of messages processed, %s", err)
		return 0
	}
	return amountInteger
}

// Add one message to the counter
func AddTotalMessages() error {
	return Redis.HSet(RedisContext, "nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1)).Err()
}

// Get prefix for a guild and return the default if there is none
func GetPrefix(gid string) string {
	prefixValue, err := Redis.HGet(RedisContext, "settings."+gid, "prefix").Result()
	if err != nil {
		if err == redis.Nil {
			return Prefix
		}
		log.Printf("Failed to obtain prefix in guild %s, %s", gid, err)
		return Prefix
	}
	if prefixValue == "" {
		return Prefix
	}
	prefixDecoded, err := base64.StdEncoding.DecodeString(prefixValue)
	if err != nil {
		log.Printf("Malformed prefix in guild %s, %s", gid, err)
		return Prefix
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

// Chat experience calculation stuffs
func LevelToExp(level int) int {
	return int(1000.0 * (math.Pow(float64(level), 1.25)))
}
func ExpToLevel(exp int) int {
	return int(math.Pow(float64(exp)/1000, 1.0/1.25))
}
