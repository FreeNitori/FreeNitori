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

// Exported variables for usage in other classes
var Config = getConfig()
var SocketPath = Config.Section("System").Key("Socket").String()
var Prefix = Config.Section("System").Key("Prefix").String()
var Shard, _ = Config.Section("System").Key("Shard").Bool()
var ShardCount, _ = Config.Section("System").Key("ShardCount").Int()
var Presence = Config.Section("System").Key("Presence").String()
var Administrator = Config.Section("System").Key("Administrator").String()
var Operator = Config.Section("System").Key("Operator").String()
var BaseURL = Config.Section("WebServer").Key("BaseURL").String()
var Host = Config.Section("WebServer").Key("Host").String()
var Port = Config.Section("WebServer").Key("Port").String()
var Debug = getDebug()
var Redis = getDatabaseClient()
var RedisContext = context.Background()

// Fetch configuration file object and generate one if needed
func getConfig() (Config *ini.File) {
	Config, err := ini.Load("/etc/nitori.conf")
	if err != nil {
		Config, err = ini.Load("nitori.conf")
		if err != nil {
			log.Printf("Error loading configuration file, %s", err)
			defaultConfigFile, err := Asset("nitori.conf")
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

// Obtain a redis client using details stored in the configuration
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

// Figure out if the execution mode happens to be debug
func getDebug() (debug bool) {
	executionMode := Config.Section("System").Key("ExecutionMode").String()
	var debugMode bool
	switch {
	case executionMode == "debug":
		debugMode = true
		break
	case executionMode == "production":
		debugMode = false
		break
	case true:
		log.Printf("Unknown execution mode: %s", executionMode)
		os.Exit(1)
	}
	return debugMode
}

// Get amount of messages totally processed
func GetTotalMessages() (amount int) {
	var err error
	messageAmount, err := Redis.HGet(RedisContext, "nitori", "total_messages").Result()
	if err != nil {
		if err == redis.Nil {
			return 0
		}
		log.Printf("Failed to obtain total amount of messages processed, %s", err)
		return 0
	}
	if messageAmount == "" {
		return 0
	}
	amountInteger, err := strconv.Atoi(messageAmount)
	if err != nil {
		log.Printf("Malformed amount of messages processed, %s", err)
		return 0
	}
	return amountInteger
}

// Add one message to the counter
func AddTotalMessages() (err error) {
	return Redis.HSet(RedisContext, "nitori", "total_messages", strconv.Itoa(GetTotalMessages()+1)).Err()
}

// Get prefix for a guild and return the default if there is none
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

// Set the prefix of a guild
func SetPrefix(gid int, prefix string) (err error) {
	return Redis.HSet(RedisContext, "settings."+strconv.Itoa(gid), "prefix", prefix).Err()
}
