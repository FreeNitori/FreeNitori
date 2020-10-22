package main

import (
	"encoding/json"
	"fmt"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/ipc"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
	"git.randomchars.net/RandomChars/FreeNitori/proc/chatbackend/state"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func ChatBackendIPCReceiver() {
	var instruction string
	var response string
	_ = vars.RPCConnection.Call("R.ChatBackendIPCResponder", []string{"furtherInstruction"}, &instruction)
	switch instruction {
	case "totalGuilds":
		response = strconv.Itoa(len(state.RawSession.State.Guilds))
	case "version":
		response = vars.Version
	case "inviteURL":
		response = fmt.Sprintf("https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot&permissions=2146958847",
			state.Application.ID)
	}
	if strings.HasPrefix(instruction, "GuildInfo") {
		var members []*ipc.UserInfo
		gid := instruction[9:]
		guildSession, err := FetchGuildSession(gid)
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
					userInfo := ipc.UserInfo{
						Name:          member.User.Username,
						ID:            member.User.ID,
						AvatarURL:     member.User.AvatarURL("128"),
						Discriminator: member.User.Discriminator,
						Bot:           member.User.Bot,
					}
					members = append(members, &userInfo)
				}
				responseBytes, err := json.Marshal(ipc.GuildInfo{
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
	_ = vars.RPCConnection.Call("R.ChatBackendIPCResponder", []string{"response", response}, nil)
}
