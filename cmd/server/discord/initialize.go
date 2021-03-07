package discord

import (
	"errors"
	"fmt"
	"git.randomchars.net/FreeNitori/FreeNitori/cmd/server/discord/sessioning"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/config"
	"git.randomchars.net/FreeNitori/FreeNitori/nitori/state"
	log "git.randomchars.net/FreeNitori/Log"
	multiplexer "git.randomchars.net/FreeNitori/Multiplexer"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strconv"
)

var err error

// Initialize early initializes Discord-related functionalities.
func Initialize() error {
	// Setup some things
	multiplexer.NoCommandMatched = func(context *multiplexer.Context) {
		// If no command was matched, resort to either being annoyed by the ping or a command not found message
		if context.HasMention {
			context.SendMessage("<a:KyokoAngryPing:757399059114885180>")
		} else {
			context.SendMessage(fmt.Sprintf("This command does not exist! Issue `%sman` for a list of command manuals.",
				context.Prefix()))
		}
	}
	state.Multiplexer.Prefix = config.Config.System.Prefix
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
	state.RawSession.UserAgent = "DiscordBot (FreeNitori " + state.Version() + ")"
	if config.TokenOverride == "" {
		state.RawSession.Token = "Bot " + config.Config.Discord.Token
	} else {
		state.RawSession.Token = "Bot " + config.TokenOverride
	}
	state.RawSession.ShouldReconnectOnError = true
	state.RawSession.State.MaxMessageCount = config.Config.Discord.CachePerChannel
	state.RawSession.Identify.Intents = discordgo.IntentsAll

	return nil
}

// LateInitialize late initializes Discord-related features.
func LateInitialize() error {
	// Authenticate and make session
	err = state.RawSession.Open()
	if err != nil {
		log.Warnf("Unable to open session with all intents, %s, Nitori will now fallback to unprivileged intents, some functionality will be unavailable.", err)
		state.RawSession.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
		err = state.RawSession.Open()
		if err != nil {
			return errors.New("unable to open session with Discord")
		}
	}
	log.Info("Raw session with Discord opened.")
	state.Multiplexer.Administrator, err = state.RawSession.User(strconv.Itoa(config.Config.System.Administrator))
	if err != nil {
		return errors.New("unable to get system administrator")
	}
	for _, id := range config.Config.System.Operator {
		user, err := state.RawSession.User(strconv.Itoa(id))
		if err == nil {
			state.Multiplexer.Operator = append(state.Multiplexer.Operator, user)
		}
	}
	state.Application, err = state.RawSession.Application("@me")
	if err != nil {
		return errors.New("unable to fetch application information")
	}
	state.InviteURL = fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8", state.Application.ID)
	go func() {
		for {
			state.DiscordReady <- true
		}
	}()

	if config.Config.Discord.Shard {
		log.Infof("Sharding is enabled, starting %v shards.", config.Config.Discord.ShardCount)
		err = sessioning.MakeSessions()
		if err != nil {
			return err
		}
	}
	log.Infof("Nitori has successfully logged in as %s#%s (%s).",
		state.RawSession.State.User.Username,
		state.RawSession.State.User.Discriminator,
		state.RawSession.State.User.ID)
	return nil
}
