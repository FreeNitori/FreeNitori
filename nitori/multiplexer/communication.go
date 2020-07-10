package multiplexer

import (
	"encoding/json"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/bwmarrin/discordgo"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var ExitCode = make(chan int)
var Initialized bool
var IPCConnection net.Conn
var RawSession, _ = discordgo.New()
var DiscordSessions []*discordgo.Session
var ExecPath string
var err error
var WebServerProcess os.Process
var ChatBackendProcess *os.Process
var ProcessAttributes = os.ProcAttr{
	Dir: ".",
	Env: os.Environ(),
	Files: []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	},
}

type IPCPacket struct {
	IssuerIdentifier   string   `json:"issuer_identifier"`
	ReceiverIdentifier string   `json:"receiver_identifier"`
	MessageIdentifier  string   `json:"message_identifier"`
	Body               []string `json:"body"`
}

func init() {
	ExecPath, err = os.Executable()
	if err != nil {
		log.Printf("Failed to get FreeNitori's executable path, %s", err)
		os.Exit(1)
	}
}

func WritePacket(connection net.Conn, outgoingPacket IPCPacket) (incomingPacket IPCPacket) {
	var err error
	var incoming IPCPacket

	jsonEncoder := json.NewEncoder(connection)
	jsonDecoder := json.NewDecoder(connection)

	// Encode the outgoing packet
	err = jsonEncoder.Encode(outgoingPacket)
	if err != nil {
		log.Printf("Failed to encode packet, %s", err)
		return incoming
	}

	// Decode and return the incoming packet
	err = jsonDecoder.Decode(&incoming)
	if err != nil {
		log.Printf("Failed to decode packet, %s", err)
		return incoming
	}
	return incoming
}

func (incomingPacket IPCPacket) SupervisorPacketHandler() (outgoingPacket IPCPacket) {
	switch incomingPacket.IssuerIdentifier {
	case "ChatBackendInitializer":

		// This should never talk to anything other than the supervisor
		if incomingPacket.ReceiverIdentifier != "Supervisor" {
			log.Println("Invalid packet from Chat backend initializer.")
			ExitCode <- 1
		}
		switch incomingPacket.MessageIdentifier {
		case "ChatBackendInitializationFinish":
			// Print out the message if not initialized and set the finish variable
			if Initialized {
				return IPCPacket{
					IssuerIdentifier:   "Supervisor",
					ReceiverIdentifier: incomingPacket.IssuerIdentifier,
					MessageIdentifier:  "MessageAcknowledgement",
					Body:               []string{"We have already initialized."},
				}
			}
			Initialized = true
			log.Printf("User: %s | ID: %s | Prefix: %s",
				incomingPacket.Body[0],
				incomingPacket.Body[1],
				incomingPacket.Body[2])
			log.Printf("FreeNitori is now ready. Press Control-C to terminate.")
		case "AbnormalExit":
			// Stop the supervisor
			if incomingPacket.Body[0] == "1" {
				log.Println("Chat backend has encountered an error, proceeding to exit.")
				// TODO: kill web server
				ExitCode <- 1
			}
		}
	case "ChatBackend":
		switch incomingPacket.ReceiverIdentifier {
		case "Supervisor":
			switch incomingPacket.MessageIdentifier {
			case "RouteLog":
				if config.Shard {
					log.Printf("(Shard %s) %s@%s > %s", incomingPacket.Body[3], incomingPacket.Body[0], incomingPacket.Body[1], incomingPacket.Body[2])
				} else {
					log.Printf("%s@%s > %s", incomingPacket.Body[0], incomingPacket.Body[1], incomingPacket.Body[2])
				}
			case "Reboot":
				// Recreate the chat backend process
				go func() {
					_, _ = ChatBackendProcess.Wait()
					ChatBackendProcess, err =
						os.StartProcess(ExecPath, []string{ExecPath, "-c", "-a", RawSession.Token}, &ProcessAttributes)
					if err != nil {
						log.Printf("Failed to recreate chat backend process, %s", err)
						ExitCode <- 0
					} else {
						log.Println("Chat backend has been restarted.")
					}
				}()
			case "FullShutdown":
				// Kill the web server and go down
				// TODO: kill web server
				log.Println("Graceful shutdown initiated by chat backend.")
				ExitCode <- 0
			default:
				return IPCPacket{
					IssuerIdentifier:   "Supervisor",
					ReceiverIdentifier: incomingPacket.IssuerIdentifier,
					MessageIdentifier:  "Error",
					Body:               []string{"Unknown message."},
				}
			}

		default:
			return IPCPacket{
				IssuerIdentifier:   "Supervisor",
				ReceiverIdentifier: incomingPacket.IssuerIdentifier,
				MessageIdentifier:  "Error",
				Body:               []string{"Unknown issuer."},
			}
		}
	default:
		return IPCPacket{
			IssuerIdentifier:   "Supervisor",
			ReceiverIdentifier: incomingPacket.IssuerIdentifier,
			MessageIdentifier:  "Error",
			Body:               []string{"Unknown issuer."},
		}
	}
	return IPCPacket{
		IssuerIdentifier:   "Supervisor",
		ReceiverIdentifier: incomingPacket.ReceiverIdentifier,
		MessageIdentifier:  "MessageAcknowledgement",
		Body:               []string{"Request has been received."},
	}
}

func MakeSessions() {
	var err error

	// Will do something about this later
	if config.ShardCount < 1 {
		gatewayBot, err := RawSession.GatewayBot()
		if err != nil {
			WritePacket(
				IPCConnection,
				IPCPacket{
					IssuerIdentifier:   "ChatBackendInitializer",
					ReceiverIdentifier: "Supervisor",
					MessageIdentifier:  "AbnormalExit",
					Body:               []string{"1"},
				})
			_ = IPCConnection.Close()
			os.Exit(1)
		}
		config.ShardCount = gatewayBot.Shards
	}

	// Make sure it doesn't end up being 0 shards
	if config.ShardCount == 0 {
		config.ShardCount = 1
	}

	// Make the sessions
	for i := 0; i < config.ShardCount; i++ {
		time.Sleep(2)
		session, _ := discordgo.New()
		session.ShardCount = config.ShardCount
		session.ShardID = i
		session.Token = RawSession.Token
		session.UserAgent = RawSession.UserAgent
		err = session.Open()
		if err != nil {
			log.Printf("Failed to open session %s, %s", strconv.Itoa(i), err)
			ExitCode <- 1
		}
		session.AddHandler(Router.OnMessageCreate)
		DiscordSessions = append(DiscordSessions, session)
	}
}
