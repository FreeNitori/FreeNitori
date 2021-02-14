package discord

import (
	"errors"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/cmd/server/discord/sessioning"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var err error

// Initialize early initializes Discord-related functionalities.
func Initialize() error {

	// Load plugins if not window
	err = loadPlugins()
	if err != nil {
		return err
	}

	// Setup some things
	discordgo.Logger = log.DiscordGoLogger
	state.RawSession.UserAgent = "DiscordBot (FreeNitori " + state.Version() + ")"
	if config.TokenOverride == "" {
		state.RawSession.Token = "Bot " + config.Config.Discord.Token
	} else {
		state.RawSession.Token = "Bot " + config.TokenOverride
	}
	state.RawSession.ShouldReconnectOnError = true
	state.RawSession.State.MaxMessageCount = config.Config.Discord.CachePerChannel
	state.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	return nil
}

// LateInitialize late initializes Discord-related features.
func LateInitialize() error {
	// Authenticate and make session
	err = state.RawSession.Open()
	if err != nil {
		log.Warnf("Unable to open session with all intents, %s, Nitori will now fallback to unprivileged intents, some functionality will be unavailable.", err)
		state.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
		err = state.RawSession.Open()
		if err != nil {
			return errors.New("unable to open session with Discord")
		}
	}
	log.Info("Raw session with Discord opened.")
	state.Administrator, err = state.RawSession.User(strconv.Itoa(config.Config.System.Administrator))
	if err != nil {
		return errors.New("unable to get system administrator")
	}
	for _, id := range config.Config.System.Operator {
		user, err := state.RawSession.User(strconv.Itoa(id))
		if err == nil {
			state.Operator = append(state.Operator, user)
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
