package communication

import (
	"encoding/json"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/session"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"syscall"
)

func ChatBackendIPCReceiver() {
	var instruction string
	var response string
	_ = state.IPCConnection.Call("IPC.ChatBackendIPCResponder", []string{"furtherInstruction"}, &instruction)
	switch instruction {
	case "totalGuilds":
		response = strconv.Itoa(len(state.RawSession.State.Guilds))
	case "version":
		response = state.Version
	case "inviteURL":
		response = fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=2146958847",
			state.Application.ID)
	}
	if strings.HasPrefix(instruction, "GuildInfo") {
		var members []*UserInfo
		gid := instruction[9:]
		guildSession, err := session.FetchGuildSession(gid)
		if err == nil {
			var guild *discordgo.Guild
			for _, guildIter := range guildSession.State.Guilds {
				if guildIter.ID == gid {
					guild = guildIter
					break
				}
			}
			if guild != nil {
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
	_ = state.IPCConnection.Call("IPC.ChatBackendIPCResponder", []string{"response", response}, nil)
}

func (ipc *IPC) SignalWebServer(args []string, reply *int) error {
	args = nil
	reply = nil
	return state.WebServerProcess.Signal(syscall.SIGUSR1)
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
