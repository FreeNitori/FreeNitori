package web

import "git.randomchars.net/RandomChars/FreeNitori/nitori/multiplexer"

func askForData(request string) string {
	var response string
	_ = multiplexer.IPCConnection.Call("IPC.RequestData", []string{"ChatBackend", request}, &response)
	return response
}
