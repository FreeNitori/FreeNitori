package config

import (
	"errors"
	"flag"
	log "git.randomchars.net/FreeNitori/Log"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	path string
	read bool
)

var (
	ErrAlreadyRead = errors.New("config already read")
	ErrFirstRun    = errors.New("first run")
)

var (
	System    system
	WebServer ws
	Discord   discord
)

func init() {
	flag.StringVar(&path, "c", "nitori.conf", "Specify path to configuration file")
}

// ReadConfig reads config if not yet read.
func ReadConfig() error {
	// Check if already read
	if read {
		return ErrAlreadyRead
	}
	defer func() { read = true }()

	// Decode config
	var payload conf
	if meta, err := toml.DecodeFile(path, &payload); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// Generate config file as required
		var file *os.File
		if file, err = os.Create(path); err != nil {
			return err
		}
		if err = toml.NewEncoder(file).Encode(def); err != nil {
			return err
		}
		first(false)
		return ErrFirstRun
	} else {
		// Warn about unused keys
		for _, key := range meta.Undecoded() {
			log.Warnf("Unused key in config: %s", key.String())
		}
	}

	// Fill out the previous variables
	System = payload.System
	WebServer = payload.WebServer
	Discord = payload.Discord

	// Parse amd set loglevel
	if level, err := logrus.ParseLevel(payload.System.LogLevel); err != nil {
		return err
	} else {
		log.SetLevel(level)
	}

	return nil
}

func CheckPlaceholders() {
	if Discord.Token == "INSERT_TOKEN_HERE" ||
		Discord.ClientSecret == "INSERT_CLIENT_SECRET_HERE" ||
		WebServer.Secret == "RANDOM_STRING" {
		log.Fatal("Configuration file still contains initial placeholders.")
		first(true)
		os.Exit(1)
	}
}
