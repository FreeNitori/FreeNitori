package multiplexer

import (
	"encoding/json"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/config"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strconv"
	"strings"
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

func FetchGuildSession(gid string) (*discordgo.Session, error) {
	if !config.Shard {
		return RawSession, nil
	}
	ID, err := strconv.ParseInt(gid, 10, 64)
	if err != nil {
		return nil, err
	}
	return DiscordSessions[(ID>>22)%int64(config.ShardCount)], nil
}

func ChatBackendIPCReceiver() {
	var instruction string
	var response string
	_ = IPCConnection.Call("IPC.ChatBackendIPCResponder", []string{"furtherInstruction"}, &instruction)
	switch instruction {
	case "totalGuilds":
		response = strconv.Itoa(len(RawSession.State.Guilds))
	case "inviteURL":
		response = fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=2146958847",
			Application.ID)
	}
	if strings.HasPrefix(instruction, "GuildInfo") {
		var members []*UserInfo
		gid := instruction[9:]
		guildSession, err := FetchGuildSession(gid)
		if err == nil {
			guild, err := guildSession.Guild(gid)
			if err == nil {
				for _, member := range guild.Members {
					userInfo := UserInfo{
						Name:          member.User.Username,
						ID:            member.User.ID,
						AvatarURL:     member.User.AvatarURL("128"),
						Discriminator: member.User.Discriminator,
					}
					members = append(members, &userInfo)
				}
				responseBytes, err := json.Marshal(GuildInfo{
					Name:    guild.Name,
					ID:      guild.ID,
					IconURL: guild.IconURL(),
					Members: members,
				})
				if err == nil {
					response = string(responseBytes)
				}
			}
		}
	}
	_ = IPCConnection.Call("IPC.ChatBackendIPCResponder", []string{"response", response}, nil)
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
