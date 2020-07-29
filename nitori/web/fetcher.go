package web

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/communication"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/state"
)

func fetchData(request string) string {
	var response string
	_ = state.IPCConnection.Call("IPC.RequestData", []string{"ChatBackend", request}, &response)
	return response
}

func fetchGuild(gid string) *communication.GuildInfo {
	var response *communication.GuildInfo
	err = state.IPCConnection.Call("IPC.RequestGuild", []string{gid}, &response)
	if err != nil {
		return nil
	}
	return response
}
