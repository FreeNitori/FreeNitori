package main

import (
	"flag"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/communication"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/session"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/web"
	"github.com/bwmarrin/discordgo"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

var StartChatBackend bool
var StartWebServer bool

func init() {
	flag.StringVar(&state.RawSession.Token, "a", "", "Discord Authorization Token")
	flag.BoolVar(&StartChatBackend, "c", false, "Start the chat backend directly")
	flag.BoolVar(&StartWebServer, "w", false, "Start the web server directly")
}

func main() {
	// Some regular initialization
	var err error
	var readyChannel = make(chan bool, 1)
	var SocketListener net.Listener
	var IPCFunctions = new(communication.IPC)
	flag.Parse()
	switch {
	case StartChatBackend && StartWebServer:
		{

			// This doesn't work, so exit
			println("Parameter \"-c\" cannot be used with \"-w\".")
			os.Exit(1)
		}
	case StartChatBackend:
		{

			// Dial the supervisor socket
			err = communication.InitializeIPC(StartChatBackend, StartWebServer)
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Failed to connect to the database, %s", err))
				os.Exit(1)
			}

			// Authenticate and make session
			if state.RawSession.Token == "" {
				configToken := config.Config.Section("System").Key("Token").String()
				if configToken != "" && configToken != "INSERT_TOKEN_HERE" {
					if config.Debug {
						state.Logger.Debug("Loaded token from configuration file.")
					}
					state.RawSession.Token = configToken
				} else {
					state.Logger.Error("Please specify an authorization token.")
					_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
					os.Exit(1)
				}
			} else {
				if config.Debug {
					state.Logger.Error("Loaded token from command parameter.")
				}
			}

			state.RawSession.UserAgent = "DiscordBot (FreeNitori " + state.Version + ")"
			state.RawSession.Token = "Bot " + state.RawSession.Token
			state.RawSession.ShouldReconnectOnError = true
			state.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
			err = state.RawSession.Open()
			if err != nil {
				state.RawSession.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)
				err = state.RawSession.Open()
				if err != nil {
					state.Logger.Error(fmt.Sprintf("An error occurred while connecting to Discord, %s", err))
					os.Exit(1)
				}
			}
			state.Initialized = true
			state.Application, err = state.RawSession.Application("@me")
			if err != nil {
				state.Logger.Error(fmt.Sprintf("An error occurred while fetching application info, %s", err))
				os.Exit(1)
			}
			_, _ = state.RawSession.UserUpdateStatus("dnd")
			_ = state.RawSession.UpdateStatus(0, config.Presence)
			if config.Shard {
				err = session.MakeSessions()
				if err != nil {
					_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
					os.Exit(1)
				}
			}

			// Log into the logger that the ChatBackend is ready to go
			_ = state.IPCConnection.Call("IPC.Log", []string{
				"INFO",
				fmt.Sprintf("User: %s | ID: %s | Default Prefix: %s",
					state.RawSession.State.User.Username+"#"+state.RawSession.State.User.Discriminator,
					state.RawSession.State.User.ID,
					config.Prefix),
			}, nil)
			_ = state.IPCConnection.Call("IPC.Log", []string{
				"INFO",
				"FreeNitori is now ready. Press Control-C to terminate.",
			}, nil)
			_ = state.IPCConnection.Call("IPC.SignalWebServer", []string{}, nil)
		}
	case StartWebServer:
		{

			// Dial the supervisor socket
			err = communication.InitializeIPC(StartChatBackend, StartWebServer)
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Unable to establish state with database, %s", err))
				os.Exit(1)
			}

			// Initialize and start the server
			web.Initialize()
			go func() {
				<-readyChannel
				err = web.Engine.Run(fmt.Sprintf("%s:%s", config.Host, config.Port))
				if err != nil {
					state.Logger.Error(fmt.Sprintf("Failed to start web server, %s", err))
					_ = state.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
					os.Exit(1)
				}
			}()
		}
	case !StartWebServer && !StartChatBackend:
		{

			// Print some fancy ASCII art
			fmt.Printf(
				`
___________                      _______  .__  __               .__ 
\_   _____/______   ____   ____  \      \ |__|/  |_  ___________|__|
 |    __) \_  __ \_/ __ \_/ __ \ /   |   \|  \   __\/  _ \_  __ \  |
 |     \   |  | \/\  ___/\  ___//    |    \  ||  | (  <_> )  | \/  |
 \___  /   |__|    \___  >\___  >____|__  /__||__|  \____/|__|  |__|
     \/                \/     \/        \/    %-16s
`+"\n", state.Version)

			// Check for an existing instance
			if _, err := os.Stat(config.SocketPath); os.IsNotExist(err) {
			} else {
				_, err := net.Dial("unix", config.SocketPath)
				if err != nil {
					err = syscall.Unlink(config.SocketPath)
					if err != nil {
						state.Logger.Error(fmt.Sprintf("Unable to remove hanging socket, %s", err))
						os.Exit(1)
					}
				} else {
					state.Logger.Error("Another instance of FreeNitori is already running.")
					os.Exit(1)
				}
			}

			// Initialize the socket
			_ = rpc.Register(IPCFunctions)
			rpc.HandleHTTP()
			SocketListener, err = net.Listen("unix", config.SocketPath)
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Failed to listen on the socket, %s", err))
				os.Exit(1)
			}
			go http.Serve(SocketListener, nil)

			// Create the chat backend process
			state.ChatBackendProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-c", "-a", state.RawSession.Token}, &state.ProcessAttributes)
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Failed to create chat backend process, %s", err))
				os.Exit(1)
			}

			// Create web server process
			state.WebServerProcess, err =
				os.StartProcess(state.ExecPath, []string{state.ExecPath, "-w"}, &state.ProcessAttributes)
			if err != nil {
				state.Logger.Error(fmt.Sprintf("Failed to create web server process, %s", err))
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
				if StartChatBackend && !StartWebServer {
					communication.ChatBackendIPCReceiver()
				} else if StartWebServer && !StartChatBackend {
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
				if !StartChatBackend && !StartWebServer {
					fmt.Print("\n")
					state.Logger.Info("Gracefully terminating...")
					_ = state.ChatBackendProcess.Signal(syscall.SIGUSR2)
					_ = state.WebServerProcess.Signal(syscall.SIGUSR2)
					_ = SocketListener.Close()
					_ = syscall.Unlink(config.SocketPath)
				} else if StartChatBackend {
					if currentSignal != os.Interrupt {
						// Only tell the supervisor if SIGUSR2 was not sent or the program was not interrupted
						_ = state.IPCConnection.Call("IPC.Restart", []string{"ChatBackend"}, nil)
					}
					for _, shardSession := range state.ShardSessions {
						_ = shardSession.Close()
					}
					_ = state.RawSession.Close()
				} else if StartWebServer {
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
	if StartChatBackend && !StartWebServer && exitCode != 0 {
		_ = state.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
	} else if StartWebServer && !StartChatBackend && exitCode != 0 {
		_ = state.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
	}
	os.Exit(exitCode)
}
