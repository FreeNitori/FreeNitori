package config

import (
	"context"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
)

// Exported variables for usage in other classes
var Config = parseConfig()
var Administrator *discordgo.User
var Operator *discordgo.User
var LogLevel = getLogLevel()
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
	LogLevel      string
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
	log.SetLevel(LogLevel)
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

// Set the log level
func getLogLevel() logrus.Level {
	switch Config.System.LogLevel {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		log.Fatalf("Unknown log level \"%s\"", Config.System.LogLevel)
		os.Exit(1)
	}
	return logrus.InfoLevel
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

// Chat experience calculation stuffs
func LevelToExp(level int) int {
	return int(1000.0 * (math.Pow(float64(level), 1.25)))
}
func ExpToLevel(exp int) int {
	return int(math.Pow(float64(exp)/1000, 1.0/1.25))
}
