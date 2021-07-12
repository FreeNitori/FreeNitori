package discord

import (
	"flag"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/config"
	log "git.randomchars.net/FreeNitori/Log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"

	// Run all init functions from internals.
	_ "git.randomchars.net/FreeNitori/FreeNitori/cmd/freenitori/discord/routes"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
)

var token string

func init() {
	flag.StringVar(&token, "a", "", "Override configured authorization token")
}

func Open() error {
	// Add the multiplexer handler to the no shard session if sharding is disabled
	if !config.Discord.Shard {
		state.Multiplexer.SessionRegisterHandlers(state.Session)
	}

	// Setup no route handler
	multiplexer.NoCommandMatched = func(context *multiplexer.Context) {
		if context.HasMention {
			context.SendMessage("<a:KyokoAngryPing:757399059114885180>")
		} else {
			context.SendMessage(fmt.Sprintf("This command does not exist! Issue `%sman` for a list of command manuals.",
				context.Prefix()))
		}
	}

	// Default prefix
	state.Multiplexer.Prefix = config.System.Prefix

	// Discordgo logger
	discordgo.Logger = func(msgL, _ int, format string, a ...interface{}) {
		var level logrus.Level
		switch msgL {
		case discordgo.LogDebug:
			level = logrus.DebugLevel
		case discordgo.LogInformational:
			level = logrus.InfoLevel
		case discordgo.LogWarning:
			level = logrus.WarnLevel
		case discordgo.LogError:
			level = logrus.ErrorLevel
		}
		log.Instance.Log(level, fmt.Sprintf(format, a...))
	}

	// Setup direct session
	state.Session.UserAgent = "DiscordBot (FreeNitori " + state.Version() + ")"
	if token == "" {
		state.Session.Token = "Bot " + config.Discord.Token
	} else {
		state.Session.Token = "Bot " + token
	}
	state.Session.ShouldReconnectOnError = true
	state.Session.State.MaxMessageCount = config.Discord.CachePerChannel
	state.Session.Identify.Intents = discordgo.IntentsAll

	// Open session
	if err := state.Session.Open(); err != nil {
		log.Warnf("Error opening session with all intents, %s", err)
		log.Warn("Nitori will fallback to unprivileged intents, some functionality will be unavailable.")
		state.Session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
		err = state.Session.Open()
		if err != nil {
			return err
		}
	}
	log.Info("Discord session opened successfully.")

	// Get administrator
	if administrator, err := state.Session.User(strconv.Itoa(config.System.Administrator)); err != nil {
		return err
	} else {
		state.Multiplexer.Administrator = administrator
	}

	// Get operators
	for _, id := range config.System.Operator {
		user, err := state.Session.User(strconv.Itoa(id))
		if err == nil {
			state.Multiplexer.Operator = append(state.Multiplexer.Operator, user)
		}
	}

	// Get application
	if application, err := state.Session.Application("@me"); err != nil {
		return err
	} else {
		state.Application = application
	}

	// Set invite URL
	state.InviteURL = fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8", state.Application.ID)

	// Open shards
	if config.Discord.Shard {
		log.Infof("Starting %v shards.", config.Discord.ShardCount)

		// Get recommended shard count from Discord as required
		if config.Discord.ShardCount < 1 {
			gatewayBot, err := state.Session.GatewayBot()
			if err != nil {
				return err
			}
			config.Discord.ShardCount = gatewayBot.Shards
		}

		// Make sure more than 0 shards
		if config.Discord.ShardCount == 0 {
			config.Discord.ShardCount = 1
		}

		// Make and open sessions
		for i := 0; i < config.Discord.ShardCount; i++ {
			time.Sleep(time.Millisecond * 100)

			// Setup session
			session, _ := discordgo.New()
			session.ShardCount = config.Discord.ShardCount
			session.ShardID = i
			session.Token = state.Session.Token
			session.UserAgent = state.Session.UserAgent
			session.ShouldReconnectOnError = state.Session.ShouldReconnectOnError
			session.Identify.Intents = state.Session.Identify.Intents
			session.State.MaxMessageCount = state.Session.State.MaxMessageCount

			// Open session
			if err := session.Open(); err != nil {
				return err
			}

			// Register handlers
			state.Multiplexer.SessionRegisterHandlers(session)

			// Add session
			state.ShardSessions = append(state.ShardSessions, session)

			log.Infof("Shard %s ready.", strconv.Itoa(i))
		}
	}

	// Perform reincarnation actions
	for _, variable := range os.Environ() {
		if strings.HasPrefix(variable, "REINCARNATION=") {
			log.Infof("Reincarnation payload found: %s", variable)
			if err := os.Unsetenv("REINCARNATION"); err != nil {
				log.Errorf("Error unsetting reincarnation variable, %s", err)
			}
			split := strings.Split(variable[14:], "\t")
			if len(split) != 3 {
				log.Error("Reincarnation payload has incorrect format.")
				break
			}
			if _, err := state.Session.ChannelMessageEdit(split[0], split[1], split[2]); err != nil {
				log.Errorf("Error editing message of previous incarnation, %s", err)
				break
			}
			break
		}
	}

	// Print final initialisation message
	log.Infof("Nitori has logged in as %s#%s (%s).",
		state.Session.State.User.Username,
		state.Session.State.User.Discriminator,
		state.Session.State.User.ID)

	return nil
}

func Close() error {
	// Close shard sessions
	for index, shardSession := range state.ShardSessions {
		if err := shardSession.Close(); err != nil {
			log.Errorf("Error closing shard %v, %s", index, err)
		}
	}

	// CLose direct session
	if err := state.Session.Close(); err != nil {
		return err
	}

	return nil
}
