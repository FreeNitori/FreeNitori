package config

import (
	"context"
	"encoding/base64"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

// Exported variables for usage in other classes
var Config = parseConfig()
var Administrator *discordgo.User
var Operator *discordgo.User
var Debug = getDebug()
var Redis = redis.NewClient(&redis.Options{
	Addr:     Config.Redis.Host + ":" + Config.Redis.Port,
	Password: Config.Redis.Password,
	DB:       Config.Redis.Database,
})
var RedisContext = context.Background()
var CustomizableMessages = map[string]string{
	"levelup": "Congratulations $USER on reaching level $LEVEL.",
}

// Configuration related types
type MessageOutOfBounds struct{}
type Conf struct {
	System    SystemSection
	Redis     RedisSection
	WebServer WebServerSection
}
type SystemSection struct {
	ExecutionMode string
	Socket        string
	Token         string
	Prefix        string
	Presence      string
	Shard         bool
	ShardCount    int
	Administrator string
	Operator      string
}
type RedisSection struct {
	Host     string
	Port     string
	Password string
	Database int
}
type WebServerSection struct {
	SecretKey string
	Host      string
	Port      string
	BaseURL   string
}

func init() {
	switch Debug {
	case true:
		log.SetLevel(logrus.DebugLevel)
	case false:
		log.SetLevel(logrus.InfoLevel)
	}
}

func (err MessageOutOfBounds) Error() string {
	return "message out of bounds"
}

// Parse or generate configuration file
func parseConfig() *Conf {
	var nitoriConf Conf
	config, err := ioutil.ReadFile("/etc/nitori.conf")
	if err != nil {
		config, err = ioutil.ReadFile("nitori.conf")
		if err != nil {
			log.Fatalf("Error loading configuration file, %s", err)
			defaultConfigFile, err := Asset("nitori.conf")
			if err != nil {
				log.Fatalf("Failed to extract the default configuration file, %s", err)
				os.Exit(1)
			}
			err = ioutil.WriteFile("nitori.conf", defaultConfigFile, 0644)
			if err != nil {
				log.Fatalf("Failed to write the default configuration file, %s", err)
				os.Exit(1)
			}
			log.Fatalf("Generated default configuration file at ./nitori.conf, " +
				"please edit it before restarting FreeNitori.")
			os.Exit(1)
		}
	}
	if _, err := toml.Decode(string(config), &nitoriConf); err != nil {
		log.Fatalf("Configuration file syntax error: %s", err)
		os.Exit(1)
	}
	return &nitoriConf
}

func ResetGuild(gid string) {
	Redis.Del(RedisContext, "settings."+gid)
	Redis.Del(RedisContext, "exp."+gid)
	Redis.Del(RedisContext, "rank."+gid)
	Redis.Del(RedisContext, "exp_bl."+gid)
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

// Figure out if the execution mode happens to be debug
func getDebug() bool {
	switch Config.System.ExecutionMode {
	case "debug":
		return true
	case "production":
		return false
	default:
		log.Fatalf("Unknown execution mode: %s", Config.System.ExecutionMode)
		os.Exit(1)
	}
	return false
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

// Add one message to the counter
func AddTotalMessages() error {
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

// Chat experience calculation stuffs
func LevelToExp(level int) int {
	return int(1000.0 * (math.Pow(float64(level), 1.25)))
}
func ExpToLevel(exp int) int {
	return int(math.Pow(float64(exp)/1000, 1.0/1.25))
}
