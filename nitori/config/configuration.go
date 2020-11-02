package config

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/assets"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math"
	"os"
)

// Exported variables for usage in other classes
var Config = parseConfig()
var NitoriConfPath string
var TokenOverride string
var LogLevel = getLogLevel()
var CustomizableMessages = map[string]string{
	"levelup": "Congratulations $USER on reaching level $LEVEL.",
}

// Configuration related types
type MessageOutOfBounds struct{}
type Conf struct {
	System    SystemSection
	WebServer WebServerSection
	Discord   DiscordSection
	LastFM    LastFMSection
}
type SystemSection struct {
	LogLevel      string
	Socket        string
	Database      string
	Prefix        string
	ChatBackend   string
	WebServer     string
	Administrator int
	Operator      []int
}
type WebServerSection struct {
	SecretKey string
	Host      string
	Port      int
	BaseURL   string
}
type DiscordSection struct {
	Token        string
	ClientID     int
	ClientSecret string
	Presence     string
	Shard        bool
	ShardCount   int
}
type LastFMSection struct {
	ApiKey    string
	ApiSecret string
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
	var config []byte
	var err error
	if NitoriConfPath == "" {
		config, err = ioutil.ReadFile("/etc/freenitori/nitori.conf")
		if err != nil {
			config, err = ioutil.ReadFile("nitori.conf")
			if err != nil {
				log.Debugf("Configuration file inaccessible, %s", err)
				defaultConfigFile, err := assets.Asset("nitori.conf")
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
					"please edit it before starting FreeNitori.")
				os.Exit(1)
			} else {
				NitoriConfPath = "nitori.conf"
			}
		} else {
			NitoriConfPath = "/etc/nitori.conf"
		}
	} else {
		config, err = ioutil.ReadFile(NitoriConfPath)
		if err != nil {
			log.Fatalf("Unable to access configuration file, %s", err)
			os.Exit(1)
		}
	}
	if _, err := toml.Decode(string(config), &nitoriConf); err != nil {
		log.Fatalf("Configuration file syntax error, %s", err)
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
