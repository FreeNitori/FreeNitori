package main

import (
	"flag"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const Version = "v0.0.1-rewrite"

var Session, _ = discordgo.New()

func init() {
	flag.StringVar(&Session.Token, "t", "", "Discord Authorization Token")
}

func main() {
	// Some regular initialization
	var err error
	fmt.Printf(
		`
___________                      _______  .__  __               .__ 
\_   _____/______   ____   ____  \      \ |__|/  |_  ___________|__|
 |    __) \_  __ \_/ __ \_/ __ \ /   |   \|  \   __\/  _ \_  __ \  |
 |     \   |  | \/\  ___/\  ___//    |    \  ||  | (  <_> )  | \/  |
 \___  /   |__|    \___  >\___  >____|__  /__||__|  \____/|__|  |__|
     \/                \/     \/        \/    %-16s
`+"\n\n", Version)
	flag.Parse()

	// Check the database
	_, err = config.Redis.Ping(config.RedisContext).Result()
	if err != nil {
		log.Printf("Unable to establish connection with database, %s", err)
		os.Exit(1)
	}

	// Authenticate and make session
	if Session.Token == "" {
		configToken := config.Config.Section("System").Key("Token").String()
		if configToken != "" && configToken != "INSERT_TOKEN_HERE" {
			if config.Debug {
				log.Println("Loaded token from configuration file.")
			}
			Session.Token = configToken
		} else {
			log.Println("Please specify an authorization token.")
			os.Exit(1)
		}
	} else {
		if config.Debug {
			log.Println("Loaded token from command parameter.")
		}
	}
	Session.UserAgent = "DiscordBot (FreeNitori " + Version + ")"
	Session.Token = "Bot " + Session.Token
	err = Session.Open()
	if err != nil {
		log.Printf("An error occurred while connecting to Discord, %s \n", err)
		os.Exit(1)
	}

	// Regular running and signal handling
	log.Printf("User: %s | ID: %s | Prefix: %s",
		Session.State.User.Username+"#"+Session.State.User.Discriminator,
		Session.State.User.ID,
		config.Prefix)
	log.Printf("FreeNitori is now running. Press Control-C to terminate.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanup stuffs
	fmt.Print("\n")
	log.Println("Gracefully terminating...")
	_ = Session.Close()
}
