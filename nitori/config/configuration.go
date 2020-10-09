package config

import (
	"context"
	"flag"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go/types"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

var _ = flags()

func flags() *types.Nil {
	flag.StringVar(&state.RawSession.Token, "a", "", "Discord Authorization Token")
	flag.StringVar(&NitoriConfPath, "c", "", "Specify configuration file path.")
	flag.BoolVar(&state.StartChatBackend, "cb", false, "Start the chat backend directly")
	flag.BoolVar(&state.StartWebServer, "ws", false, "Start the web server directly")
	flag.Parse()
	return nil
}

// Exported variables for usage in other classes
var Config = parseConfig()
var NitoriConfPath string
var LogLevel = getLogLevel()
var Redis = redis.NewClient(&redis.Options{
	Addr:     Config.Redis.Host + ":" + strconv.Itoa(Config.Redis.Port),
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
	LastFM    LastFMSection
}
type SystemSection struct {
	LogLevel      string
	Socket        string
	Token         string
	ClientID      int
	ClientSecret  string
	Prefix        string
	Presence      string
	Shard         bool
	ShardCount    int
	Administrator int
	Operator      []int
}
type RedisSection struct {
	Host     string
	Port     int
	Password string
	Database int
}
type WebServerSection struct {
	SecretKey string
	Host      string
	Port      int
	BaseURL   string
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
		config, err = ioutil.ReadFile("/etc/nitori.conf")
		if err != nil {
			config, err = ioutil.ReadFile("nitori.conf")
			if err != nil {
				log.Debugf("Configuration file inaccessible, %s", err)
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
