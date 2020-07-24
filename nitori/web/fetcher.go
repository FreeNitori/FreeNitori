package web

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

func fetchData(request string) string {
	var response string
	_ = multiplexer.IPCConnection.Call("IPC.RequestData", []string{"ChatBackend", request}, &response)
	return response
}

func fetchGuild(gid string) *multiplexer.GuildInfo {
	var response *multiplexer.GuildInfo
	err = multiplexer.IPCConnection.Call("IPC.RequestGuild", []string{gid}, &response)
	if err != nil {
		return nil
	}
	return response
}
