// Functions to interact with global and guild-specific configuration values.
package config

import (
	"flag"
	"git.randomchars.net/RandomChars/FreeNitori/binaries/confdefault"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"go/types"
	"io/ioutil"
	"math"
	"os"
)

// Early initialization quirk
var _ = flags()
var _ = checkConfig()
var _ = setLogLevel()

// Exported variables for usage in other classes
var (
	Config         = parseConfig()
	NitoriConfPath string
	TokenOverride  string
	VersionStartup bool
	LogLevel       = getLogLevel()
)

// Configuration related types
type MessageOutOfBounds struct{}
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
		Token        string
		ClientSecret string
		Presence     string
		Shard        bool
		ShardCount   int
	}
	LastFM struct {
		ApiKey    string
		ApiSecret string
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
		select {}
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

// LevelToExp calculates amount of experience from a level integer.
func LevelToExp(level int) int {
	return int(1000.0 * (math.Pow(float64(level), 1.25)))
}

// ExpToLevel calculates amount of levels from an experience integer.
func ExpToLevel(exp int) int {
	return int(math.Pow(float64(exp)/1000, 1.0/1.25))
}
