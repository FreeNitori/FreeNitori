package multiplexer

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
)

var Initialized = false
var Application *discordgo.Application

func MakeSessions() {
	var err error

	// Get recommended shard count from Discord
	if config.ShardCount < 1 {
		gatewayBot, err := RawSession.GatewayBot()
		if err != nil {
			_ = IPCConnection.Call("IPC.Error", []string{"ChatBackend"}, nil)
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
		session.ShouldReconnectOnError = RawSession.ShouldReconnectOnError
		err = session.Open()
		if err != nil {
			log.Printf("Failed to open session %s, %s", strconv.Itoa(i), err)
			ExitCode <- 1
		}
		session.AddHandler(Router.OnMessageCreate)
		DiscordSessions = append(DiscordSessions, session)
	}
}

func (ipc *IPC) SignalWebServer(args []string, reply *int) error {
	args = nil
	reply = nil
	return WebServerProcess.Signal(syscall.SIGUSR1)
}

func (ipc *IPC) ChatBackendIPCResponder(args []string, reply *string) error {
	switch args[0] {
	case "furtherInstruction":
		*reply = <-RequestInstructionChannel
	case "response":
		RequestDataChannel <- args[1]
	}
	return nil
}
