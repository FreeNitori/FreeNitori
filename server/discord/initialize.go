package discord

import (
	"errors"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

var err error

// Initialize prepares Discord-related functionality
func Initialize() error {

	// Load plugins if not window
	err = loadPlugins()
	if err != nil {
		return err
	}

	// Setup some things
	discordgo.Logger = log.DiscordGoLogger
	vars.RawSession.UserAgent = "DiscordBot (FreeNitori " + state.Version() + ")"
	if config.TokenOverride == "" {
		vars.RawSession.Token = "Bot " + config.Config.Discord.Token
	} else {
		vars.RawSession.Token = "Bot " + config.TokenOverride
	}
	vars.RawSession.ShouldReconnectOnError = true
	vars.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	return nil
}

func LateInitialize() error {
	// Authenticate and make session
	err = vars.RawSession.Open()
	if err != nil {
		vars.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
		err = vars.RawSession.Open()
		if err != nil {
			return errors.New("unable to open session with Discord")
		}
	}
	log.Info("Raw session with Discord opened.")
	vars.Administrator, err = vars.RawSession.User(strconv.Itoa(config.Config.System.Administrator))
	if err != nil {
		return errors.New("unable to get system administrator")
	}
	for _, id := range config.Config.System.Operator {
		user, err := vars.RawSession.User(strconv.Itoa(id))
		if err == nil {
			vars.Operator = append(vars.Operator, user)
		}
	}
	vars.Application, err = vars.RawSession.Application("@me")
	if err != nil {
		return errors.New("unable to fetch application information")
	}
	state.InviteURL = fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=2146958847", vars.Application.ID)
	go func() {
		for {
			state.DiscordReady <- true
		}
	}()

	if config.Config.Discord.Shard {
		err = MakeSessions()
		if err != nil {
			return err
		}
	}
	log.Infof("Nitori has successfully logged in as %s#%s (%s).",
		vars.RawSession.State.User.Username,
		vars.RawSession.State.User.Discriminator,
		vars.RawSession.State.User.ID)
	return nil
}
