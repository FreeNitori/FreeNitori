package main

import (
	"fmt"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/args"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/ipc"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	_ "git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/state"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"os/signal"
	"plugin"
	"strconv"
	"strings"
	"syscall"
)

var err error

func init() {
	vars.ProcessType = vars.ChatBackend
	func() {
		stat, err := os.Stat("plugins")
		if os.IsNotExist(err) {
			err = os.Mkdir("plugins", 0755)
			if err != nil {
				log.Fatalf("Failed to create plugin directory, %s", err)
				_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
				os.Exit(1)
			}
			return
		}
		if !stat.IsDir() {
			log.Fatal("Plugin path is not a directory.")
			_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
		}
		pluginPaths, err := ioutil.ReadDir("plugins/")
		if err != nil {
			log.Fatalf("Unable to read plugin directory, %s", err)
			_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
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
				log.Warnf("Error while looking up CommandRoute symbol in plugin %s, %s", path.Name(), err)
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
	}()
}

func main() {
	// Connect to the Supervisor
	err = ipc.InitializeIPC()
	if err != nil {
		log.Fatalf("Failed to connect to the supervisor process, %s", err)
		os.Exit(1)
	}
	defer func() { _ = vars.RPCConnection.Close() }()

	// Authenticate and make session
	discordgo.Logger = log.DiscordGoLogger
	state.RawSession.UserAgent = "DiscordBot (FreeNitori " + vars.Version + ")"
	if config.TokenOverride == "" {
		state.RawSession.Token = "Bot " + config.Config.Discord.Token
	} else {
		state.RawSession.Token = "Bot " + config.TokenOverride
	}
	state.RawSession.ShouldReconnectOnError = true
	state.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	err = state.RawSession.Open()
	if err != nil {
		state.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
		err = state.RawSession.Open()
		if err != nil {
			log.Error(fmt.Sprintf("An error occurred while connecting to Discord, %s", err))
			_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
		}
	}
	state.Administrator, err = state.RawSession.User(strconv.Itoa(config.Config.System.Administrator))
	if err != nil {
		log.Fatalf("Failed to get system administrator, %s", err)
		_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
		os.Exit(1)
	}
	for _, id := range config.Config.System.Operator {
		user, err := state.RawSession.User(strconv.Itoa(id))
		if err == nil {
			state.Operator = append(state.Operator, user)
		}
	}
	vars.Initialized = true
	state.Application, err = state.RawSession.Application("@me")
	vars.InviteURL = fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=2146958847", state.Application.ID)
	if err != nil {
		log.Error(fmt.Sprintf("An error occurred while fetching application info, %s", err))
		_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
		os.Exit(1)
	}
	_, _ = state.RawSession.UserUpdateStatus("dnd")
	_ = state.RawSession.UpdateStatus(0, config.Config.Discord.Presence)
	if config.Config.Discord.Shard {
		err = MakeSessions()
		if err != nil {
			_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
		}
	}

	// Fire the ready message and signal the WebServer
	_ = vars.RPCConnection.Call("R.FireReadyMessage", []string{
		state.RawSession.State.User.Username + "#" + state.RawSession.State.User.Discriminator,
		state.RawSession.State.User.ID}, nil)
	_ = vars.RPCConnection.Call("R.SignalWebServer", []string{}, nil)

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt, os.Kill)
	go func() {
		for {
			currentSignal := <-signalChannel
			switch currentSignal {
			case syscall.SIGUSR1:
				// Go to the supervisor to fetch further instructions
				ChatBackendIPCReceiver()
			case syscall.SIGUSR2:
				vars.ExitCode <- 0
				break
			default:
				// Cleanup stuffs
				if currentSignal != os.Interrupt {
					// Only tell the supervisor if SIGUSR2 was not sent or the program was not interrupted
					_ = vars.RPCConnection.Call("R.Restart", []string{"ChatBackend"}, nil)
				}
				for _, shardSession := range state.ShardSessions {
					_ = shardSession.Close()
				}
				_ = state.RawSession.Close()
				vars.ExitCode <- 0
				break
			}
		}
	}()

	// Tell the Supervisor and exit if there's something on that channel
	exitCode := <-vars.ExitCode
	if exitCode != 0 {
		_ = vars.RPCConnection.Call("R.Error", []string{"ChatBackend"}, nil)
	}
	os.Exit(exitCode)

}
