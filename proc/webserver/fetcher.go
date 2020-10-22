package main

import (
	"git.randomchars.net/RandomChars/FreeNitori/nitori/ipc"
	"git.randomchars.net/RandomChars/FreeNitori/nitori/vars"
)

func fetchData(request string) string {
	var response string
	_ = vars.RPCConnection.Call("R.RequestData", []string{"ChatBackend", request}, &response)
	return response
}

func fetchGuild(gid string) *ipc.GuildInfo {
	var response *ipc.GuildInfo
	err = vars.RPCConnection.Call("R.RequestGuild", []string{gid}, &response)
	if err != nil {
		return nil
	}
	return response
}
