package config

import (
	"context"
	"encoding/base64"
	"github.com/go-redis/redis/v8"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var Config = getConfig()
var Prefix = Config.Section("System").Key("Prefix").String()
var Debug = getDebug()
var Redis = getDatabaseClient()
var RedisContext = context.Background()

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

func getDatabaseClient() (client *redis.Client) {
	var err error
	db, err := strconv.Atoi(Config.Section("Redis").Key("Database").String())
	if err != nil {
		log.Printf("Failed to read redis database configuration, %s", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: Config.Section("Redis").Key("Host").String() +
			":" + Config.Section("Redis").Key("Port").String(),
		Password: Config.Section("Redis").Key("Password").String(),
		DB:       db,
	})
	return redisClient
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

func GetPrefix(gid int) (prefix string) {
	var err error
	prefixValue, err := Redis.HGet(RedisContext, "settings."+strconv.Itoa(gid), "prefix").Result()
	if err != nil {
		if err == redis.Nil {
			return Prefix
		}
		log.Printf("Failed to obtain prefix in guild %s, %s", strconv.Itoa(gid), err)
		return Prefix
	}
	if prefixValue == "" {
		return Prefix
	}
	prefixDecoded, err := base64.StdEncoding.DecodeString(prefixValue)
	if err != nil {
		log.Printf("Malformed prefix in guild %s, %s", strconv.Itoa(gid), err)
		return Prefix
	}
	return string(prefixDecoded)
}

func SetPrefix(gid int, prefix string) (err error) {
	return Redis.HSet(RedisContext, "settings."+strconv.Itoa(gid), "prefix", prefix).Err()
}
