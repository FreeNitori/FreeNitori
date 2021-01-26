package discord

import (
	"errors"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/server/discord/vars"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"plugin"
	"strconv"
	"strings"
)

var err error

// Initialize prepares Discord-related functionality
func Initialize() error {

	// Load plugins if not window
	if !state.IsWindow() {
		err = loadPlugins()
		if err != nil {
			return err
		}
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
	_, _ = vars.RawSession.UserUpdateStatus("dnd")
	_ = vars.RawSession.UpdateStatus(0, config.Config.Discord.Presence)
	if config.Config.Discord.Shard {
		err = MakeSessions()
		if err != nil {
			return err
		}
	}
	log.Infof("Nitori has successfully logged in as %s#%s (%s)",
		vars.RawSession.State.User.Username,
		vars.RawSession.State.User.Discriminator,
		vars.RawSession.State.User.ID)
	return nil
}

func loadPlugins() error {
	stat, err := os.Stat("plugins")
	if os.IsNotExist(err) {
		return errors.New("plugins directory does not exist")
	}
	if !stat.IsDir() {
		return errors.New("plugins path is not a directory")
	}
	pluginPaths, err := ioutil.ReadDir("plugins/")
	if err != nil {
		return errors.New("plugins directory unreadable")
	}
	for _, path := range pluginPaths {
		if !strings.HasSuffix(path.Name(), ".so") {
			continue
		}
		pl, err := plugin.Open("plugins/" + path.Name())
		if err != nil {
			log.Warnf("Error while loading plugin %s, %s", path.Name(), err)
			continue
		}
		symbol, err := pl.Lookup("CommandRoute")
		if err != nil {
			continue
		}
		route, ok := symbol.(*multiplexer.Route)
		if !ok {
			log.Warnf("No Route found in %s.", path.Name())
			continue
		}
		multiplexer.Router.Route(route)
		log.Infof("Loaded plugin %s implementing command %s.", path.Name(), route.Pattern)
	}
	return nil
}
