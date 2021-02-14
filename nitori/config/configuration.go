// Package config contains configuration related stuff.
package config

import (
	"flag"
	"git.randomchars.net/FreeNitori/FreeNitori/binaries/confdefault"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/log"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"go/types"
	"io/ioutil"
	"os"
)

// Early initialization quirk
var (
	_ = flags()
	_ = checkConfig()
	_ = setLogLevel()
)

var (
	// Config contains data loaded from configuration file.
	Config = parseConfig()
	// NitoriConfPath contains path to configuration file.
	NitoriConfPath string
	// TokenOverride contains the override token passed from command-line arguments.
	TokenOverride string
	// VersionStartup indicates weather the program should display version information and exit.
	VersionStartup bool
	// LogLevel is the log level.
	LogLevel = getLogLevel()
)

// MessageOutOfBounds represents an out of bounds message.
type MessageOutOfBounds struct{}

// Conf represents data of a configuration file.
type Conf struct {
	System struct {
		LogLevel      string
		LogPath       string
		Socket        string
		Database      string
		Prefix        string
		Administrator int
		Operator      []int
	}
	WebServer struct {
		Host                string
		Port                int
		BaseURL             string
		Unix                bool
		ForwardedByClientIP bool
		Secret              string
		RateLimit           int
		RateLimitPeriod     int
	}
	Discord struct {
		Token           string
		ClientSecret    string
		Presence        string
		Shard           bool
		ShardCount      int
		CachePerChannel int
	}
	LastFM struct {
		APIKey    string
		APISecret string
	}
}

func flags() *types.Nil {
	flag.BoolVar(&VersionStartup, "v", false, "Display Version information and exit")
	flag.StringVar(&TokenOverride, "a", "", "Override Discord Authorization Token")
	flag.StringVar(&NitoriConfPath, "c", "", "Specify configuration file path")
	flag.Parse()
	return nil
}

func setLogLevel() *types.Nil {
	log.SetLevel(LogLevel)
	return nil
}

func (err MessageOutOfBounds) Error() string {
	return "message out of bounds"
}

// parseConfig parses the configuration file, generating one if none is present.
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
				return nil
			}
			NitoriConfPath = "nitori.conf"
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
	if nitoriConf.Discord.Token == "INSERT_TOKEN_HERE" || nitoriConf.Discord.ClientSecret == "INSERT_CLIENT_SECRET_HERE" || nitoriConf.WebServer.Secret == "RANDOM_STRING" {
		log.Warn("Please edit the configuration file before starting.")
		firstRun(true)
	}
	return &nitoriConf
}

// checkConfig checks for a configuration file and generates default if not exists.
func checkConfig() *types.Nil {
	var nitoriConf = NitoriConfPath
	if NitoriConfPath == "" {
		nitoriConf = "nitori.conf"
	}
	if _, err := os.Stat(nitoriConf); os.IsNotExist(err) {
		defaultConfigFile, err := confdefault.Asset("nitori.conf")
		if err != nil {
			log.Fatalf("Failed to extract the default configuration file, %s", err)
			os.Exit(1)
		}
		err = ioutil.WriteFile(nitoriConf, defaultConfigFile, 0644)
		if err != nil {
			log.Fatalf("Failed to write the default configuration file, %s", err)
			os.Exit(1)
		}
		log.Warnf("Generated default configuration file at %s, "+
			"please edit it before restarting FreeNitori.", nitoriConf)
		firstRun(false)
	}
	return nil
}

// getLogLevel refers the log level configuration string to a log level integer.
func getLogLevel() logrus.Level {
	level, err := logrus.ParseLevel(Config.System.LogLevel)
	if err != nil {
		log.Fatalf("Unable to parse log level, %s", err)
		os.Exit(1)
	}
	return level
}
