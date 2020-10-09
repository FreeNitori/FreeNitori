package main

import (
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/communication"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	_ "git.randomchars.net/RandomChars/FreeNitori/nitori/handlers"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/log"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/session"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state/chatbackend"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state/supervisor"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/web"
	"github.com/bwmarrin/discordgo"
	"github.com/dgraph-io/badger/v2"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	// Some regular initialization
	var err error
	var readyChannel = make(chan bool, 1)
	var IPCFunctions = new(communication.IPC)
	switch {
	case state.StartChatBackend && state.StartWebServer:
		{

			// This doesn't work, so exit
			println("Parameter \"-cb\" cannot be used with \"-ws\".")
			os.Exit(1)
		}
	case state.StartChatBackend:
		{
			// Dial the supervisor socket
			err = communication.InitializeIPC(state.StartChatBackend, state.StartWebServer)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				log.Error(fmt.Sprintf("Failed to connect to the database, %s", err))
				os.Exit(1)
			}

			// Authenticate and make session
			if ChatBackend.RawSession.Token == "" {
				configToken := config.Config.System.Token
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
					os.Exit(1)
				}
			}
			ChatBackend.Administrator, err = ChatBackend.RawSession.User(strconv.Itoa(config.Config.System.Administrator))
			if err != nil {
				log.Fatalf("Failed to get system administrator, %s", err)
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
				os.Exit(1)
			}
			_, _ = ChatBackend.RawSession.UserUpdateStatus("dnd")
			_ = ChatBackend.RawSession.UpdateStatus(0, config.Config.System.Presence)
			if config.Config.System.Shard {
				err = session.MakeSessions()
				if err != nil {
					_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
					os.Exit(1)
				}
			}

			// Print out logs ChatBackend is ready to go
			_ = state.IPCConnection.Call("IPC.FireReadyMessage", []string{
				ChatBackend.RawSession.State.User.Username + "#" + ChatBackend.RawSession.State.User.Discriminator,
				ChatBackend.RawSession.State.User.ID}, nil)
			_ = state.IPCConnection.Call("IPC.SignalWebServer", []string{}, nil)
		}
	case state.StartWebServer:
		{

			// Dial the supervisor socket
			err = communication.InitializeIPC(state.StartChatBackend, state.StartWebServer)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				log.Error(fmt.Sprintf("Unable to establish state with database, %s", err))
				os.Exit(1)
			}

			// Initialize and start the server
			web.Initialize()
			go func() {
				<-readyChannel
				err = web.Engine.Run(fmt.Sprintf("%s:%s", config.Config.WebServer.Host, strconv.Itoa(config.Config.WebServer.Port)))
				if err != nil {
					log.Error(fmt.Sprintf("Failed to start web server, %s", err))
					_ = state.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
					os.Exit(1)
				}
			}()
		}
	case !state.StartWebServer && !state.StartChatBackend:
		{
			// Print version information and stuff
			log.Infof("Starting FreeNitori %s", state.Version)

			// Check for an existing instance
			if _, err := os.Stat(config.Config.System.Socket); os.IsNotExist(err) {
			} else {
				_, err := net.Dial("unix", config.Config.System.Socket)
				if err != nil {
					err = syscall.Unlink(config.Config.System.Socket)
					if err != nil {
						log.Error(fmt.Sprintf("Unable to remove hanging socket, %s", err))
						os.Exit(1)
					}
				} else {
					log.Error("Another instance of FreeNitori is already running.")
					os.Exit(1)
				}
			}

			// Initialize the socket
			_ = rpc.Register(IPCFunctions)
			rpc.HandleHTTP()
			SuperVisor.SocketListener, err = net.Listen("unix", config.Config.System.Socket)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to listen on the socket, %s", err))
				os.Exit(1)
			}
			go http.Serve(SuperVisor.SocketListener, nil)

			// Open the database
			dbOptions := badger.DefaultOptions(config.Config.System.Database)
			dbOptions.Logger = log.Logger
			SuperVisor.Database, err = badger.Open(dbOptions)
			if err != nil {
				log.Fatalf("Failed to open database, %s", err)
				os.Exit(1)
			}
			defer func() { _ = SuperVisor.Database.Close() }()

			// Create the chat backend process
			SuperVisor.ChatBackendProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-cb", "-a", ChatBackend.RawSession.Token, "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to create chat backend process, %s", err))
				os.Exit(1)
			}

			// Create web server process
			SuperVisor.WebServerProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-ws", "-c", config.NitoriConfPath}, &SuperVisor.ProcessAttributes)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to create web server process, %s", err))
				os.Exit(1)
			}

		}
	}

	// Signal handling
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2, os.Interrupt, os.Kill)
	go func() {
		for {
			currentSignal := <-signalChannel
			switch currentSignal {
			case syscall.SIGUSR1:
				// Go to the supervisor to fetch further instructions
				if state.StartChatBackend && !state.StartWebServer {
					communication.ChatBackendIPCReceiver()
				} else if state.StartWebServer && !state.StartChatBackend {
					if !state.Initialized {
						readyChannel <- true
						state.Initialized = true
					}
				}
			case syscall.SIGUSR2:
				state.ExitCode <- 0
				return
			default:
				// Cleanup stuffs
				if !state.StartChatBackend && !state.StartWebServer {
					fmt.Print("\n")
					log.Info("Gracefully terminating...")
					_ = SuperVisor.ChatBackendProcess.Signal(syscall.SIGUSR2)
					_ = SuperVisor.WebServerProcess.Signal(syscall.SIGUSR2)
					_ = SuperVisor.SocketListener.Close()
					_ = syscall.Unlink(config.Config.System.Socket)
				} else if state.StartChatBackend {
					if currentSignal != os.Interrupt {
						// Only tell the supervisor if SIGUSR2 was not sent or the program was not interrupted
						_ = state.IPCConnection.Call("IPC.Restart", []string{"ChatBackend"}, nil)
					}
					for _, shardSession := range ChatBackend.ShardSessions {
						_ = shardSession.Close()
					}
					_ = ChatBackend.RawSession.Close()
				} else if state.StartWebServer {
					if currentSignal != os.Interrupt {
						// Only write the packet if SIGUSR2 was not sent or the program was not interrupted
						_ = state.IPCConnection.Call("IPC.Restart", []string{"WebServer"}, nil)
					}
				}
				state.ExitCode <- 0
				return
			}
		}
	}()

	// Tell the Supervisor and exit if there's something on that channel
	exitCode := <-state.ExitCode
	if state.StartChatBackend && !state.StartWebServer && exitCode != 0 {
		_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
	} else if state.StartWebServer && !state.StartChatBackend && exitCode != 0 {
		_ = state.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
	}
	os.Exit(exitCode)
}
