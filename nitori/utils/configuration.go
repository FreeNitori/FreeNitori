package utils

import (
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
)

var Config = getConfig()
var Prefix = Config.Section("System").Key("Prefix").String()
var Debug = getDebug()

func getConfig() (Config *ini.File) {
	// Parse configuration file and generate default
	Config, err := ini.Load("/etc/nitori.conf")
	if err != nil {
		Config, err = ini.Load("nitori.conf")
		if err != nil {
			log.Printf("Error loading configuration file, %s", err)
			defaultConfigFile, err := Asset("assets/nitori.conf")
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
		if Config.Section("System").Key("ExecutionMode").String() == "debug" {
			log.Println("Loaded configuration file from current directory.")
		}
	} else {
		if Config.Section("System").Key("ExecutionMode").String() == "debug" {
			log.Println("Loaded system-wide configuration from /etc/nitori.conf.")
		}
	}
	return Config
}

func getDebug() (debug bool) {
	if Config.Section("System").Key("ExecutionMode").String() == "debug" {
		return true
	} else if Config.Section("System").Key("ExecutionMode").String() == "production" {
		return false
	} else {
		log.Printf("Unknown execution mode: %s", Config.Section("System").Key("ExecutionMode").String())
		os.Exit(1)
	}
	return false
}
