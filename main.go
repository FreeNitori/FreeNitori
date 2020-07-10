package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
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
	var SocketListener net.Listener
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
			multiplexer.IPCConnection, err = net.Dial("unix", config.SocketPath)
			if err != nil {
				log.Printf("Failed to connect to the supervisor process, %s", err)
				os.Exit(1)
			}

			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				log.Printf("Unable to establish connection with database, %s", err)
				os.Exit(1)
			}

			// Authenticate and make session
			if multiplexer.RawSession.Token == "" {
				configToken := config.Config.Section("System").Key("Token").String()
				if configToken != "" && configToken != "INSERT_TOKEN_HERE" {
					if config.Debug {
						log.Println("Loaded token from configuration file.")
					}
					multiplexer.RawSession.Token = configToken
				} else {
					log.Println("Please specify an authorization token.")
					multiplexer.WritePacket(
						multiplexer.IPCConnection,
						multiplexer.IPCPacket{
							IssuerIdentifier:   "ChatBackendInitializer",
							ReceiverIdentifier: "Supervisor",
							MessageIdentifier:  "AbnormalExit",
							Body:               []string{strconv.Itoa(1)},
						})
					os.Exit(1)
				}
			} else {
				if config.Debug {
					log.Println("Loaded token from command parameter.")
				}
			}

			multiplexer.RawSession.UserAgent = "DiscordBot (FreeNitori " + Version + ")"
			multiplexer.RawSession.Token = "Bot " + multiplexer.RawSession.Token
			err = multiplexer.RawSession.Open()
			if err != nil {
				log.Printf("An error occurred while connecting to Discord, %s \n", err)
				os.Exit(1)
			}
			_, _ = multiplexer.RawSession.UserUpdateStatus("dnd")
			_ = multiplexer.RawSession.UpdateStatus(0, config.Presence)
			if config.Shard {
				multiplexer.MakeSessions()
			}

			// Tell the supervisor we are ready to go
			_ = multiplexer.WritePacket(multiplexer.IPCConnection, multiplexer.IPCPacket{
				IssuerIdentifier:   "ChatBackendInitializer",
				ReceiverIdentifier: "Supervisor",
				MessageIdentifier:  "ChatBackendInitializationFinish",
				Body: []string{
					multiplexer.RawSession.State.User.Username + "#" + multiplexer.RawSession.State.User.Discriminator,
					multiplexer.RawSession.State.User.ID,
					config.Prefix},
			})
		}
	case StartWebServer:
		{
			// Check the database
			_, err = config.Redis.Ping(config.RedisContext).Result()
			if err != nil {
				log.Printf("Unable to establish connection with database, %s", err)
				os.Exit(1)
			}
			os.Exit(0)
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
						log.Printf("Unable to remove hanging socket, %s", err)
						os.Exit(1)
					}
				} else {
					log.Println("Another instance of FreeNitori is already running.")
					os.Exit(1)
				}
			}

			// Initialize the socket
			SocketListener, err = net.Listen("unix", config.SocketPath)
			if err != nil {
				log.Printf("Failed to listen on the socket, %s", err)
				os.Exit(1)
			}

			// Function that monitors the socket and responds
			go func() {
				for {
					descriptor, err := SocketListener.Accept()
					if err != nil {
						return
					}
					go func(connection net.Conn) {
						// defer connection.Close()
						jsonEncoder := json.NewEncoder(connection)
						jsonDecoder := json.NewDecoder(connection)
						for {
							var packet multiplexer.IPCPacket
							err := jsonDecoder.Decode(&packet)
							if err != nil {
								if err == io.EOF {
									break
								}
								log.Printf("Failed to decode packet, %s", err)
								continue
							}
							err = jsonEncoder.Encode(packet.SupervisorPacketHandler())
							if err != nil {
								log.Printf("Failed to encode packet, %s", err)
								continue
							}
						}
					}(descriptor)
				}
			}()

			// Create the chat backend process
			multiplexer.ChatBackendProcess, err =
				os.StartProcess(multiplexer.ExecPath, []string{multiplexer.ExecPath, "-c", "-a", multiplexer.RawSession.Token}, &multiplexer.ProcessAttributes)
			if err != nil {
				log.Printf("Failed to create chat backend process, %s", err)
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
				// Go to the supervisor to fetch new message
				if StartChatBackend && !StartWebServer {
					// TODO: do the stuff
				}
			default:
				// Cleanup stuffs
				if !StartChatBackend && !StartWebServer {
					fmt.Print("\n")
					log.Println("Gracefully terminating...")
					_ = multiplexer.ChatBackendProcess.Signal(syscall.SIGUSR2)
					_ = SocketListener.Close()
					_ = syscall.Unlink(config.SocketPath)
				} else if StartChatBackend {
					if currentSignal != syscall.SIGUSR2 && currentSignal != os.Interrupt {
						// Only write the packet if SIGUSR2 was not sent or the program was not interrupted
						multiplexer.WritePacket(
							multiplexer.IPCConnection,
							multiplexer.IPCPacket{
								IssuerIdentifier:   "ChatBackendInitializer",
								ReceiverIdentifier: "Supervisor",
								MessageIdentifier:  "KillSignal",
								Body:               []string{currentSignal.String()},
							})
					}
					for _, session := range multiplexer.DiscordSessions {
						_ = session.Close()
					}
					_ = multiplexer.RawSession.Close()
				}
				multiplexer.ExitCode <- 0
				return
			}
		}
	}()

	// Send packet and exit if there's something on that channel
	exitCode := <-multiplexer.ExitCode
	if StartChatBackend && !StartWebServer && exitCode != 0 {
		multiplexer.WritePacket(
			multiplexer.IPCConnection,
			multiplexer.IPCPacket{
				IssuerIdentifier:   "ChatBackendInitializer",
				ReceiverIdentifier: "Supervisor",
				MessageIdentifier:  "AbnormalExit",
				Body:               []string{strconv.Itoa(exitCode)},
			})
	}
	os.Exit(exitCode)
}
