package main

import (
	"flag"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/web"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
)

const Version = "v0.0.1-rewrite"

var StartChatBackend bool
var StartWebServer bool

func init() {
	flag.StringVar(&multiplexer.RawSession.Token, "a", "", "Discord Authorization Token")
	flag.BoolVar(&StartChatBackend, "c", false, "Start the chat backend directly")
	flag.BoolVar(&StartWebServer, "w", false, "Start the web server directly")
}

func main() {
	// Some regular initialization
	var err error
	var readyChannel = make(chan bool, 1)
	var SocketListener net.Listener
	var IPCFunctions = new(multiplexer.IPC)
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
			multiplexer.IPCConnection, err = rpc.DialHTTP("unix", config.SocketPath)
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Failed to connect to the database, %s", err))
				os.Exit(1)
			}

			// Authenticate and make session
			if multiplexer.RawSession.Token == "" {
				configToken := config.Config.Section("System").Key("Token").String()
				if configToken != "" && configToken != "INSERT_TOKEN_HERE" {
					if config.Debug {
						multiplexer.Logger.Debug("Loaded token from configuration file.")
					}
					multiplexer.RawSession.Token = configToken
				} else {
					multiplexer.Logger.Error("Please specify an authorization token.")
					_ = multiplexer.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
					os.Exit(1)
				}
			} else {
				if config.Debug {
					multiplexer.Logger.Error("Loaded token from command parameter.")
				}
			}

			multiplexer.RawSession.UserAgent = "DiscordBot (FreeNitori " + Version + ")"
			multiplexer.RawSession.Token = "Bot " + multiplexer.RawSession.Token
			multiplexer.RawSession.ShouldReconnectOnError = true
			err = multiplexer.RawSession.Open()
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("An error occurred while connecting to Discord, %s", err))
				os.Exit(1)
			}
			multiplexer.Initialized = true
			multiplexer.Application, err = multiplexer.RawSession.Application("@me")
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("An error occurred while fetching application info, %s", err))
				os.Exit(1)
			}
			_, _ = multiplexer.RawSession.UserUpdateStatus("dnd")
			_ = multiplexer.RawSession.UpdateStatus(0, config.Presence)
			if config.Shard {
				multiplexer.MakeSessions()
			}

			// Log into the logger that the ChatBackend is ready to go
			_ = multiplexer.IPCConnection.Call("IPC.Log", []string{
				"INFO",
				fmt.Sprintf("User: %s | ID: %s | Default Prefix: %s",
					multiplexer.RawSession.State.User.Username+"#"+multiplexer.RawSession.State.User.Discriminator,
					multiplexer.RawSession.State.User.ID,
					config.Prefix),
			}, nil)
			_ = multiplexer.IPCConnection.Call("IPC.Log", []string{
				"INFO",
				"FreeNitori is now ready. Press Control-C to terminate.",
			}, nil)
			_ = multiplexer.IPCConnection.Call("IPC.SignalWebServer", []string{}, nil)
		}
	case StartWebServer:
		{

			// Dial the supervisor socket
			multiplexer.IPCConnection, err = rpc.DialHTTP("unix", config.SocketPath)
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Failed to connect to the supervisor process, %s", err))
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Unable to establish connection with database, %s", err))
				os.Exit(1)
			}

			// Initialize and start the server
			web.Initialize()
			go func() {
				<-readyChannel
				err = web.Engine.Run(fmt.Sprintf("%s:%s", config.Host, config.Port))
				if err != nil {
					multiplexer.Logger.Error(fmt.Sprintf("Failed to start web server, %s", err))
					_ = multiplexer.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
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
`+"\n", Version)

			// Check for an existing instance
			if _, err := os.Stat(config.SocketPath); os.IsNotExist(err) {
			} else {
				_, err := net.Dial("unix", config.SocketPath)
				if err != nil {
					err = syscall.Unlink(config.SocketPath)
					if err != nil {
						multiplexer.Logger.Error(fmt.Sprintf("Unable to remove hanging socket, %s", err))
						os.Exit(1)
					}
				} else {
					multiplexer.Logger.Error("Another instance of FreeNitori is already running.")
					os.Exit(1)
				}
			}

			// Initialize the socket
			_ = rpc.Register(IPCFunctions)
			rpc.HandleHTTP()
			SocketListener, err = net.Listen("unix", config.SocketPath)
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Failed to listen on the socket, %s", err))
				os.Exit(1)
			}
			go http.Serve(SocketListener, nil)

			// Create the chat backend process
			multiplexer.ChatBackendProcess, err =
				os.StartProcess(multiplexer.ExecPath, []string{multiplexer.ExecPath, "-c", "-a", multiplexer.RawSession.Token}, &multiplexer.ProcessAttributes)
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Failed to create chat backend process, %s", err))
				os.Exit(1)
			}

			// Create web server process
			multiplexer.WebServerProcess, err =
				os.StartProcess(multiplexer.ExecPath, []string{multiplexer.ExecPath, "-w"}, &multiplexer.ProcessAttributes)
			if err != nil {
				multiplexer.Logger.Error(fmt.Sprintf("Failed to create web server process, %s", err))
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
					multiplexer.ChatBackendIPCReceiver()
				} else if StartWebServer && !StartChatBackend {
					if !multiplexer.Initialized {
						readyChannel <- true
						multiplexer.Initialized = true
					}
				}
			case syscall.SIGUSR2:
				multiplexer.ExitCode <- 0
				return
			default:
				// Cleanup stuffs
				if !StartChatBackend && !StartWebServer {
					fmt.Print("\n")
					multiplexer.Logger.Info("Gracefully terminating...")
					_ = multiplexer.ChatBackendProcess.Signal(syscall.SIGUSR2)
					_ = multiplexer.WebServerProcess.Signal(syscall.SIGUSR2)
					_ = SocketListener.Close()
					_ = syscall.Unlink(config.SocketPath)
				} else if StartChatBackend {
					if currentSignal != os.Interrupt {
						// Only tell the supervisor if SIGUSR2 was not sent or the program was not interrupted
						_ = multiplexer.IPCConnection.Call("IPC.Restart", []string{"ChatBackend"}, nil)
					}
					for _, session := range multiplexer.DiscordSessions {
						_ = session.Close()
					}
					_ = multiplexer.RawSession.Close()
				} else if StartWebServer {
					if currentSignal != os.Interrupt {
						// Only write the packet if SIGUSR2 was not sent or the program was not interrupted
						_ = multiplexer.IPCConnection.Call("IPC.Restart", []string{"WebServer"}, nil)
					}
				}
				multiplexer.ExitCode <- 0
				return
			}
		}
	}()

	// Tell the Supervisor and exit if there's something on that channel
	exitCode := <-multiplexer.ExitCode
	if StartChatBackend && !StartWebServer && exitCode != 0 {
		_ = multiplexer.IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
	} else if StartWebServer && !StartChatBackend && exitCode != 0 {
		_ = multiplexer.IPCConnection.Call("IPC.Error", []string{"WebServer"}, nil)
	}
	os.Exit(exitCode)
}
