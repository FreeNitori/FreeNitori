package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/communication"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/session"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	ChatBackend "git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var err error

func init() {
	// Add the multiplexer handler to the raw session if sharding is disabled
	if !config.Config.Discord.Shard {
		for _, handler := range ChatBackend.EventHandlers {
			ChatBackend.RawSession.AddHandler(handler)
		}
	}

	// Add the event handlers
	for _, handlerInfo := range multiplexer.Commands {
		multiplexer.Router.Route(
			handlerInfo.Pattern,
			handlerInfo.AliasPatterns,
			handlerInfo.Description,
			handlerInfo.Handler,
			handlerInfo.Category)
	}
}

func main() {
	state.ProcessType = state.ChatBackend

	// Connect to the Supervisor
	err = communication.InitializeIPC()
	if err != nil {
		log.Fatalf("Failed to connect to the supervisor process, %s", err)
		os.Exit(1)
	}

	// Authenticate and make session
	if ChatBackend.RawSession.Token == "" {
		configToken := config.Config.Discord.Token
		if configToken != "" && configToken != "INSERT_TOKEN_HERE" {
			log.Debug("Loaded token from configuration file.")
			ChatBackend.RawSession.Token = configToken
		} else {
			log.Error("Please specify an authorization token.")
			_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
		}
	} else {
		log.Debug("Loaded token from command parameter.")
	}
	discordgo.Logger = log.DiscordGoLogger
	ChatBackend.RawSession.UserAgent = "DiscordBot (FreeNitori " + state.Version + ")"
	ChatBackend.RawSession.Token = "Bot " + ChatBackend.RawSession.Token
	ChatBackend.RawSession.ShouldReconnectOnError = true
	ChatBackend.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	err = ChatBackend.RawSession.Open()
	if err != nil {
		ChatBackend.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
		err = ChatBackend.RawSession.Open()
		if err != nil {
			log.Error(fmt.Sprintf("An error occurred while connecting to Discord, %s", err))
			_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
		}
	}
	ChatBackend.Administrator, err = ChatBackend.RawSession.User(strconv.Itoa(config.Config.System.Administrator))
	if err != nil {
		log.Fatalf("Failed to get system administrator, %s", err)
		_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
		os.Exit(1)
	}
	for _, id := range config.Config.System.Operator {
		user, err := ChatBackend.RawSession.User(strconv.Itoa(id))
		if err == nil {
			ChatBackend.Operator = append(ChatBackend.Operator, user)
		}
	}
	state.Initialized = true
	ChatBackend.Application, err = ChatBackend.RawSession.Application("@me")
	state.InviteURL = fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=2146958847", ChatBackend.Application.ID)
	if err != nil {
		log.Error(fmt.Sprintf("An error occurred while fetching application info, %s", err))
		_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
		os.Exit(1)
	}
	_, _ = ChatBackend.RawSession.UserUpdateStatus("dnd")
	_ = ChatBackend.RawSession.UpdateStatus(0, config.Config.Discord.Presence)
	if config.Config.Discord.Shard {
		err = session.MakeSessions()
		if err != nil {
			_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
			os.Exit(1)
		}
	}

	// Fire the ready message and signal the WebServer
	_ = state.IPCConnection.Call("IPC.FireReadyMessage", []string{
		ChatBackend.RawSession.State.User.Username + "#" + ChatBackend.RawSession.State.User.Discriminator,
		ChatBackend.RawSession.State.User.ID}, nil)
	_ = state.IPCConnection.Call("IPC.SignalWebServer", []string{}, nil)

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt, os.Kill)
	go func() {
		for {
			currentSignal := <-signalChannel
			switch currentSignal {
			case syscall.SIGUSR1:
				// Go to the supervisor to fetch further instructions
				communication.ChatBackendIPCReceiver()
			case syscall.SIGUSR2:
				state.ExitCode <- 0
				return
			default:
				// Cleanup stuffs
				if currentSignal != os.Interrupt {
					// Only tell the supervisor if SIGUSR2 was not sent or the program was not interrupted
					_ = state.IPCConnection.Call("IPC.Restart", []string{"ChatBackend"}, nil)
				}
				for _, shardSession := range ChatBackend.ShardSessions {
					_ = shardSession.Close()
				}
				_ = ChatBackend.RawSession.Close()
				state.ExitCode <- 0
				return
			}
		}
	}()

	// Tell the Supervisor and exit if there's something on that channel
	exitCode := <-state.ExitCode
	if exitCode != 0 {
		_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
	}
	os.Exit(exitCode)

}
