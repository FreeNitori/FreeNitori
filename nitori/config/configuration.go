// Package config contains configuration related stuff.
package config

import (
	"flag"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"go/types"
	"io/ioutil"
	"os"
)

// Early initialization quirk
var (
	_ = flags()
)

var (
	// Config contains data loaded from configuration file.
	Config = parseConfig()
	// NitoriConfPath contains path to configuration file.
	NitoriConfPath string
	// TokenOverride contains the override token passed from command-line arguments.
	TokenOverride string
)

// MessageOutOfBounds represents an out of bounds message.
type MessageOutOfBounds struct{}

// Conf represents data of a configuration file.
type Conf struct {
	System struct {
		LogLevel       string
		LogPath        string
		Socket         string
		Database       string
		Prefix         string
		BackupInterval int
		Administrator  int
		Operator       []int
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
}

func flags() *types.Nil {
	var version bool
	flag.BoolVar(&version, "v", false, "Display version information and exit")
	flag.StringVar(&TokenOverride, "a", "", "Override Discord Authorization Token")
	flag.StringVar(&NitoriConfPath, "c", "", "Specify configuration file path")
	flag.Parse()

	if version {
		fmt.Printf("%s (%s)", state.Version(), state.Revision())
		os.Exit(0)
	}
	return nil
}

func (err MessageOutOfBounds) Error() string {
	return "message out of bounds"
}

// parseConfig parses the configuration file, generating one if none is present.
func parseConfig() *Conf {
	checkConfig()
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
	level, err := logrus.ParseLevel(nitoriConf.System.LogLevel)
	if err != nil {
		log.Fatalf("Unable to parse log level, %s", err)
		os.Exit(1)
	}
	log.SetLevel(level)
	return &nitoriConf
}

// checkConfig checks for a configuration file and generates default if not exists.
func checkConfig() {
	var nitoriConf = NitoriConfPath
	if NitoriConfPath == "" {
		nitoriConf = "nitori.conf"
	}
	if _, err := os.Stat(nitoriConf); os.IsNotExist(err) {
		file, err := os.Create(nitoriConf)
		if err != nil {
			log.Fatalf("Unable to create configuration file, %s", err)
			os.Exit(1)
		}
		encoder := toml.NewEncoder(file)
		err = encoder.Encode(confDefault)
		if err != nil {
			log.Fatalf("Unable to generate default configuration, %s", err)
			os.Exit(1)
		}
		log.Warnf("Generated default configuration file at %s, "+
			"please edit it before restarting FreeNitori.", nitoriConf)
		firstRun(false)
	}
}
